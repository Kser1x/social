package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("внутреняя ошибка сервера: %s error: %s", r.Method, r.URL.Path, err)

	writeJSON(w, http.StatusInternalServerError, "сервер обноружил проблемку")
}
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error: %s error: %s", r.Method, r.URL.Path, err)

	writeJSON(w, http.StatusBadRequest, err.Error())
}
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error: %s error: %s", r.Method, r.URL.Path, err)

	writeJSON(w, http.StatusNotFound, "Not found")
}
