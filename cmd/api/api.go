package main

import (
	"log"
	"net/http"
	"time"
)

type application struct {
	config config
}
type config struct {
	addr string
}

// метод который создаёт mux
func (app *application) mount() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET / users")

	return mux
}

func (app *application) run(mux *http.ServeMux) error {
	//mux := http.NewServeMux()

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at %s", app.config.addr)

	return srv.ListenAndServe()
}