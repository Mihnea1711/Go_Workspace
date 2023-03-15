package main

import (
	init_app "project3/demo/init_app"

	_ "github.com/go-sql-driver/mysql"
)

// Initiate web server
func main() {

	init_app.Init__main()
}
