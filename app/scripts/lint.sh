#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

cd $SCRIPT_DIR/..

# a tool failing mid-pipe (tailwindcss missing) must fail the lint, not read as "nothing found"
set -o pipefail
status=0
# our .templ files hold only templ components: templ generate copies any Go declared in them
# into the _templ.go, so it would exist twice, and it belongs in the .go file beside them.
# web/components/*/ is left out, those are shadcn-templ components kept as upstream ships them
if grep -nE '^(func|var|const|type) ' web/pages/*.templ web/components/*.templ; then
	echo "move these Go declarations to a .go file" >&2
	status=1
fi

# Tailwind renames classes between versions (break-words → wrap-break-word); canonicalize gives
# each word in a "", '' or `` string its current name, and non-class words come back unchanged.
# Generated files and tests (which feed CN odd classes on purpose) are left out
tw_files=$(grep -rLE --include='*.go' --include='*.templ' --include='*.js' --exclude='*_templ.go' --exclude='*_test.go' \
	'^// Code generated .* DO NOT EDIT\.$' internal web/pages web/components/*.templ web/components/*.go web/assets/public/*.js)
stale=$(grep -ohE "\"[^\"]*\"|'[^']*'|\`[^\`]*\`" $tw_files |
	tr -d "\"'\`" | tr ' ' '\n' | sort -u |
	tailwindcss canonicalize --stream --css web/assets/globals.css --format jsonl |
	sed -nE 's/^\{"input":"(.*)","output":"(.*)","changed":true\}$/\1 → \2/p') || status=1
if [ -n "$stale" ]; then
	echo "$stale" >&2
	echo "rename these Tailwind classes (old → new)" >&2
	status=1
fi

# a locale key nothing uses is dead text every translation still has to carry
used=$(grep -rhoE --include='*.go' --include='*.templ' --exclude='locales-keys.go' --exclude='*_templ.go' --exclude='*_test.go' \
	'\bLK_\w+' internal web cmd tools | sort -u)
unused=$(grep -oE '^\s+LK_\w+' internal/locales/locales-keys.go | tr -d '[:blank:]' | sort -u | grep -vxF "$used")
if [ -n "$unused" ]; then
	echo "$unused" >&2
	echo "remove these from the locale yaml files, or use them" >&2
	status=1
fi

staticcheck -checks "all,-ST1000" ./... || status=1
golangci-lint run || status=1
exit $status
