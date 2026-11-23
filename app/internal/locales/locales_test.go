package locales

import (
	"fmt"
	"io/fs"
	"regexp"
	"slices"
	"testing"

	"gopkg.in/yaml.v3"
)

// %% is matched first so it is not read as a verb, then fmt verbs (with flags, width, index), then %{name}
var rgxPlaceholder = regexp.MustCompile(`%%|%[-+# 0]*(?:\[\d+\])?\d*(?:\.\d+)?[a-zA-Z]|%\{\w+\}`)

// loadFlat reads <lang>/<lang>.yaml as dot keys, under its top-level language key
func loadFlat(t *testing.T, lang Language) map[string]string {
	t.Helper()
	data, err := fs.ReadFile(Content, lang.String()+"/"+lang.String()+".yaml")
	if err != nil {
		t.Fatalf("%s: %v", lang, err)
	}
	var doc map[string]any
	if err := yaml.Unmarshal(data, &doc); err != nil {
		t.Fatalf("%s: %v", lang, err)
	}
	out := map[string]string{}
	var flatten func(prefix string, node any)
	flatten = func(prefix string, node any) {
		if m, ok := node.(map[string]any); ok {
			// a plural group is one key, compared by its "other" form, so an optional "zero" in one language is fine
			if isPlural(m) {
				out[prefix[1:]] = fmt.Sprint(m["other"])
				return
			}
			for k, v := range m {
				flatten(prefix+"."+k, v)
			}
			return
		}
		out[prefix[1:]] = fmt.Sprint(node)
	}
	flatten("", doc[lang.String()])
	return out
}

// isPlural is the rule tools/generate-locales-keys uses: "other" plus only the forms ctxi18n picks from
func isPlural(m map[string]any) bool {
	if _, ok := m["other"]; !ok {
		return false
	}
	for key := range m {
		if key != "zero" && key != "one" && key != "other" {
			return false
		}
	}
	return true
}

func placeholders(s string) []string {
	found := slices.DeleteFunc(rgxPlaceholder.FindAllString(s, -1), func(p string) bool { return p == "%%" })
	slices.Sort(found)
	return found
}

func TestLocalesMatchDefault(t *testing.T) {
	base := loadFlat(t, Default)
	for _, lang := range All {
		t.Run(lang.String(), func(t *testing.T) {
			got := loadFlat(t, lang)
			for key, want := range base {
				value, ok := got[key]
				if !ok {
					t.Errorf("%s: missing key %s", lang, key)
					continue
				}
				if p, wp := placeholders(value), placeholders(want); !slices.Equal(p, wp) {
					t.Errorf("%s: %s has placeholders %v, %s has %v", lang, key, p, Default, wp)
				}
			}
			for key := range got {
				if _, ok := base[key]; !ok {
					t.Errorf("%s: extra key %s, not in %s", lang, key, Default)
				}
			}
			if got["language.name"] == "" {
				t.Errorf("%s: language.name is empty", lang)
			}
			if _, err := fs.Stat(Content, lang.String()+"/flag.svg"); err != nil {
				t.Errorf("%s: flag.svg: %v", lang, err)
			}
		})
	}
}

func TestIsPlural(t *testing.T) {
	for _, c := range []struct {
		m    map[string]any
		want bool
	}{
		{map[string]any{"one": "a", "other": "b"}, true},
		{map[string]any{"zero": "a", "one": "b", "other": "c"}, true},
		{map[string]any{"cash": "a", "card": "b", "other": "c"}, false},
		{map[string]any{"one": "a"}, false},
	} {
		if got := isPlural(c.m); got != c.want {
			t.Errorf("isPlural(%v) = %v, want %v", c.m, got, c.want)
		}
	}
}

func TestPlaceholders(t *testing.T) {
	got := placeholders("%[2]s 100%% %d %{count} %-5.2f %v")
	want := []string{"%-5.2f", "%[2]s", "%d", "%v", "%{count}"}
	if !slices.Equal(got, want) {
		t.Errorf("placeholders = %v, want %v", got, want)
	}
}
