# Roommates

If you are seeing this then that means I have not finished this project
Enjoy the WIP state of this overengineered project -- doing this for enjoyment and something to add to portfolio

## Future Setup

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

> <https://docs.sqlc.dev/en/latest/overview/install.html>

---

kafka:

- <https://towardsdev.com/how-to-build-a-highly-customizable-and-scalable-kafka-consumer-using-goroutines-4c128b6b9058>
- <https://developer.confluent.io/get-started/go/#build-producer>

will be used mostly for logging actions, possible statistics  
current plan is to just stream to "`/dev/null`" as collecting that data is not important but it's nice to have a bulk of the system ready  
Also to showcase experience with this

---

db:

- <https://docs.sqlc.dev/en/latest/howto/ddl.html>
- <https://github.com/jackc/pgx/wiki/Getting-started-with-pgx#using-a-connection-pool>

---

backend:

- <https://github.com/gin-gonic/examples/blob/master/http2/main.go>
- <https://github.com/swaggo/swag?tab=readme-ov-file#how-to-use-it-with-gin>

- <https://protobuf.dev/getting-started/gotutorial/>
- <https://github.com/protobufjs/protobuf.js/?tab=readme-ov-file#using-proto-files>

- <https://github.com/minio/minio> FILE STORAGE
- <https://hub.docker.com/r/minio/minio>
- <https://hub.docker.com/r/minio/minio>

- <https://stackoverflow.com/questions/78055110/browser-hot-reload-with-golang-templ-htmx>
- <https://dev.to/siumhossain/how-to-enable-hot-reload-in-your-gin-project-42g4>

---

frontend:

- <https://franken-ui.dev/>
- <https://htmx.org/docs/#parameters>
- <https://templ.guide/>

---

- `go install github.com/air-verse/air@latest`
- `go install github.com/a-h/templ/cmd/templ@latest`
- `go install github.com/swaggo/swag/cmd/swag@latest`
- `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`

> [setup-dev.sh](./app/setup-dev.sh) has them all for easy install

use [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/#summary) with [preferred commit types](https://github.com/conventional-changelog/commitlint/tree/master/%40commitlint/config-conventional#type-enum)
