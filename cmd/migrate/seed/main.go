package main

import (
	"github.com/Kser1x/social/internal/db"
	"github.com/Kser1x/social/internal/env"
	store2 "github.com/Kser1x/social/internal/store"
	"log"
)

func main() {
	addr := env.GetSting("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/social?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	store := store2.NewStorage(conn)

	db.Seed(store)
}
