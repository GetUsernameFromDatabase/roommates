package pages

import "github.com/a-h/templ"

var chatScript = templ.NewOnceHandle()

const chatControlsClass = "shrink-0 gap-0.5 self-center px-0 pointer-events-none opacity-0 " +
	"pointer-coarse:absolute pointer-coarse:-top-4 pointer-coarse:z-10 pointer-coarse:rounded-full pointer-coarse:bg-popover " +
	"pointer-coarse:p-0.5 pointer-coarse:shadow-md pointer-coarse:ring-1 pointer-coarse:ring-foreground/10 " +
	"group-hover/message:pointer-events-auto group-hover/message:opacity-100 " +
	"group-focus-within/message:pointer-events-auto group-focus-within/message:opacity-100 " +
	"group-data-open/message:pointer-events-auto group-data-open/message:opacity-100 " +
	"motion-safe:group-data-open/message:animate-in motion-safe:group-data-open/message:fade-in motion-safe:group-data-open/message:zoom-in-95 motion-safe:group-data-open/message:duration-120"

// chatControlClass is a labelled touch target on touch screens, a round icon with a mouse
const chatControlClass = "min-h-11 rounded-full text-muted-foreground pointer-fine:size-8 pointer-fine:min-h-8 pointer-fine:p-0"

// localTimeScript is local-time.js, a page with times loads it at page level, not in localTime
var localTimeScript = templ.NewOnceHandle()

// utcMinuteLayout is the UTC time shown until local-time.js replaces it
const utcMinuteLayout = "2006-01-02 15:04"
