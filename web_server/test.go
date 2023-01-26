package main

import "net/http"

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, reader *http.Request) {
		http.ServeFile(writer, reader, "html/unrecognized_path.html")
	})

	http.HandleFunc("/demopath1/", func(writer http.ResponseWriter, reader *http.Request) {
		http.ServeFile(writer, reader, "html/demopath1.html")
	})

	http.HandleFunc("/demopath1/subpatha", func(writer http.ResponseWriter, reader *http.Request) {
		http.ServeFile(writer, reader, "html/demopath1_subpatha.html")
	})

	http.HandleFunc("/demopath2", func(writer http.ResponseWriter, reader *http.Request) {
		http.ServeFile(writer, reader, "html/demopath2.html")
	})

	http.ListenAndServe(":8080", nil)
}
