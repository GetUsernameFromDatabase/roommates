#!/usr/bin/env bash
#
# PreToolUse hook (Bash): denies commands that mention a secret file,
# so secrets are not loaded into commands or printed into the transcript.
#
# How it works:
#   - checks the raw hook JSON (no jq needed), so a match anywhere in the
#     command counts, regardless of cd, paths or quoting
#   - quotes and backslashes are removed first, so .e""nv is still caught
#   - ALLOWED names are removed next, then each blocked pattern is tried
#   - the first match denies through hookSpecificOutput.permissionDecision,
#     with that rule's reason
#
# Limits: catches accidental use, not deliberate obfuscation (a name built
# at runtime), and not tools that read the files themselves (task dev,
# docker compose). Run those yourself with `! <command>`.
#
# To protect another file, add a `block PATTERN REASON` line below
# (and any harmless look-alike to ALLOWED).

BLOCKED_PATTERNS=()
BLOCKED_REASONS=()

# block PATTERN REASON
#   PATTERN: extended regular expression (grep -E)
#   REASON:  why it is protected, shown to Claude (any text)
block() {
	BLOCKED_PATTERNS+=("$1")
	BLOCKED_REASONS+=("$2")
}

# --- configuration ------------------------------------------------------------

# literal names that may be mentioned, removed before matching
ALLOWED=(
	'.env.example'
)

ENV_REASON="The .env file holds the app secrets (database password, session keys, API keys). \
For database commands use the container's own env instead: docker compose exec -T db sh -c 'psql -U \"\$POSTGRES_USER\" ...'."

block '\.env([^A-Za-z0-9_]|$)' "$ENV_REASON" # .env, .env.local, prod.env (not .envrc)
block '\.en[?*[]' "$ENV_REASON"              # globs: .en? .en* .en[v]
block '\.e[?*[]' "$ENV_REASON"               # globs: .e?? .e*

# appended to every reason
HOW_TO_PROCEED="Do not read or load it. If the command is needed, ask the user to run it themselves with '! <command>'."

# --- output -------------------------------------------------------------------

# %s is replaced with the JSON-escaped reason
DENY_TEMPLATE='{
  "hookSpecificOutput": {
    "hookEventName": "PreToolUse",
    "permissionDecision": "deny",
    "permissionDecisionReason": "%s"
  }
}'

# escapes backslashes and double quotes so any text fits in a JSON string
json_escape() {
	local s=${1//\\/\\\\}
	printf '%s' "${s//\"/\\\"}"
}

deny() {
	local reason="Blocked: $1 $HOW_TO_PROCEED"
	printf "$DENY_TEMPLATE\n" "$(json_escape "$reason")"
	exit 0
}

# --- check --------------------------------------------------------------------

# hook input JSON, with JSON escaping and shell quoting removed
text=$(tr -d "\"'\\\\")

for name in "${ALLOWED[@]}"; do
	text=${text//"$name"/}
done

for i in "${!BLOCKED_PATTERNS[@]}"; do
	if grep -qE -- "${BLOCKED_PATTERNS[i]}" <<<"$text"; then
		deny "${BLOCKED_REASONS[i]}"
	fi
done

exit 0
