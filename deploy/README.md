<!-- omit in toc -->
# Deploy

Running Roommates on a server with Docker Compose behind a Cloudflare Tunnel. Run every command from the repo root.

- [At a glance](#at-a-glance)
- [Details](#details)
  - [1. Required settings](#1-required-settings)
  - [2. First deploy](#2-first-deploy)
  - [3. Upgrade](#3-upgrade)
  - [4. Rollback](#4-rollback)
  - [5. Postgres major upgrade](#5-postgres-major-upgrade)
  - [6. Good to know](#6-good-to-know)
  - [7. One-time setup](#7-one-time-setup)
    - [7.1 Domain on Cloudflare](#71-domain-on-cloudflare)
    - [7.2 Tunnel](#72-tunnel)
    - [7.3 App address](#73-app-address)
    - [7.4 Web push keys](#74-web-push-keys)
    - [7.5 Email (SMTP)](#75-email-smtp)
    - [7.6 RustFS key for the app](#76-rustfs-key-for-the-app)
    - [7.7 Backups](#77-backups)

## At a glance

- **Settings:** fill `.env` as in [1](#1-required-settings); how to get each value (domain, tunnel, push keys, mail, RustFS key, backup alert) is the [one-time setup](#7-one-time-setup). Costs: a domain; Cloudflare's free plan and a free-tier mail provider are enough for one flat.
- **First deploy:** `docker compose up -d --build`, wait for `healthy`, open your hostname. See [2](#2-first-deploy).
- **Upgrade:** back up first (`docker compose run --rm backup once`), then `git pull` and `docker compose up -d --build`. See [3](#3-upgrade).
- **Rollback:** previous commit + rebuild; if the new version migrated the database, restore the pre-upgrade dump. See [4](#4-rollback).
- **Postgres major upgrade:** dump, new empty volume, restore. Never just bump the image tag. See [5](#5-postgres-major-upgrade).
- **Never** forward or publish port 8000: it is plain HTTP that trusts Cloudflare's header. See [6](#6-good-to-know).

## Details

### 1. Required settings

All in `.env` at the repo root, starting from [.env.example](../.env.example). How to get each value: [7. One-time setup](#7-one-time-setup). Nobody else reads `.env`, so you type the values in.

| Variable | Value |
| --- | --- |
| `COMPOSE_PROFILES` | `prod` (starts `cloudflared`, leaves out the dev `mailpit`) |
| `TRUSTED_PLATFORM` | `cloudflare` (TLS ends at Cloudflare; the app speaks plain HTTP). With it set, the app refuses to start on the development defaults below |
| `APP_BASE_URL` | `https://<your hostname>`, used in email links; a localhost address is refused |
| `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `SMTP_FROM` | a provider with STARTTLS (port 587); a non-local host without TLS fails, `mailpit` is refused |
| `VAPID_PUBLIC_KEY`, `VAPID_PRIVATE_KEY` | required: the app refuses to start without them. Generate once, keep forever |
| `CLOUDFLARE_TUNNEL_TOKEN` | the tunnel's token; its public hostname targets `http://app:8080` |
| `APP_S3_ACCESS_KEY`, `APP_S3_SECRET_KEY` | required: the app's own RustFS key (bucket only), see [7.6](#76-rustfs-key-for-the-app); the app never uses the root keys |
| `BACKUP_ALERT_EMAIL` | where backup failure alerts go, see [deploy/backup/README.md](./backup/README.md) |
| `DATABASE_MAX_CONNS` | optional, connection pool size, default 10 |

Plus the values every setup needs (`POSTGRES_*`, `SESSION_*`, `S3_*` root keys), as in `.env.example`. An empty value (`KEY=`) counts as unset. The app container gets only the keys listed under `app.environment` in `docker-compose.yml`, not the whole `.env`.

### 2. First deploy

```sh
docker compose up -d --build
docker compose ps
```

Wait until `db`, `rustfs` and `app` show `healthy` (`cloudflared` starts once the app is healthy, `backup` has no healthcheck). The database is migrated on start. Then open your hostname in a browser, and send a test alert: `docker compose run --rm backup test-alert`.

### 3. Upgrade

```sh
docker compose run --rm backup once
git pull
docker compose up -d --build
docker compose ps
```

The first line makes a fresh dump in `BACKUP_DIR` (it replaces today's dump if there already is one). Migrations run when the app starts, under a Postgres advisory lock, so two starting containers never migrate at the same time.

### 4. Rollback

```sh
git checkout <previous commit>
docker compose up -d --build
```

Migrations are **not** rolled back automatically. If the new version migrated the database, the old code may not work with it: stop the app and restore the dump made before the upgrade, see [Restore](./backup/README.md#restore).

### 5. Postgres major upgrade

A new major version cannot read the old data folder. Never just bump the tag on the existing `db_data` volume.

1. Back up: `docker compose run --rm backup once`, and note the dump's date.
2. Change the major version in **both** `deploy/postgres/Dockerfile` and `deploy/backup/Dockerfile` (they share a tag so `pg_dump` matches the server).
3. Stop everything: `docker compose down`. Keep the old volume until the new one works: rename it in `docker-compose.yml`, both under the top-level `volumes:` and in the `db` service (for example `db_data` → `db_data_19`), so the new server starts on an empty volume.
4. `docker compose up -d --build db`, wait until it is healthy.
5. Restore the dump into it (the database is created empty on first start):

   ```sh
   docker compose run --rm --entrypoint sh backup -c \
     'pg_restore --no-owner -d "$PGDATABASE" /backups/db/roommates-YYYY-MM-DD.dump'
   ```

6. `docker compose up -d --build`, check the app. Once happy, delete the old volume with `docker volume rm`.

### 6. Good to know

- **Port 8000:** compose publishes the app on `127.0.0.1:8000` (loopback only) for local checks. Never forward it on the router or change it to a public address: with `TRUSTED_PLATFORM` the app speaks plain HTTP and trusts the client-IP header Cloudflare sets, which anyone reaching the port directly could fake. Only the tunnel should reach the app.
- **Logs:** `docker compose logs <service>` (add `-f` to follow). Each service's log is rotated at 10 MB, 3 files kept.
- **Certificates:** only needed without `TRUSTED_PLATFORM`. The app image runs as a non-root user (uid 1000), so on a Linux host the bind-mounted `./app/certificates` must be readable by that user.
- **Backups:** what runs, the alert and how to restore: [deploy/backup/README.md](./backup/README.md).
- **Second copy on one machine:** give it its own folder name or `COMPOSE_PROJECT_NAME`, otherwise both share containers and volumes. To run both at once, give it other host ports too: `DATABASE_PORT`, `APP_PORT`, `RUSTFS_PORT`, `RUSTFS_CONSOLE_PORT`, `MAILPIT_SMTP_PORT`, `MAILPIT_WEB_PORT`, see `.env.example`. `task dev`'s server port (8090) and debugger (2345) are fixed, so only one `task dev` runs at a time.

### 7. One-time setup

Do these before the first deploy. Every value goes in `.env` at the repo root.

#### 7.1 Domain on Cloudflare

1. Create a Cloudflare account and add your domain (or buy one there).
2. If the domain was bought elsewhere, set its nameservers to the two Cloudflare shows you, and wait until Cloudflare marks the domain active.

The free plan lets uploads up to 100 MB through; the app's own limit is 25 MB per file.

#### 7.2 Tunnel

1. Cloudflare dashboard → Zero Trust → Networks → Tunnels → Create a tunnel → type **Cloudflared**.
2. Name it (for example `roommates`). On the install screen pick Docker and copy only the token, the long string after `--token`.
3. In `.env`: `CLOUDFLARE_TUNNEL_TOKEN=<the token>`.
4. In the tunnel's **Public hostname** tab add one: your subdomain (for example `flat.example.com`), service type **HTTP**, URL `app:8080`. The app always listens on 8080 inside its container.

You do not open any port on your router.

#### 7.3 App address

```txt
COMPOSE_PROFILES=prod
APP_BASE_URL=https://flat.example.com
TRUSTED_PLATFORM=cloudflare
```

`COMPOSE_PROFILES=prod` starts the tunnel (`cloudflared`) and leaves out the dev mail catcher. `APP_BASE_URL` is what links in emails point at.

#### 7.4 Web push keys

From `app/`:

```sh
go run ./tools/vapid-keys
```

Put the two values in `.env` as `VAPID_PUBLIC_KEY` and `VAPID_PRIVATE_KEY`. Generate them once and keep them: new keys silently break every phone that already turned notifications on. In production the app refuses to start without them, on purpose.

#### 7.5 Email (SMTP)

1. Make an account at a mail provider that offers SMTP with STARTTLS on port 587. `.env.example` uses Resend as its example.
2. Verify your domain there. The provider gives you DNS records (SPF, DKIM): add them in Cloudflare → DNS, with the proxy (orange cloud) **off** for those records.
3. In `.env`:

   ```txt
   SMTP_HOST=<provider host>
   SMTP_PORT=587
   SMTP_USERNAME=<from the provider>
   SMTP_PASSWORD=<from the provider>
   SMTP_FROM=Roommates <noreply@your-domain>
   ```

Test: register a new account on the real address and check that the verification email arrives and its link opens your domain.

#### 7.6 RustFS key for the app

1. `docker compose up -d --wait rustfs`, then `sh deploy/rustfs/app-key.sh`.
2. Paste the two printed lines (`APP_S3_ACCESS_KEY`, `APP_S3_SECRET_KEY`) into `.env`, then `docker compose up -d`.

The script creates the RustFS user `roommates-app` for the bucket `S3_BUCKET` of `.env`. Its policy, near the top of the script, allows only what the app does: list, read, write and delete in that bucket, and create the bucket if it is missing. It allows no other bucket and nothing admin. It always works on the RustFS of the checkout it lives in, whatever folder you run it from.

- **Changed the policy:** run it again. It re-applies the policy and keeps the secret.
- **Secret leaked:** `sh deploy/rustfs/app-key.sh rotate`, paste the new lines, `docker compose up -d`. The old secret stops working at once.
- **New RustFS volume** (new server, restore onto an empty volume): the user lives in the volume, not in the backup, so run the script again and paste.

The app refuses to start without these. The root keys (`S3_ACCESS_KEY`, `S3_SECRET_KEY`) stay with RustFS and the backup service.

#### 7.7 Backups

1. In `.env`: `BACKUP_ALERT_EMAIL=<your address>`. Then check it arrives: `docker compose run --rm backup test-alert`.
2. Recommended: point `BACKUP_DIR` at a different drive than the one holding the data.
3. Read [deploy/backup/README.md](./backup/README.md) once, so the restore steps are familiar before you need them.
