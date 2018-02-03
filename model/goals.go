package model

import (
	"encoding/json"
	"errors"
	"fmt"
)

// GetSecondaryGoals fetches the secondary goals in the database.
func GetSecondaryGoals(uid string) (GoalResponse, error) {
	var goals []SecondaryGoal
	var _goals []TransformedGoal
	var _goalResponse GoalResponse

	db.Find(&goals, "user_id = ?", uid)

	if len(goals) == 0 {
		// respond with 404
		return GoalResponse{}, ErrorNotFound
	}

	for _, item := range goals {
		id := int(item.ID)
		_goals = append(_goals, TransformedGoal{GoalID: id, UserID: item.UserID, GoalDate: item.GoalDate, GoalWeight: item.GoalWeight, Exercise: item.Exercise})
	}

	_goalResponse = GoalResponse{Data: _goals}

	return _goalResponse, nil
}

// PostSecondaryGoals creates or updates goals in the database.
func PostSecondaryGoals(uid string, g1 SecondaryGoal, g2 SecondaryGoal) (Goals, error) {

	var goals []SecondaryGoal

	db.Find(&goals, "user_id = ?", uid)
	if len(goals) > 0 {
		// overwrite existing goals
		fmt.Println("overwriting existing goals")
		db.Delete(&goals)
	}
	// create new goals
	fmt.Println(g1.GoalDate, g1.GoalWeight, g1.Exercise)
	fmt.Println(g2.GoalDate, g2.GoalWeight, g2.Exercise)

	goal1 := SecondaryGoal{UserID: uid, GoalDate: g1.GoalDate, GoalWeight: g1.GoalWeight, Exercise: g1.Exercise}
	goal2 := SecondaryGoal{UserID: uid, GoalDate: g2.GoalDate, GoalWeight: g2.GoalWeight, Exercise: g2.Exercise}
	db.Save(&goal1)
	db.Save(&goal2)

	fmt.Println("DB saved")
	fmt.Println(goal1, goal2)

	// Return a success message (maybe edit later to return the question?)
	return Goals{Goal1: goal1, Goal2: goal2}, nil
}

// GetPrimaryGoal fetches the user's primary goal from the Users db table
func GetPrimaryGoal(uid string) ([]byte, error) {
	var user UserModel
	db.First(&user, "id = ?", uid)
	if user.ID == 0 {
		return nil, errors.New("Not found")
	}

	var _goal = PrimaryGoal{Goal: user.Goal}

	js, err := json.Marshal(_goal)

	return js, err
}

// SetPrimaryGoal updates the user's primary goal in the database.
func SetPrimaryGoal(uid string, b []byte) ([]byte, error) {
	var user UserModel
	var newGoal PrimaryGoal

	err := json.Unmarshal(b, &newGoal)

	if err != nil {
		return nil, errors.New("Something went wrong")
	}
	db.First(&user, "id = ?", uid)
	db.Model(&user).Update("goal", newGoal.Goal)

	return []byte("{\"message\": \"Goal successfully updated\"}"), nil

}
