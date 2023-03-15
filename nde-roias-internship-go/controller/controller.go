package controller

import (
	"net/http"
)

// handle / home page. Temporarily just redirect to /products
func HomePage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/products", http.StatusSeeOther)
}
