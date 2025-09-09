#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd $SCRIPT_DIR

# see ".vscode/tasks.json"
sh "./locales/gen-keys/run.sh"
sqlc generate
templ generate
swag init