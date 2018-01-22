package model

import (
	"encoding/json"
	"errors"
	"fmt"
)

// GetSecondaryGoals fetches the secondary goals in the database.
func GetSecondaryGoals(uid string) ([]byte, error) {
	var goals []secondaryGoal
	var _goals []transformedGoal

	db.Find(&goals, "user_id = ?", uid)

	if len(goals) == 0 {
		// respond with 404
	}

	for _, item := range goals {
		_goals = append(_goals, transformedGoal{GoalID: item.ID, UserID: item.UserID, GoalDate: item.GoalDate, GoalWeight: item.GoalWeight, Exercise: item.Exercise})
	}

	js, err := json.Marshal(_goals)

	return js, err
}

// PostSecondaryGoals creates or updates goals in the database.
func PostSecondaryGoals(uid string, b []byte) ([]byte, error) {
	fmt.Println("Post secondary goals MODEL function")
	var postGoals postedGoals
	var goals []secondaryGoal
	// Unmarshal the JSON formatted data b into the struct
	err := json.Unmarshal(b, &postGoals)

	if err != nil {
		// handle error
		fmt.Println("JSON error")
		return nil, errors.New("error marshaling json object")
	}

	goals[0] = secondaryGoal{UserID: uid, Exercise: postGoals.Goal1.Exercise, GoalWeight: postGoals.Goal1.GoalWeight, GoalDate: postGoals.Goal1.GoalDate}
	goals[0] = secondaryGoal{UserID: uid, Exercise: postGoals.Goal2.Exercise, GoalWeight: postGoals.Goal2.GoalWeight, GoalDate: postGoals.Goal2.GoalDate}

	goals[0].UserID = uid
	goals[1].UserID = uid

	// Handle error if any
	if err != nil {
		fmt.Println("Somethign wrong")
		return []byte("{\"message\": \"Something went wrong.\"}"), err
	}

	db.Save(&goals[0])
	db.Save(&goals[1])
	fmt.Println("DB saved")

	// Return a success message (maybe edit later to return the question?)
	return []byte("{\"message\": \"Goals successfully added\"}"), nil
}
