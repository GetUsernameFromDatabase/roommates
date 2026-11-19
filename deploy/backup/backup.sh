#!/bin/sh
# Backs up the database and the RustFS bucket into /backups, see README.md.
# No argument: loop forever. "once": a single run (the loop calls this).
set -eu

KEEP_DAYS="${BACKUP_KEEP_DAYS:-14}"
DUMPS=/backups/db
FILES=/backups/files

log() { echo "$(date '+%F %T') backup: $*"; }

# the rclone remote "rustfs" and pg_dump's PG* settings come from env, see docker-compose.yml

run_once() {
    today=$(date +%F)
    tmp="$DUMPS/roommates-$today.dump.tmp"
    mkdir -p "$DUMPS" "$FILES/current" "$FILES/removed"
    # a failed dump must not leave a file that looks complete
    trap 'rm -f "$tmp"' EXIT

    log "dumping database"
    pg_dump -Fc -f "$tmp"
    mv "$tmp" "$DUMPS/roommates-$today.dump"

    log "syncing bucket $S3_BUCKET"
    rclone sync "rustfs:$S3_BUCKET" "$FILES/current" --backup-dir "$FILES/removed/$today"

    log "pruning old backups"
    # dates sort as strings
    cutoff=$(date -d "-$KEEP_DAYS days" +%F)
    for f in "$DUMPS"/roommates-*.dump; do
        d=${f##*roommates-}
        d=${d%.dump}
        [ "$d" \< "$cutoff" ] && rm -f "$f" && log "pruned $f"
    done
    for d in "$FILES"/removed/*; do
        [ -d "$d" ] || continue
        [ "${d##*/}" \< "$cutoff" ] && rm -rf "$d" && log "pruned $d"
    done
    echo "$today" > /backups/last-ok
    log "done"
}

# mails BACKUP_ALERT_EMAIL through the SMTP_* settings (same as the app); args: subject, body.
# shortcut: STARTTLS on the given port only (no implicit-TLS smtps on 465), upgrade if a provider needs it
send_mail() {
    case "${SMTP_HOST:-mailpit}" in
        localhost|127.0.0.1|::1|mailpit) ssl="" port=1025 ;;
        *) ssl="--ssl-reqd" port=587 ;;
    esac
    # "Roommates <a@b>" -> a@b
    from=$(printf '%s' "${SMTP_FROM:-}" | sed 's/.*<\(.*\)>.*/\1/')
    # --crlf turns the LF line ends into the CRLF SMTP needs
    # shellcheck disable=SC2086 # $ssl is empty or one flag
    printf 'From: %s\nTo: %s\nSubject: [%s] %s\nDate: %s\n\n%s\n' \
        "${SMTP_FROM:-}" "$BACKUP_ALERT_EMAIL" "$server" "$1" "$(date -R)" "$2" |
        curl -sS --fail --crlf $ssl \
            ${SMTP_USERNAME:+--user "$SMTP_USERNAME:${SMTP_PASSWORD:-}"} \
            --url "smtp://${SMTP_HOST:-mailpit}:${SMTP_PORT:-$port}" \
            --mail-from "$from" --mail-rcpt "$BACKUP_ALERT_EMAIL" --upload-file -
}

# the APP_BASE_URL host names the server in the subject, a container can't see the host's name
server=$(printf '%s' "${APP_BASE_URL:-}" | sed -e 's|^[a-z]*://||' -e 's|[:/].*||')
server=${server:-roommates}

# backup state for the alert body
status() {
    last=$(cat /backups/last-ok 2>/dev/null || echo never)
    echo "Last good backup: $last"
    if [ "$last" != never ]; then
        echo "Days since: $((($(date +%s) - $(date -d "$last" +%s)) / 86400))"
    fi
    newest=$(ls -t "$DUMPS"/roommates-*.dump 2>/dev/null | head -n 1)
    if [ -n "$newest" ]; then
        echo "Newest dump: ${newest##*/} ($(du -h "$newest" | cut -f1))"
    else
        echo "Newest dump: none"
    fi
    echo "Free on the backup disk: $(df -h /backups | awk 'NR==2 {print $4 " of " $2}')"
}

NEXT_STEPS='Next steps, on the server from the repo folder:
  docker compose logs --tail 100 backup     (full output)
  docker compose run --rm backup once       (run one backup by hand)
  deploy/backup/README.md                   (how it works, restore)'

# args: what happened, details; logs, and mails at most once a day while the problem lasts.
# A failing mail never stops the loop.
alert() {
    log "ALERT: $1"
    [ -n "${BACKUP_ALERT_EMAIL:-}" ] || return 0
    today=$(date +%F)
    [ "$(cat /backups/last-alert 2>/dev/null)" != "$today" ] || return 0
    if send_mail "Roommates backup failed" "$(printf '%s\n\n%s\n\n%s\n\n%s' "$1" "$2" "$(status)" "$NEXT_STEPS")"; then
        echo "$today" > /backups/last-alert
    else
        log "alert email failed"
    fi
    return 0
}

# one "works again" mail after alerts were sent; last-alert exists only while they were
recovered() {
    [ -f /backups/last-alert ] || return 0
    # alerts were turned off since, nobody to tell
    [ -n "${BACKUP_ALERT_EMAIL:-}" ] || { rm -f /backups/last-alert; return 0; }
    if send_mail "Roommates backup works again" "$(printf 'Backups succeed again.\n\n%s' "$(status)")"; then
        rm -f /backups/last-alert
    else
        log "recovery email failed, will retry in an hour"
    fi
    return 0
}

if [ "${1:-}" = once ]; then
    run_once
    exit 0
fi

# docker compose run --rm backup test-alert
if [ "${1:-}" = test-alert ]; then
    [ -n "${BACKUP_ALERT_EMAIL:-}" ] || { log "BACKUP_ALERT_EMAIL is not set"; exit 1; }
    send_mail "Roommates backup test" "$(printf 'Alert emails from the backup container work. A real alert looks like this:\n\n%s\n\n%s' "$(status)" "$NEXT_STEPS")"
    log "test alert sent to $BACKUP_ALERT_EMAIL"
    exit 0
fi

# every hour: back up unless today's run fully succeeded (dump and files), so a
# machine that was off at night still gets its backup, and restarts don't repeat it
runlog=/tmp/backup-run.log
while true; do
    last=$(cat /backups/last-ok 2>/dev/null || true)
    if [ "$last" != "$(date +%F)" ]; then
        # separate process, since set -e is off inside a condition; output kept for the alert
        ok=yes
        sh "$0" once > "$runlog" 2>&1 || ok=no
        cat "$runlog"
        if [ "$ok" = no ]; then
            log "FAILED, will retry in an hour"
            # the last "backup: <step>" line names the step that broke
            step=$(grep ' backup: ' "$runlog" | tail -n 1 | sed 's/.* backup: //')
            alert "The backup failed while ${step:-starting}." "$(printf 'Last lines of that run:\n%s' "$(tail -n 20 "$runlog")")"
        fi
    fi
    # still failing after 2 days; a missing last-ok only counts once a dump exists
    last=$(cat /backups/last-ok 2>/dev/null || true)
    if [ "$last" = "$(date +%F)" ]; then
        recovered
    elif [ -n "$last" ]; then
        [ "$last" \< "$(date -d '-2 days' +%F)" ] && alert "No successful backup since $last."
    elif ls "$DUMPS"/roommates-*.dump > /dev/null 2>&1; then
        alert "No successful backup recorded (last-ok is missing)."
    fi
    sleep 3600
done
