package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	funcUtils "rest_api.com/demo/func_utils"

	"github.com/gorilla/mux"
)

func home(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/all", http.StatusSeeOther)
}

func getAllFunctions(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "All functions endpoint hit")
}
func handleChoosenFunction(w http.ResponseWriter, r *http.Request) {
	reqVariables := mux.Vars(r)
	id := reqVariables["id"]

	fmt.Fprintf(w, "Id is: %v", id)
}

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	reqVariables := mux.Vars(r)
	id := reqVariables["id"]

	reqBody, _ := io.ReadAll(r.Body)

	fmt.Fprintf(w, "Selected function id is: %v\n", id)
	fmt.Fprintf(w, "Request body id is: %v\n", string(reqBody))
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", home)
	router.HandleFunc("/all", getAllFunctions).Methods("GET")
	router.HandleFunc("/all/{id}", handleChoosenFunction).Methods("GET")
	router.HandleFunc("/all/{id}", handleSubmit).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	fmt.Println("Hi!")
	fmt.Println(funcUtils.Fibo(5))

	handleRequests()
}
