package main

import (
	"github.com/Kser1x/social/internal/env"
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
	}
	//создаётся переменная с ссылкой на объект
	app := &application{
		config: cfg,
	}
	// При возникновении ошибки выводит лог фатал
	mux := app.mount()
	log.Fatal(app.run(mux))
}
