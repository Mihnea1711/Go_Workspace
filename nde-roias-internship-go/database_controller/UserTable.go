package database_controller

import (
	"database/sql"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"

	uuid "github.com/google/uuid"
)

// //Structs

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

func CreateUserTable(db *sql.DB, tableName string, dbName string) error {

	UseDB(dbName)
	query := `CREATE TABLE IF NOT EXISTS ` + tableName + ` (
					id int primary key auto_increment, 
					username varchar(30) NOT NULL UNIQUE,
					password varchar(80) NOT NULL,
					token varchar(36) NOT NULL UNIQUE
			)`
	_, err := db.Exec(query)

	if err != nil {
		log.Fatalf("Failed to create user table: %v", err)
		return err
	}

	return nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func CreateUser(db *sql.DB, newUser User) (User, error) {
	var err error
	newUser.Token = uuid.New().String()
	newUser.Password, err = HashPassword(newUser.Password)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
		return newUser, err
	}

	query := `INSERT INTO User (username, password, token) VALUES (?, ?, ?)`
	_, err = db.Exec(query, newUser.Username, newUser.Password, newUser.Token)
	if err != nil {
		return newUser, err
	}

	//there we will generate random reviews. Every time when we create new user, we will generate random reviews for some products.
	RandomReviewsGenerator(Db.DB)
	CreateRatingPerProduct(Db.DB)

	return newUser, nil
}

func VerifyUserExistence(db *sql.DB, newUser User) (bool, User, error) {
	var err error
	var user User
	query := `SELECT * FROM User WHERE username = ?`
	err = db.QueryRow(query, newUser.Username).Scan(&user.ID, &user.Username, &user.Password, &user.Token)
	if err != nil {
		return false, user, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newUser.Password))
	if err != nil {
		log.Fatal(err)
		return false, user, err
	}
	return true, user, nil
}

func ReturningTokenUserLogin(db *sql.DB, newUser User) (string, error) {
	var user User
	existing, user, err := VerifyUserExistence(db, newUser)
	if !existing {
		// log.Fatal(err)
		return "User doesn't exist", err
	}

	return user.Token, nil
}

func GetListUsername(db *sql.DB) ([]string, error) {
	var usernames []string
	rows, err := db.Query("SELECT username FROM User")
	if err != nil {
		return nil, errors.New("error trying to select username from User")
	}
	defer rows.Close()
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, errors.New("error trying to scan username from User")
		}
		usernames = append(usernames, username)
	}
	return usernames, nil
}

func FindUserByToken(db *sql.DB, token string) (string, error) {
	var username string
	query := `SELECT username FROM User WHERE token = ?`
	err := db.QueryRow(query, token).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}
