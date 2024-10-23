package main

import (
	"github.com/Kser1x/social/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"time"
)

// описываем класс application
type application struct {
	config  config
	storage store.Storage
}

// описываем класс config
type config struct {
	addr string
	db   dbConfig
}
type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

// метод который создаёт mux
func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
	})

	return r
}

// Функция которая запускает наш сервер
func (app *application) run(mux http.Handler) error {

	//создаем переменнную srv которая создает сервер с определенными параметрами
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute, //время простоя
	}

	log.Printf("server has started at %s", app.config.addr)

	return srv.ListenAndServe()
}
