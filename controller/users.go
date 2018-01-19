package controller

import (
	"bytes"
	"net/http"

	"github.com/katreinhart/gorilla-api/model"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Read in http request body
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	// Create user in model
	js, err := model.CreateUser(b)
	handleErrorAndRespond(js, err, w)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {

}
