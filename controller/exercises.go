package controller

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/katreinhart/push-it-api/model"
)

// ListExercises is index route for /api/exercises
func ListExercises(w http.ResponseWriter, r *http.Request) {
	var _ex []model.TransformedExercise

	_ex, err := model.ListExercises()
	if err != nil {
		handleErrorAndRespond(nil, err, w)
		return
	}

	js, err := json.Marshal(_ex)
	handleErrorAndRespond(js, err, w)
}

// FetchSingleExercise is GET one route for /api/exercises
func FetchSingleExercise(w http.ResponseWriter, r *http.Request) {
	// get the URL parameter from the http request
	vars := mux.Vars(r)
	id, _ := vars["id"]

	_ex, err := model.FetchSingleExercise(id)

	if err != nil {
		handleErrorAndRespond(nil, err, w)
		return
	}

	js, err := json.Marshal(_ex)
	handleErrorAndRespond(js, err, w)
}

// CreateExercise handles POST calls to /api/exercises
func CreateExercise(w http.ResponseWriter, r *http.Request) {
	// Create a new buffer to read the body, then parse into a []byte
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	// create data
	var ex model.Exercise

	err := json.Unmarshal(b, &ex)

	if err != nil {
		handleErrorAndRespond(nil, err, w)
		return
	}

	_ex := model.CreateExercise(ex)

	js, err := json.Marshal(_ex)
	handleErrorAndRespond(js, err, w)
}
