package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
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
		id := int(item.ID)
		_goals = append(_goals, transformedGoal{GoalID: id, UserID: item.UserID, GoalDate: item.GoalDate, GoalWeight: item.GoalWeight, Exercise: item.Exercise})
	}

	js, err := json.Marshal(_goals)

	return js, err
}

// PostSecondaryGoals creates or updates goals in the database.
func PostSecondaryGoals(uid string, b []byte) ([]byte, error) {
	fmt.Println("Post secondary goals MODEL function")

	var postGoals postedGoals
	var goal1, goal2 secondaryGoal

	// Unmarshal the JSON formatted data b into the struct
	err := json.Unmarshal(b, &postGoals)

	if err != nil {
		// handle error
		fmt.Println("JSON error")
		fmt.Fprintf(os.Stdout, "%#v", postGoals)
		return nil, errors.New("error marshaling json object")
	}

	// Calculate the dates
	date1, err := time.Parse("Jan 01, 2018", postGoals.Goal1.GoalDate)
	date2, err := time.Parse("Jan 01, 2018", postGoals.Goal2.GoalDate)

	weight1, err := strconv.Atoi(postGoals.Goal1.GoalWeight)
	weight2, err := strconv.Atoi(postGoals.Goal2.GoalWeight)

	if err != nil {
		// handle date error
	}

	// Transform data to db format
	goal1 = secondaryGoal{UserID: postGoals.Goal1.UserID, GoalDate: date1, GoalWeight: weight1, Exercise: postGoals.Goal1.Exercise}
	goal2 = secondaryGoal{UserID: postGoals.Goal2.UserID, GoalDate: date2, GoalWeight: weight2, Exercise: postGoals.Goal2.Exercise}

	// Handle error if any
	if err != nil {
		fmt.Println("Somethign wrong")
		return []byte("{\"message\": \"Something went wrong.\"}"), err
	}

	db.Save(&goal1)
	db.Save(&goal2)
	fmt.Println("DB saved")

	// Return a success message (maybe edit later to return the question?)
	return []byte("{\"message\": \"Goals successfully added\"}"), nil
}