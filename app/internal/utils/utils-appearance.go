package utils

import (
	"context"
	"slices"
)

// Appearance is the theme and mode a page is rendered in,
// written as data-theme and data-mode on <html>
//
// a theme is a file in web/assets/themes, its colours are light-dark() pairs
type Appearance struct {
	Theme string
	Mode  string
}

// the first of each is the default
var (
	Themes = []string{"home", "plain"}
	Modes  = []string{"system", "light", "dark"}
)

// themeBackgrounds are the light and dark --background of each theme, for <meta name="theme-color">
var themeBackgrounds = map[string][2]string{
	"home":  {"#e9ece4", "#16201b"},
	"plain": {"#ffffff", "#0a0a0a"},
}

// NewAppearance replaces an unknown theme or mode with the default,
// the values come from cookies and forms
func NewAppearance(theme, mode string) Appearance {
	if !slices.Contains(Themes, theme) {
		theme = Themes[0]
	}
	if !slices.Contains(Modes, mode) {
		mode = Modes[0]
	}
	return Appearance{Theme: theme, Mode: mode}
}

// Background is the theme's --background in light and in dark
func (a Appearance) Background() (light, dark string) {
	b := themeBackgrounds[a.Theme]
	return b[0], b[1]
}

type appearanceKey struct{}

func WithAppearance(ctx context.Context, a Appearance) context.Context {
	return context.WithValue(ctx, appearanceKey{}, a)
}

// GetAppearance is the default when no middleware set one
func GetAppearance(ctx context.Context) Appearance {
	if a, ok := ctx.Value(appearanceKey{}).(Appearance); ok {
		return a
	}
	return NewAppearance("", "")
}
