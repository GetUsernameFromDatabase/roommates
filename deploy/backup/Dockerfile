# same major as the db service, so pg_dump matches the server
FROM postgres:18

RUN apt-get update \
    && apt-get install -y --no-install-recommends rclone \
    && rm -rf /var/lib/apt/lists/*

COPY backup.sh /usr/local/bin/backup.sh

ENTRYPOINT ["sh", "/usr/local/bin/backup.sh"]
