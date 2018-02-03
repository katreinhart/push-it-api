package main

import (
	"net/http"
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/codegangsta/negroni"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/katreinhart/push-it-api/controller"
	"github.com/rs/cors"
)

func main() {

	// get port variable from environment or set to default
	var port string
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	} else {
		port = "8080"
	}

	// CORS middleware setup
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"Accept-Encoding", "Accept-Language", "Authorization"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "OPTIONS", "DELETE"},
		AllowCredentials: true,
	})

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", homeHandler)

	// s is a subrouter to handle question routes
	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/exercises", controller.ListExercises).Methods("GET")
	api.HandleFunc("/exercises/{id}", controller.FetchSingleExercise).Methods("GET")
	api.HandleFunc("/exercises", controller.CreateExercise).Methods("POST")

	u := r.PathPrefix("/auth").Subrouter()
	u.HandleFunc("/register", controller.CreateUser).Methods("POST")
	u.HandleFunc("/login", controller.LoginUser).Methods("POST")

	// User functions
	api.HandleFunc("/setinfo", controller.SetUserInfo).Methods("POST")

	// Goal functions
	api.HandleFunc("/user/goals", controller.GetSecondaryGoals).Methods("GET")
	api.HandleFunc("/user/goals", controller.PostSecondaryGoals).Methods("POST")
	api.HandleFunc("/user/primarygoal", controller.GetPrimaryGoal).Methods("GET")
	api.HandleFunc("/user/primarygoal", controller.SetPrimaryGoal).Methods("POST")

	// Workout functions
	api.HandleFunc("/workouts", controller.CreateWorkout).Methods("POST")
	api.HandleFunc("/workouts/{id}", controller.GetWorkout).Methods("GET")
	api.HandleFunc("/workouts/{id}", controller.MarkWorkoutAsCompleted).Methods("PUT")
	api.HandleFunc("/workouts/{id}", controller.UpdateWorkoutTimestamps).Methods("PATCH")

	// Exercise functions
	api.HandleFunc("/workouts/{id}/exercises", controller.AddExerciseToWorkout).Methods("POST")
	api.HandleFunc("/workouts/{id}/exercises/sets", controller.AddExerciseSet).Methods("POST")

	api.HandleFunc("/history", controller.History).Methods("GET")
	api.HandleFunc("/saved", controller.FetchSavedExercises).Methods("GET")

	// JWT Middleware handles authorization configuration
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET")), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	// muxRouter uses Negroni & handles the middleware for authorization
	muxRouter := http.NewServeMux()
	muxRouter.Handle("/", r)
	muxRouter.Handle("/api/", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(api),
	))

	// Negroni handles the middleware chaining with next
	n := negroni.Classic()

	// Use CORS
	n.Use(c)

	// handle routes with the muxRouter
	n.UseHandler(muxRouter)

	// listen and serve!
	http.ListenAndServe(":"+port, handlers.RecoveryHandler()(n))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"message\": \"Hello world\"}"))
}
