// Command vendor-assets downloads the third-party files listed in
// web/assets/vendor.json from jsdelivr, so they are not committed
//
//	go run ./tools/vendor-assets               download missing or changed files
//	go run ./tools/vendor-assets --check       list packages with a newer version, exit 1 if any
//	go run ./tools/vendor-assets --update [pkg...]  move packages (all if none named) to their latest version
//
// web/assets/vendor.lock.json keeps the sha256 of every file for the pinned
// version, a download that does not match fails, changing a version (by hand
// or --update) records the hashes of the new files
//
// run from the module root (app/)
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"text/tabwriter"
	"time"
)

// overridden in tests
var (
	cdnURL      = "https://cdn.jsdelivr.net/npm"
	registryURL = "https://registry.npmjs.org"
)

//exhaustruct:ignore
var client = &http.Client{Timeout: 60 * time.Second}

// manifest is vendor.json: npm package -> pinned version and files,
// files map a path in the package to a path under the assets directory
type manifest map[string]pkg

type pkg struct {
	// empty pins the latest version on the next run
	Version string            `json:"version"`
	Files   map[string]string `json:"files"`
}

// lock is vendor.lock.json, written by this tool
type lock map[string]lockedPkg

type lockedPkg struct {
	Version string `json:"version"`
	// path in the package -> sha256
	SHA256 map[string]string `json:"sha256"`
}

func main() {
	dir := flag.String("dir", "web/assets", "directory with vendor.json, destination paths are relative to it")
	check := flag.Bool("check", false, "list packages with a newer version, exit 1 if any")
	update := flag.Bool("update", false, "move the named packages, or all, to their latest version")
	flag.Parse()

	var err error
	switch {
	case *check:
		var outdated bool
		outdated, err = runCheck(*dir, os.Stdout)
		if err == nil && outdated {
			os.Exit(1)
		}
	case *update:
		err = runUpdate(*dir, flag.Args())
	default:
		err = runSync(*dir)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "vendor-assets:", err)
		os.Exit(2)
	}
}

func runSync(dir string) error {
	m, l, err := load(dir)
	if err != nil {
		return err
	}
	if err := pinLatest(dir, m, func(p pkg) bool { return p.Version == "" }); err != nil {
		return err
	}
	return syncAll(dir, m, l)
}

func runUpdate(dir string, names []string) error {
	m, l, err := load(dir)
	if err != nil {
		return err
	}
	for _, n := range names {
		if _, ok := m[n]; !ok {
			return fmt.Errorf("%s is not in vendor.json", n)
		}
	}
	err = pinLatest(dir, m, func(p pkg) bool { return true }, names...)
	if err != nil {
		return err
	}
	return syncAll(dir, m, l)
}

func runCheck(dir string, out io.Writer) (outdated bool, err error) {
	m, _, err := load(dir)
	if err != nil {
		return false, err
	}
	// write errors come back from Flush
	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(w,"package\tpinned\tlatest\t")
	for _, name := range sortedKeys(m) {
		latest, err := latestVersion(name)
		if err != nil {
			return false, err
		}
		mark := ""
		if latest != m[name].Version {
			mark, outdated = "outdated", true
		}
		_, _ = fmt.Fprintf(w,"%s\t%s\t%s\t%s\n", name, m[name].Version, latest, mark)
	}
	return outdated, w.Flush()
}

// pinLatest sets the version of the packages selected by want (only names, if given)
// to the latest and saves vendor.json if anything changed
func pinLatest(dir string, m manifest, want func(pkg) bool, names ...string) error {
	changed := false
	for _, name := range sortedKeys(m) {
		p := m[name]
		if (len(names) > 0 && !slices.Contains(names, name)) || !want(p) {
			continue
		}
		latest, err := latestVersion(name)
		if err != nil {
			return err
		}
		if latest != p.Version {
			fmt.Printf("%s: %s -> %s\n", name, p.Version, latest)
			p.Version = latest
			m[name] = p
			changed = true
		}
	}
	if !changed {
		return nil
	}
	return writeJSON(filepath.Join(dir, "vendor.json"), m)
}

// syncAll downloads every file that is missing or does not match the lock,
// then saves the lock
func syncAll(dir string, m manifest, l lock) error {
	newLock := lock{}
	for _, name := range sortedKeys(m) {
		p := m[name]
		locked, sameVersion := l[name]
		sameVersion = sameVersion && locked.Version == p.Version
		sums := map[string]string{}
		for _, src := range sortedKeys(p.Files) {
			dst := filepath.Join(dir, filepath.FromSlash(p.Files[src]))
			want := ""
			if sameVersion {
				want = locked.SHA256[src]
			}
			if want != "" && fileSHA256(dst) == want {
				sums[src] = want
				continue
			}
			body, err := download(fmt.Sprintf("%s/%s@%s/%s", cdnURL, name, p.Version, src))
			if err != nil {
				return err
			}
			sum := sha256Hex(body)
			if want != "" && sum != want {
				return fmt.Errorf("%s@%s/%s: sha256 %s does not match the lock (%s), change the version to re-lock",
					name, p.Version, src, sum, want)
			}
			if err := writeAtomic(dst, body); err != nil {
				return err
			}
			fmt.Printf("%s@%s/%s -> %s\n", name, p.Version, src, p.Files[src])
			sums[src] = sum
		}
		newLock[name] = lockedPkg{Version: p.Version, SHA256: sums}
	}
	return writeJSON(filepath.Join(dir, "vendor.lock.json"), newLock)
}

func load(dir string) (manifest, lock, error) {
	var m manifest
	if err := readJSON(filepath.Join(dir, "vendor.json"), &m); err != nil {
		return nil, nil, err
	}
	l := lock{}
	err := readJSON(filepath.Join(dir, "vendor.lock.json"), &l)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, nil, err
	}
	return m, l, nil
}

func latestVersion(name string) (string, error) {
	body, err := download(registryURL + "/" + name + "/latest")
	if err != nil {
		return "", err
	}
	var v struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(body, &v); err != nil || v.Version == "" {
		return "", fmt.Errorf("%s: no version in the registry response", name)
	}
	return v.Version, nil
}

func download(url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck // the body is fully read or the request failed
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: %s", url, resp.Status)
	}
	return io.ReadAll(resp.Body)
}

// writeAtomic writes through a temp file, a failed run never leaves half a file
func writeAtomic(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".vendor-*")
	if err != nil {
		return err
	}
	_, err = tmp.Write(data)
	if closeErr := tmp.Close(); err == nil {
		err = closeErr
	}
	if err == nil {
		err = os.Rename(tmp.Name(), path)
	}
	if err != nil {
		_ = os.Remove(tmp.Name())
	}
	return err
}

func fileSHA256(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return sha256Hex(data)
}

func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func readJSON(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// writeJSON sorts map keys (encoding/json does), so the files diff cleanly
func writeJSON(path string, v any) error {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return err
	}
	return writeAtomic(path, buf.Bytes())
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}
