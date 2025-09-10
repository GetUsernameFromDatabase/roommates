#!/bin/bash
# for generating code
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/a-h/templ/cmd/templ@latest
go install github.com/swaggo/swag/cmd/swag@latest

# only for dev, allows live reload
go install github.com/air-verse/air@latest