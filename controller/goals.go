package controller

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/katreinhart/push-it-api/model"
)

// GetSecondaryGoals responds to requests for user's secondary goals
func GetSecondaryGoals(w http.ResponseWriter, r *http.Request) {
	// parse user ID from token
	uid, err := GetUIDFromBearerToken(r)
	if err != nil {
		handleErrorAndRespond(nil, err, w)
		return
	}
	goals, err := model.GetSecondaryGoals(uid)

	if err != nil {
		handleErrorAndRespond(nil, err, w)
		return
	}

	js, err := json.Marshal(goals)
	handleErrorAndRespond(js, err, w)
}

// PostSecondaryGoals responds to requests to post or update user's secondary goals
func PostSecondaryGoals(w http.ResponseWriter, r *http.Request) {

	uid, err := GetUIDFromBearerToken(r)
	if err != nil {
		handleErrorAndRespond(nil, model.ErrorForbidden, w)
		return
	}

	// get the body from the request
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	// Fetch secondary goals from model

	var pgs model.Goals
	err = json.Unmarshal(b, &pgs)

	goal1 := pgs.Goal1
	goal2 := pgs.Goal2

	if err != nil {
		handleErrorAndRespond(nil, err, w)
		return
	}

	postedGoals, err := model.PostSecondaryGoals(uid, goal1, goal2)

	if err != nil {
		handleErrorAndRespond(nil, err, w)
		return
	}

	js, err := json.Marshal(postedGoals)
	handleErrorAndRespond(js, err, w)
}

// GetPrimaryGoal gets the user's primary goal from the database
func GetPrimaryGoal(w http.ResponseWriter, r *http.Request) {
	uid, err := GetUIDFromBearerToken(r)

	if err != nil {
		handleErrorAndRespond(nil, model.ErrorForbidden, w)
		return
	}

	goal, err := model.GetPrimaryGoal(uid)

	if err != nil {
		handleErrorAndRespond(nil, err, w)
		return
	}

	js, err := json.Marshal(goal)
	handleErrorAndRespond(js, err, w)
}

// SetPrimaryGoal sets the user's primary goal in the database
func SetPrimaryGoal(w http.ResponseWriter, r *http.Request) {

	uid, err := GetUIDFromBearerToken(r)
	if err != nil {
		handleErrorAndRespond(nil, model.ErrorNotFound, w)
		return
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	var newGoal model.PrimaryGoal
	err = json.Unmarshal(b, &newGoal)

	if err != nil {
		handleErrorAndRespond(nil, err, w)
		return
	}

	goal, err := model.SetPrimaryGoal(uid, newGoal)
	if err != nil {
		handleErrorAndRespond(nil, err, w)
		return
	}

	js, err := json.Marshal(goal)
	handleErrorAndRespond(js, err, w)
}
