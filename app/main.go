package main

import (
	"context"
	"fmt"
	"roommates/controller"
	"roommates/db"
	"roommates/logger"
	"roommates/rdb"
	"roommates/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "roommates/docs"
)

var log = logger.Main

// https://github.com/swaggo/swag?tab=readme-ov-file#declarative-comments-format
// @title           Roommates API
// Will not be changed as versioning info is not used by anything else, at least tmk
// @version         1.0

// @contact.name   API Support
// @contact.email  ryan.murulo@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

//	@securityDefinitions.apikey  ApiKeyAuth
//	@in                          header
//	@name                        Authorization
//	@description  Used for authentication of most of the access points

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	logger.InitLogger(gin.Mode() == gin.DebugMode)
	serverAddr := utils.MustGetEnv("SERVER_ADDR")
	dbURL := utils.GetDatabaseURL()

	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		panic(fmt.Errorf("unable to create connection pool: %w", err))
	}
	defer dbpool.Close()

	// ensure connection
	if err = dbpool.Ping(ctx); err != nil {
		panic(fmt.Errorf("unable to ping database: %w", err))
	}

	migrationDir := "./db/migrations"
	db.MigrateToLatest(dbpool, migrationDir)
	// to reset database should it be required
	// db.NewMigrations(dbpool, migrationDir).MigrateTo(0)

	redisHandler := rdb.New()
	controllers := controller.New(dbpool, redisHandler)

	e := InitGinEngine(controllers)
	e.RunTLS(serverAddr, "./certificates/server.pem", "./certificates/server.key")
}
