package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	db_controller "project3/demo/database_controller"
	err_check "project3/demo/utils"
	"strconv"

	"github.com/gorilla/mux"
)

func AddToCart(w http.ResponseWriter, r *http.Request) {
	reqVariables := mux.Vars(r)

	// get product from path
	productUUID, flagUUID := reqVariables["productUuid"]
	if !flagUUID {
		json.NewEncoder(w).Encode("Enter a valid product uuid")
		return
	}

	// get user id from path
	userID, err := strconv.Atoi(reqVariables["userID"])
	if err != nil {
		// render an error template for more info about the problem
		json.NewEncoder(w).Encode("// User ID has not a valid format. Please retry... \\ Error: " + err.Error())
		return
	}

	// check if product exists
	if err := err_check.CheckProductUUID(productUUID); err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// check userUUID
	if err := err_check.CheckUserID(userID); err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// validate quantity for the productUUID
	itemCurrentQTY, err := db_controller.GetCartProductQTY(productUUID, userID)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	if err := err_check.CheckProductAvailability(productUUID, strconv.Itoa(itemCurrentQTY+1)); err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	if err := err_check.CheckShoppingCartUUID(productUUID, userID); err == nil {
		_, err := db_controller.IncreaseItemQTY(db_controller.Db.DB, productUUID, userID)
		if err != nil {
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		json.NewEncoder(w).Encode("product already in cart. item qty increased")
		return
	}

	// if product and user exists, add product to user's cart
	response, err := db_controller.AddProductInCart(db_controller.Db.DB, productUUID, userID)
	if err != nil || response == nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// output
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Product added successfully")
}

func RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	reqVariables := mux.Vars(r)

	// get product uuid from path
	productUUID, flagUUID := reqVariables["productUuid"]
	if !flagUUID {
		json.NewEncoder(w).Encode(errors.New("please insert a valid product uuid"))
	}

	// get user id from path
	userID, err := strconv.Atoi(reqVariables["userID"])
	if err != nil {
		// render an error template for more info about the problem
		json.NewEncoder(w).Encode("// User ID has not a valid format. Please retry... \\ Error: " + err.Error())
		return
	}

	// check if the product is in the user's cart
	if err = err_check.CheckShoppingCartUUID(productUUID, userID); err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// remove the product from user's cart
	response, err := db_controller.RemoveProductFromCart(db_controller.Db.DB, productUUID, userID)
	if err != nil || response == nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// output
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Product removed successfully")
}

func SeeShoppingCart(w http.ResponseWriter, r *http.Request) {
	reqVariables := mux.Vars(r)

	// get user id from path
	userID, err := strconv.Atoi(reqVariables["userID"])
	if err != nil {
		json.NewEncoder(w).Encode("// User ID has not a valid format. Please retry... \\ Error: " + err.Error())
		return
	}

	// check userUUID
	if err = err_check.CheckUserID(userID); err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// get the products inside the user's cart
	products, err := db_controller.SeeUserCart(db_controller.Db.DB, userID)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// output
	w.Header().Set("Content-Type", "application/json")
	if products == nil {
		// maybe the user has no items inside
		json.NewEncoder(w).Encode("No products found...")
		return
	}
	json.NewEncoder(w).Encode(products)
}

func IncreaseItemQTY(w http.ResponseWriter, r *http.Request) {
	reqVariables := mux.Vars(r)

	// get user id from path
	userID, err := strconv.Atoi(reqVariables["userID"])
	if err != nil {
		json.NewEncoder(w).Encode(errors.New("enter a valid user uuid"))
		return
	}

	// check userUUID
	if err := err_check.CheckUserID(userID); err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// check product uuid
	productUuid, flagUUID := reqVariables["productUuid"]
	if !flagUUID {
		json.NewEncoder(w).Encode(errors.New("enter a valid product uuid"))
		return
	}

	// check if the product is in the user's cart
	if err := err_check.CheckShoppingCartUUID(productUuid, userID); err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// check if the product has enough stock to add 1 more to the cart
	if err := err_check.CheckItemAddition(productUuid, userID); err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// add 1 more item in cart
	response, err := db_controller.IncreaseItemQTY(db_controller.Db.DB, productUuid, userID)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	if rowsAffected, err := response.RowsAffected(); rowsAffected == 0 || err != nil {
		json.NewEncoder(w).Encode("error happened and nothing has been changed")
		return
	}

	// output
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("added 1 more product")
}

func DecreaseItemQTY(w http.ResponseWriter, r *http.Request) {
	// check and verify userUUID
	reqVariables := mux.Vars(r)

	// get user id from path
	userID, err := strconv.Atoi(reqVariables["userID"])
	if err != nil {
		json.NewEncoder(w).Encode("// User ID has not a valid format. Please retry... \\ Error: " + err.Error())
		return
	}

	// check user id
	if err := err_check.CheckUserID(userID); err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// check and verify product uuid
	productUuid, flagUUID := reqVariables["productUuid"]
	if !flagUUID {
		json.NewEncoder(w).Encode(errors.New("enter a valid product uuid"))
		return
	}

	// check & validate productUuid
	if err := err_check.CheckShoppingCartUUID(productUuid, userID); err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// check if the product has enough stock to add 1 more to the cart
	if err := err_check.CheckItemDeletion(productUuid, userID); err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// decrease item quantity in cart
	response, err := db_controller.DecreaseItemQTY(db_controller.Db.DB, productUuid, userID)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	if rowsAffected, err := response.RowsAffected(); rowsAffected == 0 || err != nil {
		json.NewEncoder(w).Encode("no rows were modified whule trying to decrease item qty from cart")
		return
	}

	// output
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("removed 1 product")
}

func PlaceOrder(w http.ResponseWriter, r *http.Request) {
	// check and verify userUUID
	reqVariables := mux.Vars(r)
	userID, err := strconv.Atoi(reqVariables["userID"])
	if err != nil {
		json.NewEncoder(w).Encode("// User ID has not a valid format. Please retry... \\ Error: " + err.Error())
		return
	}

	// find user shopping cart items
	products, err := db_controller.SeeUserCart(db_controller.Db.DB, userID)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// store validated products in an array to show at checkout
	// and an array of possible outdated/out-of-stock items to show them to the user in case anything goes wrong
	var validProducts, invalidProducts []db_controller.ShoppingCart

	// store total price, number of products to print at the end alongside the list
	var totalPrice, nrOfItems int

	// for each item from cart, check for left quantity
	var errorMsg = ""
	for _, item := range products {
		// this should be the original item's uuid :/
		err := err_check.CheckProductAvailability(item.ProductUuid, strconv.Itoa(item.ItemQuantity))
		if err != nil {
			errorMsg += err.Error() + " \\ "
			invalidProducts = append(invalidProducts, item)
			continue
		}

		validProducts = append(validProducts, item)

		totalPrice += item.ProductPrice * item.ItemQuantity
		nrOfItems += item.ItemQuantity
	}

	// decrement quantity of the products if everything goes ok
	// check the invalidProducts list first! Should be an atomic operation
	// if the invalid products list has products in it, then send encode products and send them without doin anything
	w.Header().Set("Content-Type", "application/json")
	if invalidProducts != nil {
		json.NewEncoder(w).Encode("Error: " + errorMsg)
		json.NewEncoder(w).Encode("the following products are invalid so the order could not be placed")
		json.NewEncoder(w).Encode(invalidProducts)
		return
	}

	for _, product := range validProducts {
		db_controller.DecrementQuantity(db_controller.Db.DB, product.ProductUuid)
		db_controller.RemoveProductFromCart(db_controller.Db.DB, product.ProductUuid, userID)
	}

	// display json with the valid product list and total price
	json.NewEncoder(w).Encode("You have successfully bought the following items\n")
	json.NewEncoder(w).Encode(validProducts)

	json.NewEncoder(w).Encode("Total price: " + strconv.Itoa(totalPrice))
	json.NewEncoder(w).Encode("Number Of Items: " + strconv.Itoa(nrOfItems))
}
