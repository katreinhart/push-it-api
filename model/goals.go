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
	var _goalResponse goalsResponse

	db.Find(&goals, "user_id = ?", uid)

	if len(goals) == 0 {
		// respond with 404
		return nil, errors.New("Not found")
	}

	for _, item := range goals {
		id := int(item.ID)
		_goals = append(_goals, transformedGoal{GoalID: id, UserID: item.UserID, GoalDate: item.GoalDate, GoalWeight: item.GoalWeight, Exercise: item.Exercise})
	}

	_goalResponse = goalsResponse{Data: _goals}

	js, err := json.Marshal(_goalResponse)

	return js, err
}

// PostSecondaryGoals creates or updates goals in the database.
func PostSecondaryGoals(uid string, b []byte) ([]byte, error) {
	fmt.Println("Post secondary goals MODEL function")
	var goals []secondaryGoal
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
	date1, err := time.Parse("2018-04-29 20:46:09 +0000", postGoals.Goal1.GoalDate)
	date2, err := time.Parse("2018-04-29 20:46:09 +0000", postGoals.Goal2.GoalDate)

	fmt.Println(date1, "\n", date2)

	weight1, err := strconv.Atoi(postGoals.Goal1.GoalWeight)
	weight2, err := strconv.Atoi(postGoals.Goal2.GoalWeight)

	if err != nil {
		fmt.Println("Something went wrong parsing data")
		return nil, errors.New("Something went wrong")
	}

	db.Find(&goals, "user_id = ?", uid)
	if len(goals) > 0 {
		// overwrite existing goals
		fmt.Println("overwriting existing goals")
		goals[0] = secondaryGoal{UserID: uid, GoalDate: date1, GoalWeight: weight1, Exercise: postGoals.Goal1.Exercise}
		goals[1] = secondaryGoal{UserID: uid, GoalDate: date2, GoalWeight: weight2, Exercise: postGoals.Goal2.Exercise}
		db.Save(&goals)
		fmt.Println("successfully overwrote goals")
	} else {
		// create new goals
		goal1 = secondaryGoal{UserID: uid, GoalDate: date1, GoalWeight: weight1, Exercise: postGoals.Goal1.Exercise}
		goal2 = secondaryGoal{UserID: uid, GoalDate: date2, GoalWeight: weight2, Exercise: postGoals.Goal2.Exercise}
		db.Save(&goal1)
		db.Save(&goal2)
	}

	fmt.Println("DB saved")

	// Return a success message (maybe edit later to return the question?)
	return []byte("{\"message\": \"Goals successfully added\"}"), nil
}

// SetPrimaryGoal updates the user's primary goal in the database.
func SetPrimaryGoal(uid string, b []byte) ([]byte, error) {
	var user userModel
	var newGoal updatePrimaryGoal

	err := json.Unmarshal(b, &newGoal)

	if err != nil {
		return nil, errors.New("Something went wrong")
	}
	db.First(&user, "id = ?", uid)
	db.Model(&user).Update("goal", newGoal.Goal)

	return []byte("{\"message\": \"Goal successfully updated\"}"), nil

}
