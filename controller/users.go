package controller

import (
	"bytes"
	"net/http"

	"github.com/katreinhart/push-it-api/model"
)

// CreateUser handles POST requests to /auth/register
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Read in http request body
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	// Create user in model
	js, err := model.CreateUser(b)
	handleErrorAndRespond(js, err, w)
}

// LoginUser handles POST requests to /auth/login
func LoginUser(w http.ResponseWriter, r *http.Request) {

}
