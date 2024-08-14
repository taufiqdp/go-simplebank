package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/taufiqdp/go-simplebank/api"
	db "github.com/taufiqdp/go-simplebank/db/sqlc"
)

func main() {
	conn, err := sql.Open("postgres", "postgresql://root:pwd@127.0.0.1:5432/simplebank?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err := server.Start("127.0.0.1:3000"); err != nil {
		log.Fatal("cannot start server: ", err)
	}

}
