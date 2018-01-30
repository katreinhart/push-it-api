package controller

import (
	"bytes"
	"errors"
	"net/http"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
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

	// Read in http request body
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	// Login user in model
	js, err := model.LoginUser(b)
	handleErrorAndRespond(js, err, w)
}

// SetUserInfo handles requests to update user info
func SetUserInfo(w http.ResponseWriter, r *http.Request) {
	// Read in http request body
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	// Get user ID from token
	uid, err := GetUIDFromBearerToken(r)
	if err != nil {
		handleErrorAndRespond(nil, errors.New("Forbidden"), w)
	}

	// Login user in model
	js, err := model.SetUserInfo(uid, b)
	handleErrorAndRespond(js, err, w)
}

// GetUIDFromBearerToken does what it says on the tin
func GetUIDFromBearerToken(r *http.Request) (string, error) {
	user := r.Context().Value("user")
	tok := user.(*jwt.Token)
	var err error

	// no token present, so this is an unauthorized request.
	if tok == nil {
		err = errors.New("Forbidden")
	}

	// get claims from token
	claims := tok.Claims.(jwt.MapClaims)

	// parse uid out of claims.
	uid, ok := claims["uid"].(float64)

	// Error parsing uid from token.
	if !ok {
		err = errors.New("Forbidden")
	}

	// UID parsed from token is of type float64; we need it as a string.
	struid := strconv.FormatFloat(uid, 'f', -1, 64)

	return struid, err
}
