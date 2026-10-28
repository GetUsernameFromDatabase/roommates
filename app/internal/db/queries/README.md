# Queries

SQL queries that [sqlc](https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html#schema-and-queries) turns into go code in [dbqueries](../dbqueries), one `<name>.sql.go` per file here

Useful links

- annotations (like `:one`): https://docs.sqlc.dev/en/latest/reference/query-annotations.html#many
- named parameters: https://docs.sqlc.dev/en/latest/howto/named_parameters.html

> [!NOTE]
> a comment right above `-- name:` becomes the doc comment of the generated function
