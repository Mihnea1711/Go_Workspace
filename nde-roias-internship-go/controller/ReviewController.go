package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	db_controller "project3/demo/database_controller"
	"strconv"

	"github.com/gorilla/mux"
)

// for postman request
type AllReviews []db_controller.Review

var ReviewsArr = AllReviews{}

func ValidateInfoReview(w http.ResponseWriter, r *http.Request, newReview *db_controller.Review) bool {
	status := true
	if newReview.Username == "" {
		status = false
		http.Error(w, "username is required", http.StatusBadRequest)
		return status
	}

	userExists := db_controller.CheckUserExists(db_controller.Db.DB, newReview.Username)
	if !userExists {
		status = false
		http.Error(w, "user does not exist", http.StatusBadRequest)
		return status
	}

	if newReview.Rating == 0 {
		status = false
		http.Error(w, "rating is required", http.StatusBadRequest)
		return status
	}
	if newReview.Rating < 1 || newReview.Rating > 5 {
		status = false
		http.Error(w, "rating must be between 1 and 5", http.StatusBadRequest)
		return status
	}

	return status
}

// in next function we get all info from postman for Review table
func GetInfoReview(w http.ResponseWriter, r *http.Request, newReview *db_controller.Review) {
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(reqBody, newReview) //not &newArray because is already pointer

	ReviewsArr = append(ReviewsArr, *newReview)

	w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(*newReview)
}

func CreateReview(w http.ResponseWriter, r *http.Request) {
	//create Review in db\
	reqVariables := mux.Vars(r)
	productUUID, status := reqVariables["productUuid"]
	if !status {
		log.Fatal("The productID parameter is missing from the request")
	}
	var newReview db_controller.Review
	GetInfoReview(w, r, &newReview)
	newReview.ProductUUID = productUUID
	if !ValidateInfoReview(w, r, &newReview) {
		http.Error(w, "Failed to create review", http.StatusBadRequest)
		return
	}
	_, err := db_controller.InsertReview(db_controller.Db.DB, newReview)
	if err != nil {
		http.Error(w, "Failed to create review", http.StatusBadRequest)
		return
	} else {
		log.Println("\nReview created")
	}
	json.NewEncoder(w).Encode("Review created successfully")
}

func ShowProductsBasedOnReview(w http.ResponseWriter, r *http.Request) {
	//show products based on review
	fmt.Println("Endpoint Hit: ShowProductsBasedOnReview")
	reqQueryStrings := r.URL.Query()
	var sorting string
	for key, value := range reqQueryStrings {
		if key == "value" {
			sorting = value[0]
		}
	}
	var products []db_controller.Product
	var err error
	if sorting == "worst" {
		products, err = db_controller.GetProductsByRatingAsc(db_controller.Db.DB)
		if err != nil {
			http.Error(w, "Failed to get products", http.StatusBadRequest)
			return
		}
	} else if sorting == "best" {
		products, err = db_controller.GetProductsByRatingDesc(db_controller.Db.DB)
		if err != nil {
			http.Error(w, "Failed to get products", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "Failed to get products. See if you choose \"worst\" or \"best\"", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(products)
}

func ShowReviewsForAProduct(w http.ResponseWriter, r *http.Request) {
	//show reviews for a product
	reqVariables := mux.Vars(r)
	productUuid, status := reqVariables["productUuid"]
	if !status {
		http.Error(w, "productUuid is missing", http.StatusBadRequest)
	}
	reviews, err := db_controller.GetReviewsForAProduct(db_controller.Db.DB, productUuid)
	if err != nil {
		http.Error(w, "Failed to get reviews", http.StatusBadRequest)
		return
	}
	if len(reviews) == 0 {
		http.Error(w, "No reviews for this product", http.StatusBadRequest)
		return
	} else {
		json.NewEncoder(w).Encode(reviews)
	}
}

func UpdateReview(w http.ResponseWriter, r *http.Request) {
	//update review
	reqVariables := mux.Vars(r)
	token, status := reqVariables["token"]
	if !status {
		http.Error(w, "token is missing", http.StatusBadRequest)
		return
	}
	productUuid, status := reqVariables["productUuid"]
	if !status {
		http.Error(w, "productUuid is missing", http.StatusBadRequest)
		return
	}
	var newReview db_controller.Review
	GetInfoReview(w, r, &newReview)

	username, err := db_controller.FindUserByToken(db_controller.Db.DB, token)
	if err != nil {
		http.Error(w, "Failed to find username", http.StatusInternalServerError)
		return
	}
	err = db_controller.UpdateReview(db_controller.Db.DB, username, productUuid, newReview.Description, strconv.Itoa(newReview.Rating))
	if err != nil {
		http.Error(w, "Failed to update review", http.StatusInternalServerError)
		return
	}

	// Send success response to client
	json.NewEncoder(w).Encode("review updated successfully")
}
