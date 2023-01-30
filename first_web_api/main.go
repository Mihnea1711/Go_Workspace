package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	funcUtils "example.com/first_web_api/func_utils"

	"github.com/julienschmidt/httprouter"
)

type Embed struct {
	Message string
	Id      int
}

type Output interface{}

func HomepageHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	templ, err := template.ParseFiles("html/home.html")
	if err != nil {
		log.Fatal(err)
		return
	}
	templ.Execute(w, nil)
}

func InputHandler(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	id, err := strconv.Atoi(p.ByName("id"))
	if err != nil {
		log.Fatal(err)
		return
	}
	templ, err := template.ParseFiles("html/input.tmpl")
	if err != nil {
		log.Fatal(err)
		return
	}
	message := ""
	switch id {
	case 1:
		message = "Enter the array so I can sort it..."
	case 2:
		message = "Enter the array so I can find the duplicates..."
	case 3:
		message = "Enter the array so I can remove the duplicates..."
	case 4:
		message = "Enter the string so I can find the longest palindromic substring..."
	case 5:
		message = "Enter the array so I can find the missing number..."
	}
	templ.Execute(w, &Embed{Message: message, Id: id})
}

func FibonacciHandler(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	nrOfTerms, err := strconv.Atoi(p.ByName("nrOfTerms"))
	if err != nil {
		log.Fatal(err)
		return
	}
	fibonacciArray := funcUtils.Fibo(nrOfTerms)
	fmt.Fprintf(w, "Your Fibonacci array of %v elements is: %v", nrOfTerms, fibonacciArray)
}

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

func OutputHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, err := strconv.Atoi(p.ByName("id"))
	if err != nil {
		log.Fatal(err)
		return
	}
	r.ParseForm()
	inputData := r.Form["inputData"]
	// fmt.Fprintln(w, "Data: ", inputData)
	// fmt.Fprintln(w, "Option: ", option)

	templ, err := template.ParseFiles("html/output.tmpl")
	if err != nil {
		log.Fatal(err)
		return
	}

	var output Output

	switch id {
	case 1:
		//merge
		output = funcUtils.SortArray(stringArrToIntArray(inputData))
		// fmt.Fprintf(w, "Data: %v", output)
	case 2:
		//merge
		output = funcUtils.GetDups(stringArrToIntArray(inputData))
		// fmt.Fprintf(w, "Data: %v", output)
	case 3:
		//merge
		output = funcUtils.ElimDups(stringArrToIntArray(inputData))
		// fmt.Fprintf(w, "Data: %v", output)
	case 4:
		//ciudat.. la input: "dada" => in loc sa dea "ada", da "dadadad"
		output = funcUtils.GetLongestSubstring(inputData[0])
		// fmt.Fprintf(w, "Data: %v", output)
	case 5:
		//merge
		output = funcUtils.GetMissingNr(stringArrToIntArray(inputData))
		// fmt.Fprintf(w, "Data: %v", output)
	}

	templ.Execute(w, output)
}

func route() {
	router := httprouter.New()

	router.GET("/functions", HomepageHandler)
	router.GET("/functions/:id", InputHandler)
	router.POST("/functions/:id", OutputHandler)

	//different method using input params in URL
	router.GET("/fibonacci/:nrOfTerms/", FibonacciHandler)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	route()
}
