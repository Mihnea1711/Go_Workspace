package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func fibo(terms int) []int {
	x := 0
	y := 1

	if terms <= 0 {
		return []int{}
	}

	if terms == 1 {
		return []int{x}
	}

	output := []int{x, y}
	for i := 0; i < terms-2; i++ {
		z := x + y
		output = append(output, z)
		x = y
		y = z
	}

	return output
}

func FibonacciHandler(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	fmt.Println(p)
	nrOfTerms, err := strconv.Atoi(p.ByName("nrOfTerms"))
	if err != nil {
		log.Fatal(err)
		return
	}
	fibonacciArray := fibo(nrOfTerms)
	fmt.Fprintf(w, "Your Fibonacci array of %v elements is: %v", nrOfTerms, fibonacciArray)
}

func main() {
	router := httprouter.New()

	router.GET("/fibonacci/:nrOfTerms/", FibonacciHandler)

	http.ListenAndServe(":8080", router)

	fmt.Print("Hi!")
}
