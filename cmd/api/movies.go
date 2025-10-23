package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Khanh1916/greenlight/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new mobie")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	movie := data.Movie{
		ID:       id,
		CreateAt: time.Now(),
		Title:    "Toiditimtoi",
		Runtime:  190,
		Genres:   []string{"drama", "romance", "comedy"},
		Version:  1,
	}

	err = app.writeJSON(w, http.StatusOK, movie, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
