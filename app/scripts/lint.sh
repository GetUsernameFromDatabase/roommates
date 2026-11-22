#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

cd $SCRIPT_DIR/..

status=0
# our .templ files hold only templ components: templ generate copies any Go declared in them
# into the _templ.go, so it would exist twice, and it belongs in the .go file beside them.
# web/components/*/ is left out, those are shadcn-templ components kept as upstream ships them
if grep -nE '^(func|var|const|type) ' web/pages/*.templ web/components/*.templ; then
	echo "move these Go declarations to a .go file" >&2
	status=1
fi

staticcheck -checks "all,-ST1000" ./...
golangci-lint run || status=1
exit $status
