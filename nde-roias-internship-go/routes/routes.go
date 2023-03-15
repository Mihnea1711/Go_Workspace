package routes

import (
	controller "project3/demo/controller"

	"github.com/gorilla/mux"
)

// Route declarations
func Router() *mux.Router {
	r := mux.NewRouter() //we create a router

	r.HandleFunc("/", controller.HomePage).Methods("GET")
	r.HandleFunc("/products", controller.ShowProducts).Methods("GET")
	r.HandleFunc("/products", controller.CreateProduct).Methods("POST")
	r.HandleFunc("/products/filter", controller.FilterByCategory).Methods("GET")
	r.HandleFunc("/products/reviews", controller.ShowProductsBasedOnReview).Methods("GET")
	r.HandleFunc("/products/{productUuid}", controller.ViewProduct).Methods("GET")
	r.HandleFunc("/products/{productUuid}", controller.UpdateProduct).Methods("PUT")
	r.HandleFunc("/products/{productUuid}", controller.DeleteProduct).Methods("DELETE")
	r.HandleFunc("/products/{productUuid}/incrementQuantity", controller.IncrementQuantity).Methods("POST")
	r.HandleFunc("/products/{productUuid}/decrementQuantity", controller.DecrementQuantity).Methods("POST")
	r.HandleFunc("/products/{productUuid}/changeProductQTY", controller.ChangeQtyByValue).Methods("POST")
	r.HandleFunc("/searchProperties/{property}", controller.ViewProductsByProperty).Methods("GET")

	r.HandleFunc("/users/register", controller.RegisterUser).Methods("POST")
	r.HandleFunc("/users/login", controller.LoginUser).Methods("POST")
	r.HandleFunc("/users/{userID}/cart", controller.SeeShoppingCart).Methods("GET")
	r.HandleFunc("/users/{userID}/cart/{productUuid}/increaseCartItemQTY", controller.IncreaseItemQTY).Methods("POST")
	r.HandleFunc("/users/{userID}/cart/{productUuid}/decreaseCartItemQTY", controller.DecreaseItemQTY).Methods("POST")
	r.HandleFunc("/users/{userID}/cart/{productUuid}/addToCart", controller.AddToCart).Methods("POST")
	r.HandleFunc("/users/{userID}/cart/{productUuid}/removeFromCart", controller.RemoveFromCart).Methods("POST")
	r.HandleFunc("/users/{userID}/placeOrder", controller.PlaceOrder).Methods("POST")

	r.HandleFunc("/review/{productUuid}", controller.CreateReview).Methods("POST")
	r.HandleFunc("/review/{productUuid}", controller.ShowReviewsForAProduct).Methods("GET")
	r.HandleFunc("/review/{token}/{productUuid}", controller.UpdateReview).Methods("PUT")

	return r
}
