//go:generate go run .

// generates locales and checks if other language files than default have missing keys

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// used to reference this file
var _, thisFile, _, _ = runtime.Caller(0)

func resolve(rel string) string {
	return filepath.Join(filepath.Dir(thisFile), filepath.FromSlash(rel))
}

// Paths relative to this file's package directory.
const (
	defaultLang = "et"

	toLocaleFolder = "../../locales"
	outFile        = "locales-languages.go"
)

// command flag arguments
var (
	baseLocaleFile = resolve(toLocaleFolder + "/et/et.yaml")
	localeFolder   = resolve(toLocaleFolder)
)

func main() {
	if thisFile == "" {
		panic("thisFile is used to resolve other files")
	}

	base, err := loadYAML(baseLocaleFile)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"%s: %v\n",
			baseLocaleFile,
			err,
		)
		os.Exit(2)
	}

	base, err = unwrapLocale(base)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"%s: %v\n",
			baseLocaleFile,
			err,
		)
		os.Exit(2)
	}

	files, err := filepath.Glob(
		filepath.Join(localeFolder, "*", "*.yaml"),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	var failed bool

	for _, file := range files {
		if file == baseLocaleFile {
			continue
		}

		target, err := loadYAML(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", file, err)
			failed = true
			continue
		}

		target, err = unwrapLocale(target)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", file, err)
			failed = true
			continue
		}

		missing := missingKeys(base, target, "")
		if len(missing) == 0 {
			fmt.Printf("OK: %s\n", file)
			continue
		}

		failed = true
		fmt.Printf("MISSING KEYS: %s\n", file)

		for _, key := range missing {
			fmt.Printf("  - %s\n", key)
		}
	}

	langs, err := discoverLanguages(localeFolder)
	if err != nil {
		fmt.Fprintf(os.Stderr, "discover languages: %v\n", err)
		os.Exit(2)
	}

	hasDefault := slices.Contains(langs, defaultLang)

	if !hasDefault {
		fmt.Fprintf(
			os.Stderr,
			"default language %q has no language directory with YAML files\n",
			defaultLang,
		)
		os.Exit(2)
	}

	// Generate even if translations are incomplete, so the Go source
	// stays synchronized with the discovered language directories.
	if err := generateLocales(localeFolder, langs); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	if failed {
		os.Exit(1)
	}

	fmt.Println("All translation files contain the baseline keys.")
}

func loadYAML(path string) (*yaml.Node, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parse YAML: %w", err)
	}

	if len(doc.Content) == 0 {
		return nil, fmt.Errorf("empty YAML document")
	}

	return doc.Content[0], nil
}

func mapping(node *yaml.Node) map[string]*yaml.Node {
	result := make(map[string]*yaml.Node)

	if node == nil || node.Kind != yaml.MappingNode {
		return result
	}

	for i := 0; i+1 < len(node.Content); i += 2 {
		key := node.Content[i]
		value := node.Content[i+1]
		result[key.Value] = value
	}

	return result
}

func unwrapLocale(node *yaml.Node) (*yaml.Node, error) {
	if node.Kind != yaml.MappingNode || len(node.Content) != 2 {
		return nil, fmt.Errorf(
			"expected a YAML mapping with one top-level locale key",
		)
	}

	return node.Content[1], nil
}

func missingKeys(base, target *yaml.Node, prefix string) []string {
	baseMap := mapping(base)
	targetMap := mapping(target)

	var missing []string

	// Sort keys so traversal and output are deterministic.
	keys := make([]string, 0, len(baseMap))
	for key := range baseMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		baseValue := baseMap[key]

		path := key
		if prefix != "" {
			path = prefix + "." + key
		}

		targetValue, exists := targetMap[key]
		if !exists {
			missing = append(missing, path)

			// Report every nested key when an entire mapping is missing.
			if baseValue.Kind == yaml.MappingNode {
				missing = append(
					missing,
					missingKeys(baseValue, nil, path)...,
				)
			}
			continue
		}

		if baseValue.Kind == yaml.MappingNode {
			if targetValue.Kind != yaml.MappingNode {
				missing = append(missing, path+" (expected mapping)")

				// Also report the nested keys that cannot be present.
				missing = append(
					missing,
					missingKeys(baseValue, nil, path)...,
				)
				continue
			}

			missing = append(
				missing,
				missingKeys(baseValue, targetValue, path)...,
			)
		}
	}

	return missing
}

// discoverLanguages finds language directories containing YAML files.
func discoverLanguages(root string) ([]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	var langs []string

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		lang := entry.Name()

		// This generator expects two-letter lowercase language codes.
		if len(lang) != 2 ||
			lang[0] < 'a' || lang[0] > 'z' ||
			lang[1] < 'a' || lang[1] > 'z' {
			continue
		}

		files, err := filepath.Glob(
			filepath.Join(root, lang, "*.yaml"),
		)
		if err != nil {
			return nil, err
		}

		if len(files) > 0 {
			langs = append(langs, lang)
		}
	}

	sort.Strings(langs)
	return langs, nil
}

// generateLocales writes the entire [localeFolder]/[outFile] file.
func generateLocales(root string, langs []string) error {
	var b bytes.Buffer

	// header
	fmt.Fprintf(&b, "// Code generated by %s; DO NOT EDIT.\n\n", filepath.Base(thisFile))
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "package locales")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, `import "embed"`)
	fmt.Fprintln(&b)

	// language embeds
	for _, lang := range langs {
		fmt.Fprintf(&b, "//go:embed %s\n", lang)
	}
	fmt.Fprintln(&b, "var Content embed.FS")
	fmt.Fprintln(&b)

	// language constants
	fmt.Fprintln(&b, "type Language string")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "const (")
	for _, lang := range langs {
		fmt.Fprintf(
			&b,
			"\t%s Language = %q\n",
			strings.ToUpper(lang),
			lang,
		)
	}
	fmt.Fprintf(
		&b,
		"\tDefault Language = %s\n",
		strings.ToUpper(defaultLang),
	)
	fmt.Fprintln(&b, ")")
	fmt.Fprintln(&b)

	// to string
	fmt.Fprintln(&b, "func (e Language) String() string {")
	fmt.Fprintln(&b, "\treturn string(e)")
	fmt.Fprintln(&b, "}")
	fmt.Fprintln(&b)

	// isValid
	fmt.Fprintln(&b, "func (e Language) IsValid() bool {")
	fmt.Fprintln(&b, "\tswitch e {")
	for _, lang := range langs {
		fmt.Fprintf(&b, "\tcase %s:\n", strings.ToUpper(lang))
		fmt.Fprintln(&b, "\t\treturn true")
	}
	fmt.Fprintln(&b, "\t}")
	fmt.Fprintln(&b, "\treturn false")
	fmt.Fprintln(&b, "}")

	src, err := format.Source(b.Bytes())
	if err != nil {
		return fmt.Errorf("format generated Go code: %w", err)
	}

	output := filepath.Join(root, outFile)
	if err := os.WriteFile(output, src, 0644); err != nil {
		return fmt.Errorf("write %s: %w", output, err)
	}

	fmt.Printf("Generated: %s\n", output)
	return nil
}
