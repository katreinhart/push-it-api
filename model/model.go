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

	// UpdateWorkout format sent from front end at completion of workout
	UpdateWorkout struct {
		ID        string `json:"id"`
		Completed bool   `json:"completed"`
		Rating    int    `json:"rating"`
		Comments  string `json:"comments"`
	}

	// CompletedWorkout is used to send all nested info about finished workouts back to front end
	CompletedWorkout struct {
		User      string                       `json:"uid"`
		WorkoutID string                       `json:"workout_id"`
		Start     time.Time                    `json:"start_time"`
		End       time.Time                    `json:"finish_time"`
		Rating    int                          `json:"rating"`
		Comments  string                       `json:"comments"`
		Exercises []TransformedWorkoutExercise `json:"exercises"`
		Sets      []WorkoutSet                 `json:"sets"`
	}

	// SavedWorkout returns information about workout that has not been completed.
	SavedWorkout struct {
		UserID    string                       `json:"uid"`
		CreatedAt time.Time                    `json:"created"`
		WorkoutID string                       `json:"workout_id"`
		Exercises []TransformedWorkoutExercise `json:"exercises"`
	}

	// Exercise is the storage representation of an individual exercise.
	Exercise struct {
		gorm.Model
		Name string `json:"ex_name"`
		Link string `json:"info_url"`
	}

	// TransformedExercise is the version sent back to the front end.
	TransformedExercise struct {
		ID   uint   `json:"id"`
		Name string `json:"ex_name"`
		Link string `json:"info_url"`
	}

	// WorkoutExercise is the database representation of each exercise performed.
	WorkoutExercise struct {
		gorm.Model
		WorkoutID      string `json:"workout_id"`
		ExerciseID     uint   `json:"exercise_id"`
		GoalWeight     int    `json:"goal_weight"`
		GoalSets       int    `json:"goal_sets"`
		GoalRepsPerSet int    `json:"goal_reps_per_set"`
	}

	// TransformedWorkoutExercise is the representation sent back to the front end
	TransformedWorkoutExercise struct {
		WorkoutID      string `json:"workout_id"`
		ExerciseID     uint   `json:"exercise_id"`
		ExerciseName   string `json:"exercise_name"`
		GoalWeight     int    `json:"goal_weight"`
		GoalSets       int    `json:"goal_sets"`
		GoalRepsPerSet int    `json:"goal_reps_per_set"`
	}

	// WorkoutExerciseAsPosted is the representation sent in from the front end
	WorkoutExerciseAsPosted struct {
		ExerciseName   string `json:"exercise_name"`
		GoalWeight     int    `json:"goal_weight"`
		GoalSets       int    `json:"goal_sets"`
		GoalRepsPerSet int    `json:"goal_reps"`
	}

	// WorkoutExerciseSet is the database representation of a single set
	WorkoutExerciseSet struct {
		gorm.Model
		WorkoutID     string `json:"workout_id"`
		ExerciseName  string `json:"exercise_name"`
		Weight        int    `json:"weight"`
		RepsAttempted int    `json:"reps_att"`
		RepsCompleted int    `json:"reps_comp"`
	}

	// WorkoutSet is the representation sent back and forth from the front end
	WorkoutSet struct {
		Weight        int    `json:"weight"`
		Exercise      string `json:"exercise_name"`
		RepsAttempted int    `json:"reps_att"`
		RepsCompleted int    `json:"reps_comp"`
	}

	// PostedTimestamps is the format that timestamps come in
	PostedTimestamps struct {
		StartedAt  time.Time `json:"started_at"`
		FinishedAt time.Time `json:"finished_at"`
	}

	// Goals takes in two goals from front end
	Goals struct {
		Goal1 SecondaryGoal `json:"goal1"`
		Goal2 SecondaryGoal `json:"goal2"`
	}

	// GoalResponse is two goals sent back to front end
	GoalResponse struct {
		Data []TransformedGoal `json:"data"`
	}

	// SecondaryGoal is format saved in the database
	SecondaryGoal struct {
		gorm.Model
		UserID     string    `json:"uid"`
		GoalDate   time.Time `json:"goal_date"`
		Exercise   string    `json:"exercise"`
		GoalWeight int       `json:"goal_weight"`
	}

	// TransformedGoal is the format sent back to front end
	TransformedGoal struct {
		GoalID     int       `json:"goal_id"`
		UserID     string    `json:"uid"`
		GoalDate   time.Time `json:"goal_date"`
		Exercise   string    `json:"exercise"`
		GoalWeight int       `json:"goal_weight"`
	}

	// PrimaryGoal is format for passing primary goal back and forth
	PrimaryGoal struct {
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
	db.AutoMigrate(&Exercise{})
	db.AutoMigrate(&SecondaryGoal{})
	db.AutoMigrate(&WorkoutModel{})
	db.AutoMigrate(&WorkoutExercise{})
	db.AutoMigrate(&WorkoutExerciseSet{})
}
