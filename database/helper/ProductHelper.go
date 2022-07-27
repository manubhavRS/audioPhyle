package helper

import (
	"audioPhile/database"
	"audioPhile/models"
	"audioPhile/utilities"
	"github.com/elgris/sqrl"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"log"
	"strconv"
)

func AddCategoryHelper(category models.CategoryModel) (models.FetchCategoryModel, error) {
	//language=SQL
	SQL := `INSERT INTO categories(category)
			values($1) 
			Returning id`
	var Fetchcategory models.FetchCategoryModel
	err := database.Aph.Get(&Fetchcategory.ID, SQL, category.Category)
	if err != nil {
		log.Printf("AddCategoryHelper Error: %v", err)
		return Fetchcategory, err
	}
	Fetchcategory.Category = category.Category

	return Fetchcategory, nil
}
func FetchCategoryHelper() ([]models.FetchCategoryModel, error) {
	//language=SQL
	SQL := `SELECT id,category 
			FROM categories 
            WHERE archived_at IS NULL`
	categoryList := make([]models.FetchCategoryModel, 0)
	err := database.Aph.Select(&categoryList, SQL)
	if err != nil {
		log.Printf("FetchCategoryHelper Error: %v", err)
		return categoryList, err
	}
	return categoryList, nil
}

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
func FetchProductsHelper(pageNo string) (models.ProductListModel, error) {
	//language=SQL
	SQL := `WITH cte AS (
    		SELECT *
    		FROM   products
			)
			SELECT id,name,category,price,about,quantity,full_count
			FROM  (
          		TABLE  cte
              	ORDER  BY id
              	LIMIT  10
              	OFFSET $1
      		) sub
          	RIGHT  JOIN (SELECT count(*) FROM cte) c(full_count) ON true;`
	products := make([]models.ProductModelWithCount, 0)
	pdtList := make([]models.ProductModel, 0)
	var productsList models.ProductListModel
	page, err := strconv.ParseInt(pageNo, 10, 64)
	if err != nil {
		log.Printf("FetchProductsHelper Error: %v", err)
		return productsList, err
	}
	err = database.Aph.Select(&products, SQL, (page-1)*10)
	if err != nil {
		log.Printf("FetchProductsHelper Error: %v", err)
		return productsList, err
	}
	var totalCount int
	for _, pdt := range products {
		var product models.ProductModel
		product.ID = pdt.ID
		product.Name = pdt.Name
		product.Price = pdt.Price
		product.About = pdt.About
		product.Category = pdt.Category
		product.Quantity = pdt.Quantity
		pdtList = append(pdtList, product)
		totalCount = pdt.TotalCount
	}
	productsList.TotalCount = totalCount
	productsList.Products = pdtList
	return productsList, nil
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

func FetchProductsSearchHelper(productSearch models.ProductSearchModel) (models.ProductListModel, error) {
	//language=SQL
	SQL := `WITH cte AS (
    		SELECT *
    		FROM   products
    		WHERE ($1 OR category= $2)  AND
          	name ILIKE '%' || $3 || '%' AND 
        	archived_at IS NULL
			)
			SELECT id,name,category,price,about,quantity,full_count
			FROM  (
          		TABLE  cte
              	ORDER  BY $4
              	LIMIT  $5
              	OFFSET $6
      		) sub
          	RIGHT  JOIN (SELECT count(*) FROM products) c(full_count) ON true;`
	products := make([]models.ProductModelWithCount, 0)
	pdtList := make([]models.ProductModel, 0)
	var productsList models.ProductListModel
	err := database.Aph.Select(&products, SQL, productSearch.CheckCategory, utilities.NewNullString(productSearch.Category), productSearch.Search, productSearch.OrderBy, productSearch.Limit, (productSearch.PageNo-1)*10)
	if err != nil {
		log.Printf("FetchProductsSearchHelper Error: %v", err)
		return productsList, err
	}
	var totalCount int
	for _, pdt := range products {
		var product models.ProductModel
		product.ID = pdt.ID
		product.Name = pdt.Name
		product.Price = pdt.Price
		product.About = pdt.About
		product.Category = pdt.Category
		product.Quantity = pdt.Quantity
		pdtList = append(pdtList, product)
		totalCount = pdt.TotalCount
	}
	productsList.TotalCount = totalCount
	productsList.Products = pdtList
	return productsList, nil
}
func FetchProductsListSearchHelper(productSearch models.ProductListSearchModel) (models.ProductListModel, error) {
	//language=SQL
	SQL := `WITH cte AS (
    		SELECT *
    		FROM   products
    		WHERE ($1 OR category=ANY($2))  AND
          	name ILIKE '%' || $3 || '%' AND 
        	archived_at IS NULL
			)
			SELECT id,name,category,price,about,quantity,full_count
			FROM  (
          		TABLE  cte
              	ORDER  BY $4
              	LIMIT  $5
              	OFFSET $6
      		) sub
          	RIGHT  JOIN (SELECT count(*) FROM products) c(full_count) ON true;`
	products := make([]models.ProductModelWithCount, 0)
	pdtList := make([]models.ProductModel, 0)
	var productsList models.ProductListModel
	err := database.Aph.Select(&products, SQL, productSearch.CheckCategory, pq.StringArray(productSearch.Category), productSearch.Search, productSearch.OrderBy, productSearch.Limit, (productSearch.PageNo-1)*10)
	if err != nil {
		log.Printf("FetchProductsListSearchHelper Error: %v", err)
		return productsList, err
	}
	var totalCount int
	for _, pdt := range products {
		var product models.ProductModel
		product.ID = pdt.ID
		product.Name = pdt.Name
		product.Price = pdt.Price
		product.About = pdt.About
		product.Category = pdt.Category
		product.Quantity = pdt.Quantity
		pdtList = append(pdtList, product)
		totalCount = pdt.TotalCount
	}
	productsList.TotalCount = totalCount
	productsList.Products = pdtList
	return productsList, nil
}
func RemoveCategoryHelper(category models.CategoryID) error {
	//language=SQL
	SQL := `UPDATE categories
  		    SET archived_at=CURRENT_TIMESTAMP
  		    WHERE id=$1 
  		    RETURNING id`
	_, err := database.Aph.Exec(SQL, category.ID)
	if err != nil {
		log.Printf("RemoveCategoryHelper Error: %v", err)
		return err
	}
	return err
}
func RemoveProductHelper(product models.ProductID) error {
	//language=SQL
	SQL := `UPDATE products
  		    SET archived_at=CURRENT_TIMESTAMP
  		    WHERE id=$1 
  		    RETURNING id`
	_, err := database.Aph.Exec(SQL, product.ID)
	if err != nil {
		log.Printf("RemoveProductHelper Error: %v", err)
		return err
	}
	return err
}
