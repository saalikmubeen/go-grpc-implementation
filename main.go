package main

import (
	"database/sql"
	"log"

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

	store := generated_db.NewStore(conn)
	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
