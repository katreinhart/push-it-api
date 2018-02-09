package model

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
		db.Delete(&goals)
	}

	goal1 := SecondaryGoal{UserID: uid, GoalDate: g1.GoalDate, GoalWeight: g1.GoalWeight, Exercise: g1.Exercise}
	goal2 := SecondaryGoal{UserID: uid, GoalDate: g2.GoalDate, GoalWeight: g2.GoalWeight, Exercise: g2.Exercise}

	db.Save(&goal1)
	db.Save(&goal2)

	// Return a success message (maybe edit later to return the question?)
	return Goals{Goal1: goal1, Goal2: goal2}, nil
}

// GetPrimaryGoal fetches the user's primary goal from the Users db table
func GetPrimaryGoal(uid string) (PrimaryGoal, error) {
	var user UserModel

	db.First(&user, "id = ?", uid)
	if user.ID == 0 {
		return PrimaryGoal{}, ErrorNotFound
	}

	var _goal = PrimaryGoal{Goal: user.Goal}

	return _goal, nil
}

// SetPrimaryGoal updates the user's primary goal in the database.
func SetPrimaryGoal(uid string, newGoal PrimaryGoal) (PrimaryGoal, error) {
	var user UserModel

	db.First(&user, "id = ?", uid)
	db.Model(&user).Update("goal", newGoal.Goal)

	return newGoal, nil
}
