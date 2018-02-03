package model

import (
	"errors"
	"fmt"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"

	// importing postgres dialect for GORM
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// declare DB
var db *gorm.DB

type (
	// UserModel is the database model for users
	UserModel struct {
		gorm.Model
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Goal     string `json:"goal"`
		Level    string `json:"level"`
	}

	// TransformedUser is the version that gets sent back over API
	TransformedUser struct {
		ID    uint   `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
		Goal  string `json:"goal"`
		Level string `json:"level"`
		Token string `json:"token"`
	}

	// WorkoutModel is the database model for storing workouts
	WorkoutModel struct {
		gorm.Model
		UserID    string    `json:"user_id"`
		Start     time.Time `json:"start_time"`
		End       time.Time `json:"finish_time"`
		Rating    int       `json:"rating"`
		Comments  string    `json:"comments"`
		Completed bool      `json:"completed"`
	}

	updateWorkoutModel struct {
		ID        string `json:"id"`
		Completed bool   `json:"completed"`
		Rating    int    `json:"rating"`
		Comments  string `json:"comments"`
	}

	CompletedWorkout struct {
		User      string                       `json:"uid"`
		WorkoutID string                       `json:"workout_id"`
		Start     time.Time                    `json:"start_time"`
		End       time.Time                    `json:"finish_time"`
		Rating    int                          `json:"rating"`
		Comments  string                       `json:"comments"`
		Exercises []TransformedWorkoutExercise `json:"exercises"`
		Sets      []transformedWorkoutSet      `json:"sets"`
	}

	savedWorkout struct {
		User      string                       `json:"uid"`
		WorkoutID string                       `json:"workout_id"`
		Exercises []TransformedWorkoutExercise `json:"exercises"`
	}

	transformedSavedWorkout struct {
		ID        uint                         `json:"id"`
		User      string                       `json:"uid"`
		WorkoutID string                       `json:"workout_id"`
		Exercises []TransformedWorkoutExercise `json:"exercises"`
	}

	exercise struct {
		gorm.Model
		Name string `json:"ex_name"`
		Link string `json:"info_url"`
	}

	transformedExercise struct {
		ID   uint   `json:"id"`
		Name string `json:"ex_name"`
		Link string `json:"info_url"`
	}

	WorkoutExercise struct {
		gorm.Model
		WorkoutID      string `json:"workout_id"`
		ExerciseID     uint   `json:"exercise_id"`
		GoalSets       int    `json:"goal_sets"`
		GoalRepsPerSet int    `json:"goal_reps_per_set"`
	}

	TransformedWorkoutExercise struct {
		WorkoutID      string `json:"workout_id"`
		ExerciseID     uint   `json:"exercise_id"`
		ExerciseName   string `json:"exercise_name"`
		GoalSets       int    `json:"goal_sets"`
		GoalRepsPerSet int    `json:"goal_reps_per_set"`
	}

	WorkoutExerciseAsPosted struct {
		ExerciseName   string `json:"exercise_name"`
		GoalSets       int    `json:"goal_sets"`
		GoalRepsPerSet int    `json:"goal_reps"`
	}

	workoutExerciseSet struct {
		gorm.Model
		WorkoutExerciseID string `json:"workout_exercise_id"`
		WorkoutID         string `json:"workout_id"`
		ExerciseName      string `json:"exercise_name"`
		Weight            int    `json:"weight"`
		RepsAttempted     int    `json:"reps_att"`
		RepsCompleted     int    `json:"reps_comp"`
	}

	workoutSetAsPosted struct {
		Weight        int    `json:"weight"`
		Exercise      string `json:"exercise_name"`
		RepsAttempted int    `json:"reps_att"`
		RepsCompleted int    `json:"reps_comp"`
	}

	transformedWorkoutSet struct {
		ExerciseName  string `json:"exercise"`
		Weight        int    `json:"weight"`
		RepsAttempted int    `json:"reps_att"`
		RepsCompleted int    `json:"reps_comp"`
	}

	postedTimestamps struct {
		StartedAt  time.Time `json:"started_at"`
		FinishedAt time.Time `json:"finished_at"`
	}

	postedGoals struct {
		Goal1 postedBasicGoal `json:"goal1"`
		Goal2 postedBasicGoal `json:"goal2"`
	}

	postedBasicGoal struct {
		UserID     string `json:"uid"`
		GoalDate   string `json:"goal_date"`
		Exercise   string `json:"exercise"`
		GoalWeight string `json:"goal_weight"`
	}

	secondaryGoal struct {
		gorm.Model
		UserID     string    `json:"uid"`
		GoalDate   time.Time `json:"goal_date"`
		Exercise   string    `json:"exercise"`
		GoalWeight int       `json:"goal_weight"`
	}

	transformedGoal struct {
		GoalID     int       `json:"goal_id"`
		UserID     string    `json:"uid"`
		SetDate    time.Time `json:"set_date"`
		GoalDate   time.Time `json:"goal_date"`
		Exercise   string    `json:"exercise"`
		GoalWeight int       `json:"goal_weight"`
	}

	goalsResponse struct {
		Data []transformedGoal `json:"data"`
	}

	transformedPrimaryGoal struct {
		Goal string `json:"goal"`
	}

	updatePrimaryGoal struct {
		Goal string `json:"goal"`
	}

	// CustomClaims for JWT handling
	CustomClaims struct {
		UID uint `json:"uid"`
		jwt.StandardClaims
	}
)

// Errors
var ErrorBadRequest = errors.New("Bad request")
var ErrorUserExists = errors.New("User exists in db")
var ErrorInternalServer = errors.New("Something went wrong")
var ErrorForbidden = errors.New("Forbidden")
var ErrorNotFound = errors.New("Not found")

// init function runs at setup; connects to database
func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading env file")
	}

	hostname := os.Getenv("HOST")
	dbname := os.Getenv("DBNAME")
	uname := os.Getenv("DBUSER")
	password := os.Getenv("PASSWORD")

	dbString := "host=" + hostname + " user=" + uname + " dbname=" + dbname + " sslmode=disable password=" + password

	// var err error
	db, err = gorm.Open("postgres", dbString)
	if err != nil {
		fmt.Println(err.Error())
		panic("Unable to connect to DB")
	}

	db.AutoMigrate(&UserModel{})
	db.AutoMigrate(&exercise{})
	db.AutoMigrate(&secondaryGoal{})
	db.AutoMigrate(&WorkoutModel{})
	db.AutoMigrate(&WorkoutExercise{})
	db.AutoMigrate(&workoutExerciseSet{})
}
