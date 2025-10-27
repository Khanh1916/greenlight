package main

import (
	"fmt"
	"net/http"
)

// logging an error message
func (app *application) logError(r *http.Request, err error) {
	app.logger.Println(err)
}

// generic helper sending JSON formatted error message with given status
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// log the error at runtime, send 500 status and JSON response (contaning generic error message) to client
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	message := "The server encountered a problem and could not process your request."
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// send Not found status 404 and JSON response to client
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "The request resource could not be found."
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// send method not allowed status and JSON response
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The %s method is not supported for this resourse.", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// send bad request JSON response to client
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}
