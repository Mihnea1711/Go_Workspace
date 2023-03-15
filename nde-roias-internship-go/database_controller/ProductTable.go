package database_controller

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"

	uuid "github.com/google/uuid"
)

// //Structs
type Product struct {
	ID                 int     `json:"id"`
	Uuid               string  `json:"uuid"`
	Title              string  `json:"title"`
	Description        string  `json:"description"`
	Price              int     `json:"price"`
	DiscountPercentage float64 `json:"discountPercentage"`
	Rating             string  `json:"rating"`
	Stock              int     `json:"stock"`
	Quantity           int     `json:"quantity"`
	Brand              string  `json:"brand"`
	Category           string  `json:"category"`
	Thumbnail          string  `json:"thumbnail"`
}

// response because the details from Product are extracted from a public API
type Response struct {
	Products []Product `json:"products"`
}

// Functions

func CreateProductTable(db *sql.DB, dbTableName string, dbName string) error {

	UseDB("project_publicapi")
	query := `CREATE TABLE IF NOT EXISTS ` + dbTableName + ` (
					id int primary key auto_increment, 
					uuid varchar(36) NOT NULL UNIQUE,
					title text NOT NULL, 
					description text NOT NULL, 
					price int NOT NULL, 
					discountPercentage float, 
					rating text, 
					stock int NOT NULL, 
					quantity int NOT NULL,
					brand text NOT NULL, 
					category text NOT NULL, 
					thumbnail text NOT NULL
			)`
	_, err := db.Exec(query)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func GetLastProductID() (int, error) {
	UseDB("project_publicapi")
	query := `SELECT MAX(id) FROM product`

	lastIndex := 0
	queryRes, err := Db.DB.Query(query)
	if err != nil {
		return 0, errors.New("error in sql while trying to get last product id")
	}

	if queryRes.Next() {
		if err := queryRes.Scan(&lastIndex); err != nil {
			return 0, errors.New("error trying to scan last product id")
		}
	}

	return lastIndex, nil
}

// insert one product in the database
func InsertProduct(db *sql.DB, p Product) (Product, error) {
	UseDB("project_publicapi")
	query := "INSERT INTO Product (uuid, title, description, price, discountPercentage, rating, stock, quantity, brand, category, thumbnail) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	p.Uuid = uuid.New().String()
	p.Quantity = rand.Intn(100) + 1

	if _, err := db.Exec(query, p.Uuid, p.Title, p.Description, p.Price, p.DiscountPercentage, p.Rating, p.Stock, p.Quantity, p.Brand, p.Category, p.Thumbnail); err != nil {
		return Product{}, errors.New("error executing sql to insert new product")
	}

	id, err := GetLastProductID()
	if err != nil {
		return Product{}, err
	}
	p.ID = id

	return p, nil
}

// select all products inside database
func GetProducts(db *sql.DB) ([]Product, error) {
	UseDB("project_publicapi")
	query := "SELECT * FROM product"

	// query db to get all products
	queryRows, err := db.Query(query)
	if err != nil {
		return nil, errors.New("error in sql while trying to show all products")
	}

	// store products from db, otherwise return error
	var productList []Product
	for queryRows.Next() {
		var product Product
		if err := queryRows.Scan(&product.ID, &product.Uuid, &product.Title, &product.Description, &product.Price, &product.DiscountPercentage,
			&product.Rating, &product.Stock, &product.Quantity, &product.Brand, &product.Category, &product.Thumbnail); err != nil {
			return nil, errors.New("error trying to scan all the products")
		}
		productList = append(productList, product)
	}

	// return the product list if no error
	return productList, nil
}

func GetProduct(db *sql.DB, uuid string) (Product, error) {
	UseDB("project_publicapi")
	query := "SELECT * FROM product WHERE uuid = '" + uuid + "';"

	// query db to find product based on uuid
	queryRows, err := db.Query(query)
	if err != nil {
		return Product{}, errors.New("error in sql while trying to get the product based on uuid")
	}

	// store found product, otherwise return error
	var product Product
	if queryRows.Next() {
		if err := queryRows.Scan(&product.ID, &product.Uuid, &product.Title, &product.Description, &product.Price, &product.DiscountPercentage,
			&product.Rating, &product.Stock, &product.Quantity, &product.Brand, &product.Category, &product.Thumbnail); err != nil {
			return Product{}, errors.New("error trying to scan the query for a specific product")
		}
	}

	// return found product if no error
	return product, nil
}

func GetProductByUUID(db *sql.DB, uuid string) (Product, error) {
	var product Product
	// query to select product by id
	queryRows, err := db.Query("SELECT * FROM product WHERE uuid=?", uuid)
	if err != nil {
		return product, err
	}

	// loop through the rows and get the product
	for queryRows.Next() {
		err := queryRows.Scan(&product.ID, &product.Uuid, &product.Title, &product.Description, &product.Price, &product.DiscountPercentage,
			&product.Rating, &product.Stock, &product.Quantity, &product.Brand, &product.Category, &product.Thumbnail)
		if err != nil {
			log.Fatal(err)
		}
	}

	return product, nil
}

func GetProductsByProperty(db *sql.DB, productProperty string) ([]string, error) {
	UseDB("project_publicapi")
	query := "SELECT DISTINCT " + productProperty + " FROM product"

	// query db to find the property of a product
	queryRows, err := db.Query(query)
	if err != nil {
		return nil, errors.New("error in sql trying to get product property")
	}

	// store property, otherwise return error
	var propertyValues []string
	for queryRows.Next() {
		var propertyValue string
		if err := queryRows.Scan(&propertyValue); err != nil {
			return nil, errors.New("error trying to scan product property")
		}
		propertyValues = append(propertyValues, propertyValue)
	}

	// return property of product if no error
	return propertyValues, nil
}

func UpdateProduct(db *sql.DB, newProduct Product) (Product, error) {
	UseDB("project_publicapi")
	query := "UPDATE product SET title=?, description=?, price=?, discountPercentage=?, rating=?, stock=?, quantity=?, brand=?, category=?, thumbnail=? WHERE uuid=?"
	_, err := db.Exec(query, newProduct.Title, newProduct.Description, newProduct.Price, newProduct.DiscountPercentage, newProduct.Rating, newProduct.Stock, newProduct.Quantity, newProduct.Brand, newProduct.Category, newProduct.Thumbnail, newProduct.Uuid)
	if err != nil {
		return Product{}, errors.New("error in sql while updating the product")
	}

	// Return the updated product
	updatedProduct, err := GetProduct(db, newProduct.Uuid)
	if err != nil {
		return Product{}, err
	}

	return updatedProduct, nil
}

func DeleteProduct(db *sql.DB, uuid string) error {
	UseDB("project_publicapi")
	result, err := db.Exec("DELETE FROM product WHERE uuid = '" + uuid + "'")
	if err != nil {
		return errors.New("error in sql while deleting a product")
	}
	if rowsAffected, err := result.RowsAffected(); rowsAffected == 0 || err != nil {
		return errors.New("error deleting the product: product either doesn't exist or something bad happened")
	}
	return nil
}

// method to filter data in db with a (filter, filterValue) pair
func GetProductsByCategory(db *sql.DB, filter string, filterValue string) ([]Product, error) {
	// select products based on filter
	UseDB("project_publicapi")
	query := "SELECT * FROM product WHERE " + filter + " = '" + filterValue + "';"

	// query db to get all products based on a category
	queryRows, err := db.Query(query)
	if err != nil {
		return nil, errors.New("error trying to get the products by category")
	}

	// store products, otherwise return error
	var productList []Product
	for queryRows.Next() {
		var product Product
		if err := queryRows.Scan(&product.ID, &product.Uuid, &product.Title, &product.Description, &product.Price, &product.DiscountPercentage,
			&product.Rating, &product.Stock, &product.Quantity, &product.Brand, &product.Category, &product.Thumbnail); err != nil {
			return nil, errors.New("error trying to scan found products by category")
		}
		productList = append(productList, product)
	}

	// return product list if no error
	return productList, nil
}

func IncrementQuantity(db *sql.DB, productUuid string) (sql.Result, error) {
	UseDB("project_publicapi")

	// query to increment product quantity
	incrementQuantityQuery := `
	UPDATE product
	SET quantity = (
		SELECT quantity FROM product WHERE uuid = '` + productUuid + `'
		) + 1
	WHERE UUID = '` + productUuid + `'
	`
	// execute query to increment product qty
	result, err := db.Exec(incrementQuantityQuery)
	if err != nil {
		return nil, errors.New("something bad happened while trying to increment product quantity")
	}

	// return the sql result if no error
	return result, nil
}

func DecrementQuantity(db *sql.DB, productUuid string) (sql.Result, error) {
	UseDB("project_publicapi")

	// query to decrement product quantity
	decrementQuantityQuery := `
	UPDATE product
	SET quantity = (
		SELECT quantity FROM product WHERE uuid = '` + productUuid + `'
		) - 1
	WHERE UUID = '` + productUuid + `'
	`
	// execute query to decrement product qty
	result, err := db.Exec(decrementQuantityQuery)
	if err != nil {
		return nil, errors.New("something bad happened while trying to decrement product quantity")
	}

	// return sql result if no error
	return result, nil
}

func ChangeQtyByValue(db *sql.DB, productUuid string, qtyValue string) (sql.Result, error) {
	UseDB("project_publicapi")

	// query to decrement product quantity
	changeQuantityQuery := `
	UPDATE product
	SET quantity = (
		SELECT quantity FROM product WHERE uuid = '` + productUuid + `'
		) + '` + qtyValue + `'
	WHERE UUID = '` + productUuid + `'
	`
	// execute query to add/subtract product qty
	result, err := db.Exec(changeQuantityQuery)
	if err != nil {
		return nil, errors.New("something bad happened while trying to update product quantity")
	}

	// return sql result if no error
	return result, nil
}

func GetListProductUuid(db *sql.DB) ([]string, error) {
	var productUUIDs []string
	rows, err := db.Query("SELECT uuid FROM Product")
	if err != nil {
		return nil, errors.New("error trying to select uuid from Product")
	}
	defer rows.Close()
	for rows.Next() {
		var uuid string
		if err := rows.Scan(&uuid); err != nil {
			return nil, errors.New("error trying to scan uuid from Product")
		}
		productUUIDs = append(productUUIDs, uuid)
	}

	return productUUIDs, nil
}

func CreateRatingPerProduct(db *sql.DB) {
	productUUIDs, err := GetListProductUuid(db)
	if err != nil {
		log.Fatalf("Error in getting list of product uuids: %v", err)
	}

	for _, product_uuid := range productUUIDs {
		rows, err := db.Query("SELECT rating FROM Review WHERE product_uuid=?", product_uuid)
		if err != nil {
			log.Fatalf("Error in getting list of ratings: %v", err)
		}
		defer rows.Close()

		var totalRatings int
		var ratingSum float64
		for rows.Next() {
			var rating int
			if err := rows.Scan(&rating); err != nil {
				log.Fatalf("Error in scanning rating: %v", err)
			}
			ratingSum += float64(rating)
			totalRatings++
		}
		if err := rows.Err(); err != nil {
			log.Fatalf("Error in getting list of ratings: %v", err)
		}

		var productRating float64
		if totalRatings > 0 {
			productRating = ratingSum / float64(totalRatings)
		}

		var finalRating string
		if totalRatings == 0 {
			finalRating = "not rated"
		} else {
			finalRating = fmt.Sprintf("%.2f", productRating)
		}
		_, err = db.Exec("UPDATE Product SET rating=? WHERE uuid=?", finalRating, product_uuid)
		if err != nil {
			log.Fatalf("Error in updating rating: %v", err)
		}
	}
}

func GetProductsByRatingDesc(db *sql.DB) ([]Product, error) {
	query := `SELECT id, uuid, title, description, price, discountPercentage, rating, stock, quantity, brand, category, thumbnail FROM Product ORDER BY CASE WHEN rating = 'not rated' THEN 1 ELSE 0 END, rating DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Uuid, &product.Title, &product.Description, &product.Price, &product.DiscountPercentage, &product.Rating, &product.Stock, &product.Quantity, &product.Brand, &product.Category, &product.Thumbnail)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func GetProductsByRatingAsc(db *sql.DB) ([]Product, error) {
	query := `SELECT id, uuid, title, description, price, discountPercentage, rating, stock, quantity, brand, category, thumbnail FROM Product ORDER BY CASE WHEN rating = 'not rated' THEN 0 ELSE 1 END, rating ASC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Uuid, &product.Title, &product.Description, &product.Price, &product.DiscountPercentage, &product.Rating, &product.Stock, &product.Quantity, &product.Brand, &product.Category, &product.Thumbnail)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}
