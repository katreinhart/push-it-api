package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type (
	userModel struct {
		gorm.Model
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Goal     string `json:"primary_goal"`
	}

	transformedUser struct {
		Email string `json:"email"`
		Name  string `json:"name"`
		Goal  string `json:"goal"`
	}

	workout struct {
		gorm.Model
		User  uint      `json:"uid"`
		Start time.Time `json:"start_time"`
		End   time.Time `json:"finish_time"`
	}

	exercise struct {
		gorm.Model
		Name string `json:"ex_name"`
		Link string `json:"info_url"`
	}

	workoutExercise struct {
		gorm.Model
		WorkoutID  uint `json:"workout_id"`
		ExerciseID uint `json:"exercise_id"`
	}

	workoutExerciseSet struct {
		gorm.Model
		WorkoutExerciseID uint `json:"workout_exercise_id"`
		Weight            uint `json:"weight"`
		RepsAttempted     uint `json:"reps_att"`
		RepsCompleted     uint `json:"reps_comp"`
	}

	secondaryGoal struct {
		gorm.Model
		SetDate    time.Time `json:"set_date"`
		GoalDate   time.Time `json:"goal_date"`
		Exercise   uint      `json:"exercise_id"`
		GoalWeight uint      `json:"goal_weight"`
	}
)
