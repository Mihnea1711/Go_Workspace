package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	db_controller "project3/demo/database_controller"
	"strings"
)

// for postman request
type AllUsers []db_controller.User

var UsersArr = AllUsers{}

// in next function we get all info from postman for User table
func GetInfoUser(w http.ResponseWriter, r *http.Request, newUser *db_controller.User) {
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(reqBody, newUser) //not &newArray because is already pointer

	// exist := db_controller.CheckIfUsernameExists(db_controller.Db.DB, newUser.Username)

	if newUser.Username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	if newUser.Password == "" {
		http.Error(w, "password is required", http.StatusBadRequest)
		return
	}
	UsersArr = append(UsersArr, *newUser)

	w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(*newProduct)
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	//register User in db, we will pass the username and password in the request and the token will be an UUID
	//generated before saving the user in db
	var newUser db_controller.User
	GetInfoUser(w, r, &newUser)

	w.Header().Set("Content-Type", "application/json")
	user, err := db_controller.CreateUser(db_controller.Db.DB, newUser)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") && strings.Contains(err.Error(), "username") {
			http.Error(w, "Username already exists", http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to create user", http.StatusBadRequest)
		}
		return
	}

	json.NewEncoder(w).Encode("User created")
	json.NewEncoder(w).Encode(user)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	//login User in db, we will pass the username and password in the request and the token will be an UUID
	//generated before saving the user in db
	var newUser db_controller.User
	GetInfoUser(w, r, &newUser)
	result, err := db_controller.ReturningTokenUserLogin(db_controller.Db.DB, newUser)
	if err != nil {
		// log.Fatal(err)
		fmt.Fprintln(w, "Username or password is incorrect")
	} else {
		log.Println("\nUser logged in")
		fmt.Fprintln(w, "token is: ")
		json.NewEncoder(w).Encode(result)
	}
}
