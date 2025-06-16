package database

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func ConnectDB(ctx context.Context, connStr string) *sql.DB {
	// connect to db
	log.Println("connecting to db")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	log.Println("db connected")
	err = db.PingContext(ctx)
	if err != nil {
		panic(err)
	}
	log.Println("db pinged")
	return db
}
