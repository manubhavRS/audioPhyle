package handlers

import (
	"audioPhile/database"
	"audioPhile/database/helper"
	"audioPhile/models"
	"audioPhile/utilities"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

func AddCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var category models.CategoryModel
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		log.Printf("AddCategoryHandler: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	Fetchcategory, err := helper.AddCategoryHelper(category)
	if err != nil {
		log.Printf("AddCategoryHandler: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	jsonResponse, err := utilities.JsonData(Fetchcategory)
	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	_, err = w.Write(jsonResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func FetchCategoryHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := helper.FetchCategoryHelper()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonResponse, err := utilities.JsonData(categories)
	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	_, err = w.Write(jsonResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func AddProductHandler(w http.ResponseWriter, r *http.Request) {
	var product models.AddProductModel
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		log.Printf("AddProductHandler: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var productID string
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		productID, err := helper.AddProductHelper(product, tx)
		if err != nil {
			return err
		}

		err = helper.AddInTheBoxHelper(product.InTheBox, productID, tx)
		return err
	})
	if txErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write([]byte(productID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

//func FetchProductsHandler(w http.ResponseWriter, r *http.Request) {
//	pageNo := r.URL.Query().Get("pageNo")
//	products, err := helper.FetchProductsHelper(pageNo)
//	if err != nil {
//		log.Printf("FetchProductsHandler : %v", err)
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//
//	jsonResponse, err := utilities.JsonData(products)
//	if err != nil {
//		log.Printf("Error happened in JSON marshal. Err: %v", err)
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//
//	}
//	_, err = w.Write(jsonResponse)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//}

//func FetchProductsCategoryHandler(w http.ResponseWriter, r *http.Request) {
//	pageNo := r.URL.Query().Get("pageNo")
//	category := r.URL.Query().Get("category")
//	products := make([]models.ProductModel, 0)
//	products, err := helper.FetchProductsCategoryHelper(pageNo, category)
//	if err != nil {
//		log.Printf("FetchProductsHandler : %v", err)
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	jsonResponse, err := utilities.JsonData(products)
//	if err != nil {
//		log.Printf("Error happened in JSON marshal. Err: %v", err)
//		w.WriteHeader(http.StatusInternalServerError)
//	}
//	_, err = w.Write(jsonResponse)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//}
//func FetchLatestProductHandler(w http.ResponseWriter, r *http.Request) {
//	product, err := helper.FetchLatestProductHelper()
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	jsonResponse, err := utilities.JsonData(product)
//	if err != nil {
//		log.Printf("Error happened in JSON marshal. Err: %v", err)
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	_, err = w.Write(jsonResponse)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//}
//func FetchYouMayLike(w http.ResponseWriter, r *http.Request) {
//	product, err := helper.FetchYouMayLikeHelper()
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	jsonResponse, err := utilities.JsonData(product)
//	if err != nil {
//		log.Printf("Error happened in JSON marshal. Err: %v", err)
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	_, err = w.Write(jsonResponse)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//}
func FetchProductAssetsHandler(w http.ResponseWriter, r *http.Request) {
	var productID models.ProductID
	err := json.NewDecoder(r.Body).Decode(&productID)
	if err != nil {
		log.Printf("FetchProductAssetsHandler: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	productAssets, err := helper.FetchProductAssetHelper(productID.ID)
	if err != nil {
		log.Printf("FetchProductAssetsHandler : %v", err)
	}
	productAssetURLS := make([]models.ImageStructure, 0)
	for _, productAsset := range productAssets {
		productAssetURL := utilities.CreateImageUrl(productAsset)
		productAssetURLS = append(productAssetURLS, productAssetURL)
	}
	jsonResponse, err := utilities.JsonData(productAssetURLS)
	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

//func FetchProductsSearchHandler(w http.ResponseWriter, r *http.Request) {
//	var productSearch models.ProductSearchModel
//	pageNo := r.URL.Query().Get("pageNo")
//	productSearch.Category = r.URL.Query().Get("category")
//	productSearch.Search = r.URL.Query().Get("search")
//	limit := r.URL.Query().Get("limit")
//	productSearch.Latest = r.URL.Query().Get("latest")
//	productSearch.YML = r.URL.Query().Get("youMayLike")
//	if productSearch.Category == "" {
//		productSearch.CheckCategory = "true"
//	} else {
//		productSearch.CheckCategory = "false"
//	}
//	if productSearch.Latest == "true" && productSearch.YML == "true" {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	if productSearch.Latest == "true" {
//		productSearch.Limit = 1
//		productSearch.OrderBy = "created_by DESC"
//		productSearch.PageNo = 1
//	} else if productSearch.YML == "true" {
//		productSearch.Limit = 5
//		productSearch.OrderBy = "RANDOM ()"
//		productSearch.PageNo = 1
//	} else {
//		productSearch.PageNo, _ = strconv.Atoi(pageNo)
//		productSearch.Limit, _ = strconv.Atoi(limit)
//	}
//	products, err := helper.FetchProductsSearchHelper(productSearch)
//	if err != nil {
//		log.Printf("FetchProductsSearchHandler : %v", err)
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	jsonResponse, err := utilities.JsonData(products)
//	if err != nil {
//		log.Printf("Error happened in JSON marshal. Err: %v", err)
//		w.WriteHeader(http.StatusInternalServerError)
//	}
//	_, err = w.Write(jsonResponse)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//}

func FetchProductsListSearchHandler(w http.ResponseWriter, r *http.Request) {
	productSearch, err, ret := utilities.ArgumentsMapping(r)
	if err != nil {
		log.Printf("FetchProductsListSearchHandler : %v", err)
		w.WriteHeader(ret)
		return
	}
	products, err := helper.FetchProductsListSearchHelper(productSearch)
	if err != nil {
		log.Printf("FetchProductsListSearchHandler : %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonResponse, err := utilities.JsonData(products)
	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	_, err = w.Write(jsonResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func RemoveCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var category models.CategoryID
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		log.Printf("AddAddressHandler: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = helper.RemoveCategoryHelper(category)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func RemoveProductHandler(w http.ResponseWriter, r *http.Request) {
	var product models.ProductID
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		log.Printf("RemoveProductHandler: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = helper.RemoveProductHelper(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
