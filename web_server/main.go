package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// )

//reading the body of a post response, etc..

//required for http requests responses..

func main() {
	//get request takes an url
	//returns pointer to response and an error object
	//if err is nil, response contains response body, else failure

	// response, err := http.Get("https://url.com")
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Println(response)

	//post request contains a body requests stored as a binary
	// data is sent to server which sends response back

	// postBody, _ := json.Marshal(map[string]string{
	// 	"name":  "Mihnea",
	// 	"email": "mihneabejinaru@yahoo.com",
	// })
	// requestBody := bytes.NewBuffer(postBody)
	// response, err := http.Post("https://postman-echo.com/post", "application/json", requestBody)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	//postform used to post to a URL
	// URL encoded key value pairs

	//head
	// used to issue HEAD(metadata for request/response) to a URL
}
