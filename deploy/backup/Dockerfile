# same tag as deploy/postgres/Dockerfile, so pg_dump matches the server
FROM postgres:18.7

RUN apt-get update \
    && apt-get install -y --no-install-recommends rclone curl ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY backup.sh /usr/local/bin/backup.sh

ENTRYPOINT ["sh", "/usr/local/bin/backup.sh"]
