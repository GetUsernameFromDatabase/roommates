package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
)

// fakeNPM serves pkg@1.0.0 and pkg@2.0.0 (latest) with one file each
func fakeNPM(t *testing.T) *atomic.Int32 {
	t.Helper()
	var downloads atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/npm/pkg@1.0.0/dist/a.js":
			downloads.Add(1)
			_, _ = w.Write([]byte("v1"))
		case "/npm/pkg@2.0.0/dist/a.js":
			downloads.Add(1)
			_, _ = w.Write([]byte("v2"))
		case "/registry/pkg/latest":
			_, _ = w.Write([]byte(`{"version":"2.0.0"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	oldCDN, oldRegistry := cdnURL, registryURL
	cdnURL, registryURL = srv.URL+"/npm", srv.URL+"/registry"
	t.Cleanup(func() { cdnURL, registryURL = oldCDN, oldRegistry })
	return &downloads
}

func writeManifest(t *testing.T, dir, version string) {
	t.Helper()
	m := manifest{"pkg": {Version: version, Files: map[string]string{"dist/a.js": "lib/a.js"}}}
	if err := writeJSON(filepath.Join(dir, "vendor.json"), m); err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestSync(t *testing.T) {
	downloads := fakeNPM(t)
	dir := t.TempDir()
	asset := filepath.Join(dir, "lib", "a.js")
	writeManifest(t, dir, "1.0.0")

	if err := runSync(dir); err != nil {
		t.Fatal(err)
	}
	if readFile(t, asset) != "v1" || downloads.Load() != 1 {
		t.Fatalf("first sync: file %q, %d downloads", readFile(t, asset), downloads.Load())
	}

	// present and matching the lock: nothing is downloaded
	if err := runSync(dir); err != nil {
		t.Fatal(err)
	}
	if downloads.Load() != 1 {
		t.Fatalf("second sync downloaded again, %d downloads", downloads.Load())
	}

	// a download that does not match the lock fails and keeps the old file
	if err := os.WriteFile(asset, []byte("edited"), 0o644); err != nil {
		t.Fatal(err)
	}
	lockPath := filepath.Join(dir, "vendor.lock.json")
	lockText := readFile(t, lockPath)
	tampered := strings.Replace(lockText, sha256Hex([]byte("v1")), sha256Hex([]byte("other")), 1)
	if err := os.WriteFile(lockPath, []byte(tampered), 0o644); err != nil {
		t.Fatal(err)
	}
	err := runSync(dir)
	if err == nil || !strings.Contains(err.Error(), "does not match the lock") {
		t.Fatalf("mismatching download: got error %v", err)
	}
	if readFile(t, asset) != "edited" {
		t.Fatal("a failed sync must not replace the file")
	}

	// changing the version re-locks
	writeManifest(t, dir, "2.0.0")
	if err := runSync(dir); err != nil {
		t.Fatal(err)
	}
	if readFile(t, asset) != "v2" || !strings.Contains(readFile(t, lockPath), sha256Hex([]byte("v2"))) {
		t.Fatal("a new version should be downloaded and locked")
	}
}

func TestCheckAndUpdate(t *testing.T) {
	fakeNPM(t)
	dir := t.TempDir()
	writeManifest(t, dir, "1.0.0")

	var out strings.Builder
	outdated, err := runCheck(dir, &out)
	if err != nil || !outdated || !strings.Contains(out.String(), "outdated") {
		t.Fatalf("check: outdated=%v err=%v\n%s", outdated, err, out.String())
	}

	if err := runUpdate(dir, []string{"pkg"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(readFile(t, filepath.Join(dir, "vendor.json")), `"version": "2.0.0"`) ||
		readFile(t, filepath.Join(dir, "lib", "a.js")) != "v2" {
		t.Fatal("update should pin and download the latest version")
	}
	if err := runUpdate(dir, []string{"missing"}); err == nil {
		t.Fatal("update of a package not in vendor.json should fail")
	}

	// an empty version pins the latest
	writeManifest(t, dir, "")
	if err := runSync(dir); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(readFile(t, filepath.Join(dir, "vendor.json")), `"version": "2.0.0"`) {
		t.Fatal("an empty version should be pinned to the latest")
	}
}
