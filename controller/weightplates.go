package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/katreinhart/push-it-api/model"
)

// FindWeightPlates handles incoming requests to calculate weight plates
func FindWeightPlates(w http.ResponseWriter, r *http.Request) {
	// get the URL parameter from the http request
	vars := mux.Vars(r)
	weight, _ := vars["weight"]

	intWeight, _ := strconv.Atoi(weight)

	var weightResult = model.FindWeightPlates(intWeight, 45)

	js, err := json.Marshal(&weightResult)

	handleErrorAndRespond(js, err, w)
}
