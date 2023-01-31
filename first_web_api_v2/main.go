package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	funcUtils "rest_api.com/demo/func_utils"

	"github.com/gorilla/mux"
)

// this is a struct
type Response struct {
	Output interface{}
}

// helper function
func stringArrToIntArray(inputArr []string) []int {
	outputArr := []int{}
	for _, character := range strings.Split(inputArr[0], ", ") {
		integer, err := strconv.Atoi(character)
		if err != nil {
			log.Fatal(err)
			return []int{}
		}
		outputArr = append(outputArr, integer)
	}
	return outputArr
}

// home route
func home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("html/home.html")
	if err != nil {
		log.Fatal(err)
		return
	}
	tmpl.Execute(w, nil)
}

// fibonacci handler
func handleFibonacci(w http.ResponseWriter, r *http.Request) {
	reqVariables := mux.Vars(r)
	fiboArrayLength, err := strconv.Atoi(reqVariables["length"])
	if err != nil {
		log.Fatal(err)
		return
	}

	fiboArray := funcUtils.Fibo(fiboArrayLength)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// fmt.Fprintln(w, "Change array size from URL params!")
	// fmt.Fprintf(w, "Fibonacci array of size %v is: %v\n", fiboArrayLength, fiboArray)

	json.NewEncoder(w).Encode(Response{Output: fiboArray})
}

// renders a form which dynamically redirects the input to the specific handler
func renderForm(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	//one way
	// templateString := "<div style='text-align: center;'><h1>Input the array!</h1><form action='{{.}}' method='POST'><input type='text' name='inputData'><button>Submit</button></form></div>"
	// tp := template.New("Input Template")
	// tp, _ = tp.Parse(templateString)
	// tp.Execute(w, path)

	//another way
	templ, _ := template.ParseFiles("html/input.html")
	templ.Execute(w, path)
}

// sort handler
func handleSort(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	inputArray := r.Form["inputData"]

	sorted := funcUtils.SortArray(stringArrToIntArray(inputArray))

	fmt.Fprintf(w, "Input array is: %v\n", inputArray)
	fmt.Fprintf(w, "Sorted array is: %v\n", sorted)

	json.NewEncoder(w).Encode(Response{Output: sorted})
}

// duplicate handler ( showing the duplicates of an array)
func handleGetDuplicates(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	inputArray := r.Form["inputData"]

	dict := funcUtils.GetDups(stringArrToIntArray(inputArray))

	fmt.Fprintf(w, "Input array is: %v\n", inputArray)
	fmt.Fprintf(w, "Duplicates are: \n")
	for key, value := range dict {
		if value > 1 {
			fmt.Fprintf(w, "Elementul %v apare de %v ori\n", key, value)
		}
	}

	json.NewEncoder(w).Encode(Response{Output: dict})
}

// duplicate handler (removing the duplicates from an array)
func handleRemoveDuplicates(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	inputArray := r.Form["inputData"]

	arraySet := funcUtils.ElimDups(stringArrToIntArray(inputArray))

	fmt.Fprintf(w, "Input array is: %v\n", inputArray)
	fmt.Fprintf(w, "Array without duplicates is: %v\n", arraySet)

	json.NewEncoder(w).Encode(Response{Output: arraySet})
}

// substring handler (finds the longest palindromic substring from a given string)
func handlesubstring(w http.ResponseWriter, r *http.Request) {
	reqVariables := mux.Vars(r)
	inputString := reqVariables["string"]

	substring := funcUtils.GetLongestSubstring(inputString)

	fmt.Fprintln(w, "Change the input string from URL params!")
	fmt.Fprintf(w, "Longest palindromic substring of %v is: %v\n", inputString, substring)

	json.NewEncoder(w).Encode(Response{Output: substring})
}

// missing number handler (finds the missing number from an array starting from 0)
func handleFindMissing(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	inputArray := r.Form["inputData"]

	missing := funcUtils.GetMissingNr(stringArrToIntArray(inputArray))

	fmt.Fprintf(w, "Input array is: %v\n", inputArray)
	fmt.Fprintf(w, "Missing number is: %v\n", missing)

	json.NewEncoder(w).Encode(Response{Output: missing})
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	//endpoint for selecting the functions
	router.HandleFunc("/", home)

	//endpoint for the fibonacci handler
	//input the size of the fibonacci array in the URL path
	router.HandleFunc("/fibonacci/{length}", handleFibonacci).Methods("GET")

	//route for the sort handler
	router.HandleFunc("/sort", renderForm).Methods("GET")
	router.HandleFunc("/sort", handleSort).Methods("POST")

	//route for the first duplicates handler
	router.HandleFunc("/get-duplicates", renderForm).Methods("GET")
	router.HandleFunc("/get-duplicates", handleGetDuplicates).Methods("POST")

	//route for the second duplicates handler
	router.HandleFunc("/remove-duplicates", renderForm).Methods("GET")
	router.HandleFunc("/remove-duplicates", handleRemoveDuplicates).Methods("POST")

	//route for the longest palindromic substring handler
	router.HandleFunc("/longest-substring/{string}", handlesubstring).Methods("GET")

	//route for the missing number handler
	router.HandleFunc("/find-missing", renderForm).Methods("GET")
	router.HandleFunc("/find-missing", handleFindMissing).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	handleRequests()
}
