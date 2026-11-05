# Backup

The `backup` compose service keeps a local copy of the database and the files.

Every hour it checks `last-ok` (the date of the last fully successful run), and if it is not today it:

1. `pg_dump -Fc` of the `db` service to `/backups/db/roommates-YYYY-MM-DD.dump` (written as `.tmp`, then renamed).
2. `rclone sync` of the bucket (`S3_BUCKET`, over the S3 API of the `rustfs` service) to `/backups/files/current`. Files deleted or replaced since the last run move to `/backups/files/removed/YYYY-MM-DD`, so an older dump still finds its files.
3. Deletes dumps and `removed/<date>` folders older than `BACKUP_KEEP_DAYS`.

A failed run is logged and retried an hour later: `docker compose logs backup`. Nothing else alerts you.

## Config (`.env`)

- `BACKUP_DIR`: host folder mounted at `/backups`, default `./backups`.
- `BACKUP_KEEP_DAYS`: days to keep, default 14.
- DB and S3 settings are the ones the app uses (`POSTGRES_*`, `S3_*`).

shortcut: same disk, point BACKUP_DIR at another drive or add off-site when affordable

## Restore

Run from the repo root. Replace the date. Stop the app first (`docker compose stop app`).

Database, into the running `db` service:

```sh
docker compose run --rm --entrypoint sh backup -c \
  'pg_restore --clean --if-exists --no-owner -d "$PGDATABASE" /backups/db/roommates-2026-01-31.dump'
```

Try it on a throwaway database first by setting `PGDATABASE` to a new, created database.

Files, copied back into the bucket:

```sh
docker compose run --rm --entrypoint sh backup -c \
  'rclone copy /backups/files/current rustfs:"$S3_BUCKET"'
```

Which folder: `current` is the newest state. To get back something deleted after the dump you want, look in `/backups/files/removed/`: each `<date>` folder holds what the run on that date deleted or replaced. Copy the folder from the first date on or after the files were lost:

```sh
docker compose run --rm --entrypoint sh backup -c \
  'rclone copy /backups/files/removed/2026-01-31 rustfs:"$S3_BUCKET"'
```

The `rustfs:` remote is configured by the `RCLONE_CONFIG_RUSTFS_*` variables of the `backup` service in `docker-compose.yml`.
