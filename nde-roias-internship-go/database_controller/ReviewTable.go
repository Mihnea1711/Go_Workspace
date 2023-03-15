package database_controller

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
)

type Review struct {
	ID          int    `json:"id"`
	ProductUUID string `json:"product_uuid"`
	Description string `json:"description"`
	Username    string `json:"username"`
	Rating      int    `json:"rating"`
}

func CreateReviewTable(db *sql.DB, tableName string, dbName string) error {
	UseDB(dbName)
	query := `CREATE TABLE IF NOT EXISTS ` + tableName + ` (
					id int primary key auto_increment, 
					product_uuid varchar(36) NOT NULL,
					description text,
					username varchar(30) NOT NULL,
					rating int NOT NULL		
					)`
	_, err := db.Exec(query)

	if err != nil {
		log.Fatalf("\nError in creating Review Table: %v", err)
	}

	return nil
}

func InsertReview(db *sql.DB, review Review) (Review, error) {
	var err error
	query := `INSERT INTO Review (product_uuid, description, username, rating) VALUES (?, ?, ?, ?)`
	_, err = db.Exec(query, review.ProductUUID, review.Description, review.Username, review.Rating)
	if err != nil {
		return review, err
	}
	return review, nil
}

func CheckUserExists(db *sql.DB, username string) bool {
	var userExists bool
	query := `SELECT EXISTS(SELECT 1 FROM User WHERE username = ?)`
	err := db.QueryRow(query, username).Scan(&userExists)
	if err != nil {
		log.Fatal(err)
	}
	return userExists
}

func RandomReviewsGenerator(db *sql.DB) error {
	count, err := CountRowsInTable("Product", db)
	if err != nil {
		return errors.New("failed to count rows in product table")
	}

	productUUIDs, err := GetListProductUuid(db)
	if err != nil {
		return errors.New("failed to get product uuids")
	}

	usernames, err := GetListUsername(db)
	if err != nil {
		return errors.New("failed to get usernames")
	}

	//Because I generate this function every time when a user is created, I don't want to have a lot of reviews for each product for each user. That's why I divide the count by 5.
	for i := 0; i < count/5; i++ {
		var newReview Review
		if len(productUUIDs) > 0 {
			newReview.ProductUUID = productUUIDs[rand.Intn(len(productUUIDs))]
		} else {
			return errors.New("empty list of product uuids")
		}
		if len(usernames) > 0 {
			newReview.Username = usernames[rand.Intn(len(usernames))]
		} else {
			return errors.New("empty list of usernames")
		}
		newReview.Rating = rand.Intn(5) + 1
		newReview.Description = fmt.Sprintf("Review %d for product %s by user %s", i, newReview.ProductUUID, newReview.Username)
		_, err := InsertReview(db, newReview)
		if err != nil {
			log.Fatalf("Failed to insert random review: %v", err)
		}
	}
	return nil
}

func GetReviewsForAProduct(db *sql.DB, productUuid string) ([]Review, error) {
	query := `SELECT * FROM Review WHERE product_uuid = ?`
	rows, err := db.Query(query, productUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []Review
	for rows.Next() {
		var review Review
		err := rows.Scan(&review.ID, &review.ProductUUID, &review.Description, &review.Username, &review.Rating)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}

func UpdateReview(db *sql.DB, username string, productUuid string, description string, rating string) error {

	query := `UPDATE Review SET description = ?, rating = ? WHERE username = ? AND product_uuid = ?`
	_, err := db.Exec(query, description, rating, username, productUuid)
	if err != nil {
		return err
	}
	return nil
}
