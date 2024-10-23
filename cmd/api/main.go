package main

import (
	"github.com/Kser1x/social/internal/db"
	"github.com/Kser1x/social/internal/env"
	"github.com/Kser1x/social/internal/store"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	//создаём переменную с объектом конфигурации
	cfg := config{
		addr: env.GetSting("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetSting("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetSting("DB_MAX_IDLE_TIME", "1m"),
		},
	}

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Panicf("Ошибка подключения %v", err)
	}

	defer db.Close()
	log.Println("соединение с БД установлено")

	storage := store.NewStorage(db)
	//создаётся переменная с ссылкой на объект
	app := &application{
		config:  cfg,
		storage: storage,
	}

	// При возникновении ошибки выводит лог фатал
	mux := app.mount()
	log.Fatal(app.run(mux))
}
