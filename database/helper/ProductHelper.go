package helper

import (
	"audioPhile/database"
	"audioPhile/models"
	"github.com/elgris/sqrl"
	"github.com/jmoiron/sqlx"
	"log"
	"strconv"
)

func AddProductHelper(product models.AddProductModel, tx *sqlx.Tx) (string, error) {
	//language=SQL
	SQL := `INSERT INTO products(name,category,price,feature,about,quantity) 
			values($1,$2,$3,$4,$5,$6) 
			Returning id`
	var productID string
	err := tx.Get(&productID, SQL, product.Name, product.Category, product.Price, product.Feature, product.About, product.Quantity)
	if err != nil {
		log.Printf("AddProductHelper Error: %v", err)
		return "", err
	}
	return productID, nil
}
func AddInTheBoxHelper(inTheBoxList []models.AddInTheBoxModel, productID string, tx *sqlx.Tx) error {
	//language=SQL
	psql := sqrl.StatementBuilder.PlaceholderFormat(sqrl.Dollar)
	insertBuilder := psql.Insert("in_the_box").Columns("pdt_id", "name", "quantity")
	for _, inTheBox := range inTheBoxList {
		insertBuilder.Values(productID, inTheBox.Name, inTheBox.Quantity)
	}
	sql, args, err := insertBuilder.ToSql()
	if err != nil {
		log.Printf("AddInTheBoxHelper Error: %v", err)
		return err
	}
	_, err = tx.Exec(sql, args...)
	if err != nil {
		log.Printf("AddInTheBoxHelper Error: %v", err)
		return err
	}
	return nil
}
func FetchProductsHelper(pageNo string) ([]models.ProductModel, error) {
	//language=SQL
	SQL := `SELECT id,name,category,price,about,quantity
          FROM products
          WHERE archived_at IS NULL
          LIMIT 10
          OFFSET $1
          `
	products := make([]models.ProductModel, 0)
	page, err := strconv.ParseInt(pageNo, 10, 64)
	if err != nil {
		log.Printf("FetchProductsHelper Error: %v", err)
		return nil, err
	}
	err = database.Aph.Select(&products, SQL, (page-1)*10)
	if err != nil {
		log.Printf("FetchProductsHelper Error: %v", err)
		return nil, err
	}
	return products, nil
}
func FetchProductsCategoryHelper(pageNo, category string) ([]models.ProductModel, error) {
	//language=SQL
	SQL := `SELECT id,name,category,price,about,quantity
          FROM products
          WHERE category=$1 and
          archived_at IS NULL
          LIMIT 10
          OFFSET $2`
	products := make([]models.ProductModel, 0)
	page, err := strconv.ParseInt(pageNo, 10, 64)
	if err != nil {
		log.Printf("FetchProductsHelper Error: %v", err)
		return nil, err
	}
	err = database.Aph.Select(&products, SQL, category, (page-1)*10)
	if err != nil {
		log.Printf("FetchProductsHelper Error: %v", err)
		return nil, err
	}
	return products, nil
}
func UpdateProductQuantity(cart models.CartModel, tx *sqlx.Tx) error {
	//language=SQL
	SQL := `SELECT quantity 
			from products p 
			where p.id in (?) and
          	archived_at IS NULL`
	quantity := make([]int, 0)
	productIDs := make([]string, 0)
	for _, product := range cart.Products {
		productIDs = append(productIDs, product.ProductID)
	}
	sqlStatement, args, err := sqlx.In(SQL, productIDs)
	sqlStatement = database.Aph.Rebind(sqlStatement)
	err = database.Aph.Select(&quantity, sqlStatement, args...)
	if err != nil {
		log.Printf("UpdateProductQuantity: %v", err)
		return err
	}
	psql := sqrl.StatementBuilder
	insertBuilder := psql.Update("products")
	insertCase := sqrl.Case("id")
	for i, productID := range productIDs {
		insertCase.When(`'`+productID+`'`, strconv.Itoa(quantity[i]-cart.Products[i].Quantity))
	}
	insertBuilder.Set("quantity", insertCase)
	insertBuilder.Where("id IN (?)")
	sql, args, err := insertBuilder.ToSql()
	if err != nil {
		log.Printf("UpdateProductQuantity Error: %v", err)
		return err
	}
	sql, args, err = sqlx.In(sql, productIDs)
	sql = database.Aph.Rebind(sql)
	_, err = tx.Exec(sql, args...)
	if err != nil {
		log.Printf("UpdateProductQuantity Error: %v", err)
		return err
	}
	return nil
}
func FetchLatestProductHelper() (models.ProductModel, error) {
	//language=SQL
	SQL := `SELECT id, name, category,price,about,quantity 
			FROM products
            WHERE archived_at IS NULL
			ORDER BY created_at DESC                                        
			LIMIT 1`
	var product models.ProductModel
	err := database.Aph.Get(&product, SQL)
	if err != nil {
		log.Printf("FetchLatestProductHelper Error: %v", err)
		return product, err
	}
	return product, nil
}
func FetchYouMayLikeHelper() ([]models.ProductModel, error) {
	//language=SQL
	SQL := `SELECT id, name, category,price,about,quantity 
			FROM products
			WHERE archived_at IS NULL
			ORDER BY RANDOM ()                                    
			LIMIT 5`
	product := make([]models.ProductModel, 0)
	err := database.Aph.Select(&product, SQL)
	if err != nil {
		log.Printf("FetchYouMayLikeHelper Error: %v", err)
		return product, err
	}
	return product, nil
}
func AddProductAssetHelper(productAsset models.ProductAssetModel, tx *sqlx.Tx) (string, error) {
	//language=SQL
	SQL := `INSERT INTO product_assets(pdt_id,name) 
		  VALUES($1,$2)
		  RETURNING id`
	var id string
	err := tx.Get(&id, SQL, productAsset.ID, productAsset.AssetName)
	if err != nil {
		log.Printf("AddProductAssetHelper Error: %v", err)
		return "", err
	}
	return id, nil
}
func FetchProductAssetHelper(productID string) ([]string, error) {
	//language=SQL
	SQL := `SELECT name 
		  FROM product_assets 
          WHERE pdt_id=$1 AND archived_at IS NULL`
	names := make([]string, 0)
	err := database.Aph.Select(&names, SQL, productID)
	if err != nil {
		log.Printf("FetchProductAssetHelper Error: %v", err)
		return names, err
	}
	return names, nil
}
