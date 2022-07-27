package helper

import (
	"audioPhile/database"
	"audioPhile/models"
	"audioPhile/utilities"
	"github.com/elgris/sqrl"
	"github.com/jmoiron/sqlx"
	"log"
	"strconv"
)

func FetchTotalCostHelper(products []models.AddCartProductModel) (float64, error) {
	//language=SQL
	sqlStatement := `SELECT price 
					 from products p 
					 where p.id 
					 in (?) and archived_at IS NULL`
	costs := make([]float64, 0)
	productsIDs := make([]string, 0)
	for _, product := range products {
		productsIDs = append(productsIDs, product.ProductID)
	}
	sqlStatement, args, err := sqlx.In(sqlStatement, productsIDs)
	sqlStatement = database.Aph.Rebind(sqlStatement)
	err = database.Aph.Select(&costs, sqlStatement, args...)
	if err != nil {
		log.Printf("FetchTotalCostHelper: %v", err)
		return 0, err
	}
	var totalCost float64
	for i, cost := range costs {
		totalCost = totalCost + ((float64(products[i].Quantity)) * cost)
	}
	return totalCost, err
}
func AddCartHelper(cart models.AddCartModel, tx *sqlx.Tx) (string, error) {
	//language=SQL
	SQL := `INSERT INTO carts(user_id,total_cost) 
			values($1,$2) 
			RETURNING id`
	var cartID string
	err := tx.Get(&cartID, SQL, cart.UserID, cart.TotalCost)
	if err != nil {
		log.Printf("AddCartHelper Error: %v", err)
		return "", err
	}
	log.Printf(cartID)
	return cartID, nil
}
func AddCartProductsHelper(cartID string, products []models.AddCartProductModel, tx *sqlx.Tx) error {
	psql := sqrl.StatementBuilder.PlaceholderFormat(sqrl.Dollar)
	insertBuilder := psql.Insert("cart_products").Columns("product_id", "cart_id", "quantity")
	for _, product := range products {
		insertBuilder.Values(product.ProductID, cartID, product.Quantity)
	}
	sql, args, err := insertBuilder.ToSql()
	if err != nil {
		log.Printf("AddCartProductsHelper Error: %v", err)
		return err
	}
	_, err = tx.Exec(sql, args...)
	if err != nil {
		log.Printf("AddCartProductsHelper Error: %v", err)
		return err
	}
	return nil
}
func FetchCartDetailsHelper(cartID string) (models.CartModel, error) {
	//language=SQL
	SQL := `SELECT cp.product_id,cp.quantity,p.price
			FROM cart_products cp 
			INNER JOIN products p
			on cp.product_id=p.id
			WHERE cart_id=$1 AND 
            cp.archived_at IS NULL`
	var cartDetails models.CartModel
	cartProducts := make([]models.AddCartProductModel, 0)
	err := database.Aph.Select(&cartProducts, SQL, cartID)
	if err != nil {
		log.Printf("FetchCartDetailsHelper Error: %v ", err)
		return cartDetails, err
	}
	SQL = `SELECT total_cost 
		   from carts 
			where id=$1 and
			archived_at IS NULL`
	err = database.Aph.Get(&cartDetails.TotalCostWithShipping, SQL, cartID)
	if err != nil {
		log.Printf("FetchCartDetailsHelper Error: %v ", err)
		return cartDetails, err
	}
	var num int
	cartDetails.Products = cartProducts
	for _, cartProduct := range cartDetails.Products {
		num = num + cartProduct.Quantity
	}
	cartDetails.TotalCostWithShipping = cartDetails.TotalCostWithShipping + float64(num*utilities.ShippingCharges)
	cartDetails.Products = cartProducts
	return cartDetails, nil
}
func UpdateCartProductHelper(cart models.UpdateCartModel, tx *sqlx.Tx) error {
	productIDs := make([]string, 0)
	for _, product := range cart.Products {
		productIDs = append(productIDs, product.ProductID)
	}
	psql := sqrl.StatementBuilder
	insertBuilder := psql.Update("cart_products")
	insertCase := sqrl.Case("product_id")
	for i, productID := range productIDs {
		insertCase.When(`'`+productID+`'`, strconv.Itoa(cart.Products[i].Quantity))
	}
	insertBuilder.Set("quantity", insertCase)
	insertBuilder.Where("cart_id=" + `'` + cart.CartID + `'`)
	insertBuilder.Where("product_id in (?)")
	insertBuilder.Where("archived_at IS NULL")
	sql, args, err := insertBuilder.ToSql()
	if err != nil {
		log.Printf("UpdateCartProductHelper Error: %v", err)
		return err
	}
	sql, args, err = sqlx.In(sql, productIDs)
	sql = database.Aph.Rebind(sql)
	_, err = tx.Exec(sql, args...)
	if err != nil {
		log.Printf("UpdateCartProductHelper Error: %v", err)
		return err
	}
	return nil
}
func RemoveCartProductHelper(tx *sqlx.Tx) error {
	//language=SQL
	SQL := `Update cart_products 
			SET archived_at=CURRENT_TIMESTAMP
			WHERE quantity=0 and
			archived_at IS NULL`
	_, err := tx.Exec(SQL)
	if err != nil {
		log.Printf("RemoveCartProductHelper Error: %v", err)
		return err
	}
	return nil
}
func UpdateCartCostHelper(cartID string, cost float64, tx *sqlx.Tx) error {
	//language=SQL
	SQL := `UPDATE carts 
		  SET total_cost=$1 
		  WHERE id=$2 and
		  archived_at IS NULL`
	_, err := tx.Exec(SQL, cost, cartID)
	if err != nil {
		log.Printf("UpdateCartCostHelper Error: %v", err)
		return err
	}
	return nil
}
func RemoveCartHelper(cartID string, tx *sqlx.Tx) error {
	//language=SQL
	SQL := `UPDATE carts
		   SET archived_at=CURRENT_TIMESTAMP
		   WHERE id=$1`
	_, err := tx.Exec(SQL, cartID)
	if err != nil {
		log.Printf("RemoveCartHelper Error: %v", err)
		return err
	}
	return nil
}
