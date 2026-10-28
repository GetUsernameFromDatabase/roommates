#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

templ fmt $SCRIPT_DIR/..
gofmt -s -w $SCRIPT_DIR/..
