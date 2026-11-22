package pages

import "github.com/a-h/templ"

var shoppingReuseScript = templ.NewOnceHandle()

// shoppingBuyAgainShown is how many of the old items show as buttons, all of them are suggested while typing
const shoppingBuyAgainShown = 8

// padInputClass writes on the pad's line, the line shows the focus (.pad-add)
const padInputClass = "h-11 rounded-none border-0 bg-transparent px-0 text-lg shadow-none focus-visible:ring-0 dark:bg-transparent md:text-lg"

const boughtLinkClass = "inline-block py-2 underline underline-offset-4 hover:text-foreground"
