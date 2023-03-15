package init_app

import (
	"log"
	"net/http"
	db_controller "project3/demo/database_controller"
	r "project3/demo/routes"
)

func DB_init() {
	db_controller.DBworking()
}

func Init_server() {
	router := r.Router()
	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8080",
	}

	log.Fatal(srv.ListenAndServe())
}

func Init__main() {
	DB_init()

	Init_server()

	db_controller.DB_close()
}
