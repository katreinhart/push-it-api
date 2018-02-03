package controller

import (
	"net/http"

	"github.com/katreinhart/push-it-api/model"
)

func handleErrorAndRespond(js []byte, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	// handle the error cases
	if err != nil {
		if err == model.ErrorBadRequest {
			w.WriteHeader(http.StatusBadRequest)
		} else if err == model.ErrorUserExists {
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else if err == model.ErrorInternalServer {
			w.WriteHeader(http.StatusInternalServerError)
		} else if err == model.ErrorForbidden {
			w.WriteHeader(http.StatusUnauthorized)
		} else if err == model.ErrorNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		// in any case, send back the message if there is one
		w.Write(js)
		return
	}

	// Handle the success case
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
