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

if [ "${1:-}" = once ]; then
    run_once
    exit 0
fi

# every hour: back up unless today's run fully succeeded (dump and files), so a
# machine that was off at night still gets its backup, and restarts don't repeat it
while true; do
    if [ "$(cat /backups/last-ok 2>/dev/null)" != "$(date +%F)" ]; then
        # separate process, since set -e is off inside a condition
        sh "$0" once || log "FAILED, will retry in an hour"
    fi
    sleep 3600
done
