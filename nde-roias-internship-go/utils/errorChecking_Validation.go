package error_checking

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	db_controller "project3/demo/database_controller"
)

func CheckFormProductCreation(r *http.Request) (db_controller.Product, error) {
	title := r.FormValue("title")
	description := r.FormValue("description")
	price, err := strconv.Atoi(r.FormValue("price"))
	if err != nil {
		return db_controller.Product{}, errors.New("price is not valid")
	}
	discount, err := strconv.ParseFloat(r.FormValue("discountPercentage"), 64)
	if err != nil {
		return db_controller.Product{}, errors.New("discount is not valid")
	}
	rating := r.FormValue("rating")
	stock, err := strconv.Atoi(r.FormValue("stock"))
	if err != nil {
		return db_controller.Product{}, errors.New("stock is not valid")
	}
	quantity, err := strconv.Atoi(r.FormValue("quantity"))
	if err != nil {
		return db_controller.Product{}, errors.New("quantity is not valid")
	}
	brand := r.FormValue("brand")
	category := r.FormValue("category")
	thumbnail := r.FormValue("thumbnail")

	// error if something is missing
	errorMessage := ""
	if title == "" {
		errorMessage = "Title is required."
	} else if description == "" {
		errorMessage = "Description is required."
	} else if price <= 0 {
		errorMessage = "Price is required and must be a price greater than 0."
	} else if stock <= 0 {
		errorMessage = "Stock is required and must be a number greater than 0."
	} else if quantity <= 0 {
		errorMessage = "Quantity is required and must be a number greater than 0."
	} else if brand == "" {
		errorMessage = "Brand is required."
	} else if category == "" {
		errorMessage = "Category is required."
	} else if thumbnail == "" {
		errorMessage = "Thumbnail is required."
	}

	if errorMessage != "" {
		return db_controller.Product{}, errors.New(errorMessage)
	}

	var product = db_controller.Product{
		Title:              title,
		Description:        description,
		Price:              price,
		DiscountPercentage: discount,
		Rating:             rating,
		Stock:              stock,
		Quantity:           quantity,
		Brand:              brand,
		Category:           category,
		Thumbnail:          thumbnail,
	}

	return product, nil
}

func CheckProductUUID(uuid string) error {
	// check if product exists
	db_controller.UseDB("project_publicapi")
	query := "SELECT 2 FROM product WHERE uuid = '" + uuid + "';"

	var product interface{}
	if err := db_controller.Db.DB.QueryRow(query).Scan(&product); err == nil {
		return nil
	} else if err == sql.ErrNoRows {
		return errors.New("// product not found in the database. Insert a valid product UUID \\")
	} else {
		return errors.New("// an error happened while trying to find the product in the database \\ Error: " + err.Error())
	}
}

func CheckUserID(id int) error {
	// check if user exists
	db_controller.UseDB("project_publicapi")
	query := "SELECT 2 FROM user WHERE id = " + strconv.Itoa(id) + ";"

	var user interface{}
	if err := db_controller.Db.DB.QueryRow(query).Scan(&user); err == nil {
		return nil
	} else if err == sql.ErrNoRows {
		return errors.New("// user not found in the database. Insert a valid user ID \\")
	} else {
		return errors.New("// an error happened while trying to find the user in the database \\ Error: " + err.Error())
	}
}

func CheckShoppingCartUUID(productUUID string, userID int) error {
	// check if cart item exists in user's cart
	db_controller.UseDB("project_publicapi")
	query := "SELECT 2 FROM shoppingcart WHERE product_uuid = '" + productUUID + "' and user_id = " + strconv.Itoa(userID) + ";"

	var cartItem interface{}
	if err := db_controller.Db.DB.QueryRow(query).Scan(&cartItem); err == nil {
		return nil
	} else if err == sql.ErrNoRows {
		return errors.New("// product not found in the user's shopping cart. Insert a valid user ID or a valid item UUID \\ Error: " + err.Error())
	} else {
		return errors.New("// an error happened while trying to find the shopping item in the user's cart \\ Error: " + err.Error())
	}
}

func CheckProductAvailability(productUUID string, qty string) error {
	if productExistErr := CheckProductUUID(productUUID); productExistErr != nil {
		return productExistErr
	}

	var enough bool
	// check if quantity is > 0
	if err := db_controller.Db.DB.QueryRow("SELECT (quantity >= " + qty + ") FROM product WHERE uuid = '" + productUUID + "';").Scan(&enough); err != nil || !enough {
		if err == sql.ErrNoRows {
			return errors.New("product not existing")
		}
		return errors.New("product out of stock")
	}

	return nil
}

func CheckItemAddition(productUUID string, userID int) error {
	sql := `
		SELECT (sc.item_quantity + 1 < p.quantity)
		FROM product p JOIN shoppingcart sc ON p.uuid = sc.product_uuid
		WHERE sc.product_uuid = '` + productUUID + `' and sc.user_id = ` + strconv.Itoa(userID)

	var enough bool
	// check if we have enough stock
	if err := db_controller.Db.DB.QueryRow(sql).Scan(&enough); err != nil || !enough {
		return errors.New("no more products in stock or an error has occured while trying to add 1 more item")
	}

	return nil
}

func CheckItemDeletion(productUUID string, userID int) error {
	sql := `
		SELECT (sc.item_quantity - 1 > 0)
		FROM shoppingcart sc
		WHERE sc.product_uuid = '` + productUUID + `' and sc.user_id = ` + strconv.Itoa(userID)

	var enough bool
	// check if we have enough stock
	if err := db_controller.Db.DB.QueryRow(sql).Scan(&enough); err != nil {
		return errors.New("error trying to access db in check item deletion")
	}

	if !enough {
		// delete product from cart
		if result, err := db_controller.RemoveProductFromCart(db_controller.Db.DB, productUUID, userID); err != nil {
			return err
		} else {
			if rowsAffected, err := result.RowsAffected(); err != nil || rowsAffected == 0 {
				return errors.New("something bad happened trying to delete item from cart.. no changes were made")
			}
		}
		return errors.New("item was deleted from cart")
	}

	return nil
}

func CheckProductQTY(productUUID, qty string) error {
	query := `
	SELECT (quantity + ` + qty + ` > 0)
	FROM product
	WHERE UUID = '` + productUUID + "'"

	var enough bool
	// check if we have enough stock
	if err := db_controller.Db.DB.QueryRow(query).Scan(&enough); err != nil {
		return errors.New("error trying to access db in check item deletion")
	}

	if !enough {
		// delete product from cart
		if err := db_controller.DeleteProduct(db_controller.Db.DB, productUUID); err != nil {
			return err
		}
		return errors.New("item was deleted from stock")
	}

	return nil
}
