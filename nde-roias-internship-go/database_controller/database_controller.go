package database_controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Database_struct struct {
	DB *sql.DB
}

var Db Database_struct

func UseDB(dbName string) {
	_, err := Db.DB.Exec("USE " + dbName)
	if err != nil {
		log.Fatal(err)
	}
}

func DBworking() {
	//open db connexion
	var err error
	Db.DB, err = sql.Open("mysql", "root:ManzoC@841037@tcp(127.0.0.1:3306)/")
	if err != nil {
		log.Fatal(err)
	}

	//varify is a valid connexion
	err = Db.DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	//verify if the databases exists
	err = VerifyDB(Db.DB, "project_publicapi")
	if err != nil {
		log.Fatal(err)
	}

	VerifyAndCreateProductTable()

	//we need to create the other tables
	CreateOtherTables(Db.DB)

}

func VerifyDB(db *sql.DB, dbName string) error {
	//verifying the existence of the database
	var databaseExists string
	err := db.QueryRow("SHOW DATABASES LIKE '" + dbName + "'").Scan(&databaseExists)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("The database does not exist.")
			_, err = db.Exec("CREATE DATABASE " + dbName)
			if err != nil {
				log.Fatal(err)
				return err
			}
			fmt.Println("The database has been created successfully.")
		} else {
			log.Fatal(err)
			return err
		}
	}

	return err
}

func VerifyAndCreateProductTable() {
	UseDB("project_publicapi")
	tableName := "Product"
	tableExists := TableExists(tableName, "project_publicapi", Db.DB)
	if !tableExists {
		err := CreateProductTable(Db.DB, "Product", "project_publicapi")
		if err == nil {
			InsertInfoFromPublicAPI()
			fmt.Printf("Table %s was created succesffuly", tableName)
		}
	} else {
		fmt.Printf("Table %s already exists", tableName)
	}
}

func InsertInfoFromPublicAPI() {
	//we call the public API
	response, err := http.Get("https://dummyjson.com/products/")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject Response
	json.Unmarshal(responseData, &responseObject)

	for _, i := range responseObject.Products {
		InsertProduct(Db.DB, i)
	}
}

func TableExists(tableName string, dbName string, db *sql.DB) bool {
	var exists bool = false
	if err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_schema= ? AND table_name = ?)", dbName, tableName).Scan(&exists); err != nil {
		log.Fatalf("Failed to check if table exists: %v", err)
	}
	return exists
}

func MessageForCreateTables(err error, tableName string) {
	if err == nil {
		fmt.Printf("\nTable %s was created succesffuly", tableName)
	}
}

func CreateOtherTables(db *sql.DB) {
	dbName := "project_publicapi"
	UseDB(dbName)
	// Table User
	tableName := "User"
	tableExists := TableExists(tableName, dbName, Db.DB)
	if !tableExists {
		err := CreateUserTable(Db.DB, tableName, dbName)
		MessageForCreateTables(err, tableName)
	} else {
		fmt.Printf("\nTable %s already exists", tableName)
	}

	// Table ShoppingCart
	tableName = "ShoppingCart"
	tableExists = TableExists(tableName, dbName, Db.DB)
	if !tableExists {
		err := CreateShoppingCartTable(Db.DB, tableName, dbName)
		MessageForCreateTables(err, tableName)
	} else {
		fmt.Printf("\nTable %s already exists", tableName)
	}

	// Table Review
	tableName = "Review"
	tableExists = TableExists(tableName, dbName, Db.DB)
	if !tableExists {
		err := CreateReviewTable(Db.DB, tableName, dbName)

		MessageForCreateTables(err, tableName)
	} else {
		fmt.Printf("\nTable %s already exists", tableName)
	}
}

func CountRowsInTable(tableName string, db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func DB_close() {
	UseDB("project_publicapi")
	defer Db.DB.Close()
}
