package pages

import "github.com/a-h/templ"

var expenseSplitScript = templ.NewOnceHandle()

// tillFieldClass is a line to write on, the input's ring stays as the focus mark because
// a read only share has no border to show it
const tillFieldClass = "h-11 rounded-none border-0 border-b-2 border-dotted border-foreground/40 bg-transparent px-1 shadow-none dark:bg-transparent"

var datePickerScript = templ.NewOnceHandle()
