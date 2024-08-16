package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/taufiqdp/go-simplebank/api"
	db "github.com/taufiqdp/go-simplebank/db/sqlc"
	"github.com/taufiqdp/go-simplebank/utils"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	if err := server.Start(config.ServerAddress); err != nil {
		log.Fatal("cannot start server: ", err)
	}

}
