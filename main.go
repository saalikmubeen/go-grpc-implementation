package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"

	// Our Migration Source is the file system (file://db/migrations)
	// where our migration files are stored.
	// Other sources include: gitlab, s3, github, etc.
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/lib/pq"

	"github.com/saalikmubeen/backend-masterclass-go/api"
	generated_db "github.com/saalikmubeen/backend-masterclass-go/db/sqlc"
	"github.com/saalikmubeen/backend-masterclass-go/utils"
)

func main() {

	config, err := utils.LoadConfig(".")

	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// fmt.Println(config)

	conn, err := sql.Open(config.DBDriver, config.DBURI)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	// Another way of running database migrations
	// (directly from the GO code)
	runDBMigration(config.DBMigrationsURL, config.DBURI)

	store := generated_db.NewStore(conn)

	startGinServer(config, store)
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create new migrate instance:", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up:", err)
	}

	log.Println("db migrated successfully")
}

func startGinServer(config utils.Config, store generated_db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	fmt.Printf("Starting server at %s....!\n", config.HTTPServerAddress)
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
