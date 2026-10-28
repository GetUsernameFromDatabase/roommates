package utils

import "testing"

// cases from the comments in shadcn-templ.go
func TestCN(t *testing.T) {
	tests := []struct {
		in   []any
		want string
	}{
		{[]any{"bg-red-500 hover:bg-blue-500", "bg-green-500"}, "hover:bg-blue-500 bg-green-500"},
		// input order is kept, "pb-11 p-2" would lose the bottom padding
		{[]any{"p-2", "pb-11"}, "p-2 pb-11"},
		{[]any{"pb-11", "p-2"}, "p-2"},
		// v4 variable shorthand
		{[]any{"px-(--card-spacing)", "px-4"}, "px-4"},
		{[]any{"px-4", "px-(--card-spacing)"}, "px-(--card-spacing)"},
		// v4 trailing important marker
		{[]any{"group-data-[x]:p-2!", "group-data-[x]:p-4!"}, "group-data-[x]:p-4!"},
		// clsx-style inputs
		{[]any{"p-2", map[string]bool{"m-1": true, "m-2": false}, []string{"", "text-sm"}, nil}, "p-2 m-1 text-sm"},
	}
	for _, tt := range tests {
		if got := CN(tt.in...); got != tt.want {
			t.Errorf("CN(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
