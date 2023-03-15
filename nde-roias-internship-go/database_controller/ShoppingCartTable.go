package database_controller

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
)

type ShoppingCart struct {
	ID              int
	ProductUuid     string
	ProductName     string
	ProductDetails  string
	ProductPrice    int
	ProductDiscount float64
	ItemQuantity    int
	UserID          int
}

func CreateShoppingCartTable(db *sql.DB, tableName string, dbName string) error {
	UseDB(dbName)
	query := `CREATE TABLE IF NOT EXISTS ` + tableName + ` (
					id int primary key auto_increment, 
					product_uuid varchar(36) NOT NULL,
					product_name text NOT NULL,
					product_details text NOT NULL,
					product_price int NOT NULL,
					product_discount float NOT NULL,
					item_quantity int NOT NULL,
					user_id int NOT NULL,
					FOREIGN KEY (user_id) REFERENCES User(id) ON DELETE CASCADE
			)`
	_, err := db.Exec(query)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

// method to add a product in the user's cart
func AddProductInCart(db *sql.DB, productUUID string, userID int) (sql.Result, error) {
	//check if product exists and return product to store properties
	product, err := GetProduct(db, productUUID)
	if err != nil {
		return nil, err
	}

	// when adding an item in the cart, start with quantity 1
	itemQuantity := 1

	// insert into shoppin cart the product properties and the user's id
	// nu mai stau sa fac o functie care ia user id-ul doar pt ca e gresit tabelul
	UseDB("project_publicapi")
	query := "INSERT INTO shoppingCart (product_uuid, product_name, product_details, product_price, product_discount, item_quantity, user_id) VALUES ('" +
		productUUID + "', '" + product.Title + "', '" +
		product.Description + "', " + strconv.Itoa(product.Price) + ", '" +
		strconv.FormatFloat(product.DiscountPercentage, 'g', 3, 64) + "', " + strconv.Itoa(itemQuantity) + ", " + strconv.Itoa(userID) + ")"

	// query the db to add a product in cart
	response, err := db.Exec(query)
	if err != nil {
		return nil, errors.New("error trying to execute query")
	}

	// if error or no rows affected, return err
	rowsAffected, err := response.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return nil, errors.New("an error occured. item was not added in cart")
	}

	// return sql response if no error
	return response, nil
}

// method to remove an item from the cart based on its uuid and on the user id
func RemoveProductFromCart(db *sql.DB, productUUID string, userID int) (sql.Result, error) {
	// select and delete the corresponding item
	UseDB("project_publicapi")
	query := "DELETE FROM shoppingCart WHERE product_uuid = '" + productUUID + "' and user_id = " + strconv.Itoa(userID) + ";"

	// query the db to remove an item from the cart
	response, err := db.Exec(query)
	if err != nil {
		return nil, errors.New("error in sql trying to delete the product")
	}

	// if error or no rows affected, return err
	rowsAffected, err := response.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return nil, errors.New("product does not exist so it cannot be deleted")
	}

	// return sql response if no error
	return response, nil
}

// method to see the user's shopping cart
func SeeUserCart(db *sql.DB, userID int) ([]ShoppingCart, error) {
	UseDB("project_publicapi")
	query := "SELECT * FROM shoppingcart WHERE user_id = " + strconv.Itoa(userID) + ";"

	// query db to see all items in user's cart
	queryRows, err := Db.DB.Query(query)
	if err != nil {
		return nil, errors.New("error while trying to access user's cart")
	}

	// store items found, otherwise return error
	var cartItems []ShoppingCart
	for queryRows.Next() {
		var cartItem ShoppingCart
		err := queryRows.Scan(&cartItem.ID, &cartItem.ProductUuid, &cartItem.ProductName, &cartItem.ProductDetails, &cartItem.ProductPrice,
			&cartItem.ProductDiscount, &cartItem.ItemQuantity, &cartItem.UserID)
		if err != nil {
			return nil, errors.New("error trying to scan user's cart")
		}
		cartItems = append(cartItems, cartItem)
	}

	// return products from user's cart if no error
	return cartItems, nil
}

func IncreaseItemQTY(db *sql.DB, productUuid string, userID int) (sql.Result, error) {
	UseDB("project_publicapi")
	query := `
		UPDATE shoppingcart 
		SET item_quantity = (
			SELECT item_quantity FROM shoppingcart WHERE product_uuid = '` + productUuid + `' AND user_id = ` + strconv.Itoa(userID) + `
		) + 1
		WHERE product_uuid = '` + productUuid + `' AND user_id = ` + strconv.Itoa(userID)

	// query db to increment qty of a product in cart, otherwise return error
	result, err := db.Exec(query)
	if err != nil {
		return nil, errors.New("error in sql while adding 1 more product")
	} else {
		if rowsAffected, err := result.RowsAffected(); err != nil || rowsAffected == 0 {
			return nil, errors.New("an error occured while tring to increase item qty in cart.. nothing has changed")
		}
		return result, nil
	}
}

func DecreaseItemQTY(db *sql.DB, productUuid string, userID int) (sql.Result, error) {
	UseDB("project_publicapi")

	query := `
		UPDATE shoppingcart 
		SET item_quantity = (
			SELECT item_quantity FROM shoppingcart WHERE product_uuid = '` + productUuid + `' AND user_id = ` + strconv.Itoa(userID) + `
		) - 1
		WHERE product_uuid = '` + productUuid + `' AND user_id = ` + strconv.Itoa(userID)

	// query db to idecrement qty of a product in cart, otherwise return error
	result, err := db.Exec(query)
	if err != nil {
		return nil, errors.New("error in sql while adding 1 more product")
	} else {
		if rowsAffected, err := result.RowsAffected(); err != nil || rowsAffected == 0 {
			return nil, errors.New("an error occured while tring to increase item qty in cart.. nothing has changed")
		}
		return result, nil
	}
}

func GetCartProductQTY(productUUID string, userID int) (int, error) {
	UseDB("project_publicapi")
	query := `SELECT item_quantity FROM shoppingcart WHERE product_uuid = '` + productUUID + `' and user_id = ` + strconv.Itoa(userID)

	// query db to see all items in user's cart
	queryRows, err := Db.DB.Query(query)
	if err != nil {
		return 0, errors.New("error while trying to get cart item qty")
	}

	var itemQty = 0
	if queryRows.Next() {
		err := queryRows.Scan(&itemQty)
		if err != nil {
			return 0, errors.New("error trying to scan item qty")
		}
	}

	return itemQty, nil
}
