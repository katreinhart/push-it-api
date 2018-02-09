package controller

import (
	"bytes"
	"encoding/json"
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

	var u model.UserModel
	err := json.Unmarshal(b, &u)

	if err != nil {
		handleErrorAndRespond(nil, model.ErrorBadRequest, w)
		return
	}

	_u, err := model.CreateUser(u)
	if err != nil {
		handleErrorAndRespond(nil, model.ErrorUserExists, w)
		return
	}

	// Create user in model
	js, err := json.Marshal(_u)
	handleErrorAndRespond(js, err, w)
}

// LoginUser handles POST requests to /auth/login
func LoginUser(w http.ResponseWriter, r *http.Request) {

	// Read in http request body
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	var u model.UserModel
	var _u model.TransformedUser

	err := json.Unmarshal(b, &u)
	if err != nil {
		handleErrorAndRespond(nil, model.ErrorBadRequest, w)
		return
	}

	// Login user in model
	_u, err = model.LoginUser(u)

	if err != nil {
		handleErrorAndRespond(nil, model.ErrorForbidden, w)
		return
	}

	js, err := json.Marshal(_u)
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

	userID, err := strconv.Atoi(uid)

	if err != nil {
		handleErrorAndRespond(nil, model.ErrorInternalServer, w)
		return
	}
	_userID := uint(userID)

	var u model.UserModel

	err = json.Unmarshal(b, &u)
	if err != nil {
		handleErrorAndRespond(nil, model.ErrorBadRequest, w)
		return
	}

	_u, err := model.SetUserInfo(_userID, u)
	if err != nil {
		handleErrorAndRespond(nil, err, w)
		return
	}

	js, err := json.Marshal(_u)

	handleErrorAndRespond(js, err, w)
}

// GetUIDFromBearerToken does what it says on the tin
func GetUIDFromBearerToken(r *http.Request) (string, error) {
	user := r.Context().Value("user")
	tok := user.(*jwt.Token)
	var err error

	// no token present, so this is an unauthorized request.
	if tok == nil {
		err = model.ErrorForbidden
	}

	// get claims from token
	claims := tok.Claims.(jwt.MapClaims)

	// parse uid out of claims.
	uid, ok := claims["uid"].(float64)

	// Error parsing uid from token.
	if !ok {
		err = model.ErrorForbidden
	}

	// UID parsed from token is of type float64; we need it as a string.
	struid := strconv.FormatFloat(uid, 'f', -1, 64)
	return struid, err
}
