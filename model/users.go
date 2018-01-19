package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser creates a new user instance in the database
func CreateUser(b []byte) ([]byte, error) {

	// declare data structures to be used
	var user, dbUser userModel

	// Unmarshal the json body into the user data structure
	err := json.Unmarshal(b, &user)

	// handle error if any
	if err != nil {
		return []byte("{\"message\": \"Something went wrong.\"}"), err
	}

	// See if the user exists in the database. If so, return an error (no duplicates allowed)
	db.First(&dbUser, "email = ?", user.Email)
	if dbUser.ID != 0 {
		return []byte("{\"message\": \"User already exists in DB.\"}"), errors.New("User already exists")
	}

	// Hash the user's password using bcrypt helper function and handle any error
	hash, err := hashPassword(user.Password)
	if err != nil {
		return []byte("{\"message\": \"Sorry, something went wrong.\"}"), err
	}

	// Overwrite the user's password with the hashed version (no plaintext storage of passwords)
	user.Password = hash

	// save the user in the DB
	db.Save(&user)

	// Get the user back from the database so you have the correct ID
	db.First(&dbUser, "email = ?", user.Email)

	// create and sign the JWT
	t, err := createAndSignJWT(dbUser)

	// Handle error in JWT creation/signing
	if err != nil {
		fmt.Println(err.Error())
		return []byte("{\"message\": \"Sorry, something went wrong.\"}"), err
	}

	// create transformed version of user structure, marshal it into JSON and return
	_user := transformedUser{ID: user.ID, Email: user.Email, Password: user.Password, Goal: user.Goal}
	js, err := json.Marshal(_user)
	return js, err
}

// --------------------- Helper Functions ---------------------
// user login password helper functions
// from https://gowebexamples.com/password-hashing/
func hashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(b), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// JWT helper function
func createAndSignJWT(user userModel) (string, error) {

	// create the expiration time, build claim, create and sign token, and return.
	e := time.Now().Add(time.Hour * 24).Unix()
	c := CustomClaims{
		user.ID,
		user.Admin,
		jwt.StandardClaims{
			ExpiresAt: e,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	secret := []byte(os.Getenv("SECRET"))
	t, err := token.SignedString(secret)
	return t, err
}
