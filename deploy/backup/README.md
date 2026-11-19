<!-- omit in toc -->
# Backup

The `backup` compose service keeps a local copy of the database and the files.

Every hour it checks `last-ok` (the date of the last fully successful run), and if it is not today it:

1. `pg_dump -Fc` of the `db` service to `/backups/db/roommates-YYYY-MM-DD.dump` (written as `.tmp`, then renamed).
2. `rclone sync` of the bucket (`S3_BUCKET`, over the S3 API of the `rustfs` service) to `/backups/files/current`. A file deleted or replaced in the bucket moves from `current` to `/backups/files/removed/YYYY-MM-DD`, dated by the run that noticed it, so an older dump still finds its files.
3. Deletes dumps and `removed/<date>` folders older than `BACKUP_KEEP_DAYS`.
4. Writes today's date to `/backups/last-ok`.

A failed run is logged and retried an hour later: `docker compose logs backup`.

Run one backup now (for example before an upgrade): `docker compose run --rm backup once`.

- [Alert](#alert)
- [Config (`.env`)](#config-env)
- [Restore](#restore)
  - [Newest dump](#newest-dump)
  - [An older dump](#an-older-dump)
  - [Try it on scratch copies first](#try-it-on-scratch-copies-first)

## Alert

With `BACKUP_ALERT_EMAIL` set, the service emails that address, at most once a day, when a run fails or when `last-ok` is older than 2 days, and once more when a backup succeeds again. The subject names the server (the host of `APP_BASE_URL`). A failure mail says which step broke, quotes the last 20 lines of that run, and lists the last good backup, the newest dump and its size, free space on the backup disk, and the commands to look further. It sends through the same `SMTP_*` settings as the app (STARTTLS unless the host is local). Without it, problems only go to the log.

Test the mail once after setting it up:

```sh
docker compose run --rm backup test-alert
```

## Config (`.env`)

- `BACKUP_DIR`: host folder mounted at `/backups`, default `./backups`.
- `BACKUP_KEEP_DAYS`: days to keep, default 14.
- `BACKUP_ALERT_EMAIL`: where failure alerts go, unset = log only.
- DB, S3 and SMTP settings are the ones the app uses (`POSTGRES_*`, `S3_*` root keys, `SMTP_*`).

shortcut: same disk, point BACKUP_DIR at another drive or add off-site when affordable

## Restore

Run from the repo root. The `rustfs:` remote is configured by the `RCLONE_CONFIG_RUSTFS_*` variables of the `backup` service in `docker-compose.yml`. Inside `-c '...'` refer to variables by name (`$PGDATABASE`, `$S3_BUCKET`) so their values are never typed or printed. Stop the app first: `docker compose stop app`, and `docker compose start app` when done.

### Newest dump

Database, into the running `db` service:

```sh
docker compose run --rm --entrypoint sh backup -c \
  'pg_restore --clean --if-exists --no-owner -d "$PGDATABASE" /backups/db/roommates-2026-01-31.dump'
```

Files, copied back into the bucket:

```sh
docker compose run --rm --entrypoint sh backup -c \
  'rclone copy /backups/files/current rustfs:"$S3_BUCKET"'
```

### An older dump

`current` matches the newest run. Files that the older dump still points at but that were deleted or replaced later sit in the `removed/<date>` folders of the runs **after** the dump date (the folder of the dump date itself holds files that were already gone before that dump).

Copy `current` first, then every `removed/<date>` folder later than the dump date, **newest date first**, overwriting each time. A file replaced several times then ends up with the version from the oldest of those folders, which is the one the dump knew. Set `D` to the dump date:

```sh
docker compose run --rm --entrypoint sh backup -c '
  D=2026-01-31
  rclone copy --ignore-times /backups/files/current rustfs:"$S3_BUCKET"
  for r in $(ls -r /backups/files/removed); do
    if [ "$r" \> "$D" ]; then rclone copy --ignore-times "/backups/files/removed/$r" rustfs:"$S3_BUCKET"; fi
  done'
```

`--ignore-times` makes every copy overwrite, so the order alone decides which version wins. Files uploaded after the dump stay in the bucket; the restored database does not reference them, they only take space.

Then restore that dump with the database command above.

### Try it on scratch copies first

This is the drill that was run on 2026-11-19 (row counts and files matched). It never writes to the real database or bucket.

Database: restore into a new database, compare row counts per table with the live one, drop it:

```sh
docker compose run --rm -T --entrypoint sh backup -c '
  createdb restore_drill
  pg_restore --no-owner -d restore_drill /backups/db/roommates-2026-01-31.dump
  Q="select table_name, (xpath(\$\$/row/c/text()\$\$, query_to_xml(format(\$\$select count(*) as c from public.%I\$\$, table_name), false, true, \$\$\$\$)))[1]::text from information_schema.tables where table_schema = \$\$public\$\$ and table_type = \$\$BASE TABLE\$\$ order by 1"
  psql -X -At -c "$Q" > /tmp/live.txt
  psql -X -At -d restore_drill -c "$Q" > /tmp/drill.txt
  diff /tmp/live.txt /tmp/drill.txt && echo "row counts match"
  dropdb restore_drill'
```

Rows written after the dump (new sessions, messages) show up as differences; on a quiet system there are none.

Files: copy into a new bucket, compare with the real one (read-only), delete it:

```sh
docker compose run --rm --entrypoint sh backup -c '
  rclone mkdir rustfs:restore-drill
  rclone copy /backups/files/current rustfs:restore-drill
  rclone check rustfs:restore-drill rustfs:"$S3_BUCKET"
  rclone purge rustfs:restore-drill'
```

`rclone check` ends with `0 differences found` when the copy is complete. For an older dump, run the loop above with `rustfs:restore-drill` in place of `rustfs:"$S3_BUCKET"` and compare it with the `files` table of the scratch database instead.
