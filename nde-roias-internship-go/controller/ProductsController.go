package controller

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	db_controller "project3/demo/database_controller"
	error_checking "project3/demo/utils"

	"github.com/gorilla/mux"
)

// for postman request
type AllProducts []db_controller.Product

var ProductsArr = AllProducts{}

// in next function we get all info from postman for Product table
func GetInfoProduct(w http.ResponseWriter, r *http.Request, newProduct *db_controller.Product) {
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(reqBody, newProduct) //not &newArray because is already pointer

	ProductsArr = append(ProductsArr, *newProduct)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newProduct)
}

// get route for showing all products in db
func ShowProducts(w http.ResponseWriter, r *http.Request) {
	// get all products
	products, err := db_controller.GetProducts(db_controller.Db.DB)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// post route for creating a new project and adding it to the db
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	// store from request values in the product
	formProduct, err := error_checking.CheckFormProductCreation(r)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
	}

	// insert product
	insertedProduct, err := db_controller.InsertProduct(db_controller.Db.DB, formProduct)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Successfully added 1 product")
	json.NewEncoder(w).Encode(insertedProduct)
}

// get route for showing specific product in db
func ViewProduct(w http.ResponseWriter, r *http.Request) {
	reqVariables := mux.Vars(r)

	// get product uuid from path
	productUuid, status := reqVariables["productUuid"]
	if !status {
		json.NewEncoder(w).Encode("The productUuid parameter is missing from the request")
		return
	}

	// check if product exists
	if err := error_checking.CheckProductUUID(productUuid); err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// gert product based on uuid
	product, err := db_controller.GetProduct(db_controller.Db.DB, productUuid)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func ViewProductsByProperty(w http.ResponseWriter, r *http.Request) {
	reqVariables := mux.Vars(r)

	// get product property from path
	productProperty, statusProperty := reqVariables["property"]
	if !statusProperty {
		json.NewEncoder(w).Encode("The property parameter is missing from the request")
		return
	}

	// get product property from db
	propertyValues, err := db_controller.GetProductsByProperty(db_controller.Db.DB, productProperty)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// output
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(propertyValues)
}

// put route for updating a product in db // make it with postman, not browser
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	reqVariables := mux.Vars(r)
	product, err := error_checking.CheckFormProductCreation(r)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	productUuid, flag := reqVariables["productUuid"]
	if !flag {
		json.NewEncoder(w).Encode(errors.New("insert a valid uuid"))
		return
	}

	product.Uuid = productUuid

	// update product in db
	w.Header().Set("Content-Type", "application/json")
	updatedProduct, err := db_controller.UpdateProduct(db_controller.Db.DB, product)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
	}
	json.NewEncoder(w).Encode("Product updated successfully")
	json.NewEncoder(w).Encode(updatedProduct)
}

// delete route for deleting a specific route from db
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	reqVariables := mux.Vars(r)
	productUuid, status := reqVariables["productUuid"]
	if !status {
		json.NewEncoder(w).Encode("The productUuid parameter is missing from the request")
	}

	err := db_controller.DeleteProduct(db_controller.Db.DB, productUuid)
	if err != nil {
		RespondWithJSON(w, http.StatusInternalServerError, err.Error())
	} else {
		RespondWithJSON(w, http.StatusOK, "Product deleted succesfully")
	}
}

// callback function for filtering data based on a query string
func FilterByCategory(w http.ResponseWriter, r *http.Request) {
	reqQueryStrings := r.URL.Query()

	var filter string
	var filterValue string

	// we assume ther is only one filter for simplicity (otherwise, we filter the data by the last filter used)
	for key, value := range reqQueryStrings {
		filter = key
		filterValue = value[0] // assume the filter has only one value
	}

	// filter products and store in an array
	products, err := db_controller.GetProductsByCategory(db_controller.Db.DB, filter, filterValue)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// output
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func IncrementQuantity(w http.ResponseWriter, r *http.Request) {
	reqVariables := mux.Vars(r)

	// get product uuid
	productUuid, flagUUID := reqVariables["productUuid"]
	if !flagUUID {
		json.NewEncoder(w).Encode(errors.New("insert a valid product uuid"))
		return
	}

	// check if product exists
	productExistErr := error_checking.CheckProductUUID(productUuid)
	if productExistErr != nil {
		json.NewEncoder(w).Encode(productExistErr.Error())
		return
	}

	// if product exists
	queryResult, queryErr := db_controller.IncrementQuantity(db_controller.Db.DB, productUuid)
	if queryErr != nil || queryResult == nil {
		json.NewEncoder(w).Encode(queryErr.Error())
		return
	}

	// output
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("product qty has been incremented successfully")
}

func DecrementQuantity(w http.ResponseWriter, r *http.Request) {
	reqVariables := mux.Vars(r)

	// get product uuid
	productUuid, flagUUID := reqVariables["productUuid"]
	if !flagUUID {
		json.NewEncoder(w).Encode(errors.New("insert a valid product uuid"))
		return
	}

	// check if product exists
	productExistErr := error_checking.CheckProductUUID(productUuid)
	if productExistErr != nil {
		json.NewEncoder(w).Encode(productExistErr.Error())
		return
	}

	if err := error_checking.CheckProductQTY(productUuid, "-1"); err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// if product exists
	queryResult, queryErr := db_controller.DecrementQuantity(db_controller.Db.DB, productUuid)
	if queryErr != nil || queryResult == nil {
		json.NewEncoder(w).Encode(queryErr.Error())
		return
	}

	// output
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("product qty has been decremented successfully")
}

func ChangeQtyByValue(w http.ResponseWriter, r *http.Request) {
	reqVariables := mux.Vars(r)

	// query string
	reqQueryStrings := r.URL.Query()
	var qtyValue string

	// we take the value only from the value query and ignore others
	for key, value := range reqQueryStrings {
		if key == "value" {
			qtyValue = value[0]
		}
	}

	// get product uuid
	productUuid, flagUUID := reqVariables["productUuid"]
	if !flagUUID {
		json.NewEncoder(w).Encode("insert a valid product uuid")
		return
	}

	// check if product exists
	productExistErr := error_checking.CheckProductUUID(productUuid)
	if productExistErr != nil {
		json.NewEncoder(w).Encode(productExistErr.Error())
		return
	}

	if err := error_checking.CheckProductQTY(productUuid, qtyValue); err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// if product exists
	queryResult, queryErr := db_controller.ChangeQtyByValue(db_controller.Db.DB, productUuid, qtyValue)
	if queryErr != nil || queryResult == nil {
		json.NewEncoder(w).Encode(queryErr.Error())
		return
	}

	// output
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("product qty has been updated successfully")
}

func RespondWithJSON(w http.ResponseWriter, statusCode int, message string) {
	result, _ := json.Marshal(message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(result)
}
