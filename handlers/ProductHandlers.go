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
	w.Write([]byte(productID))
}

func FetchProductsHandler(w http.ResponseWriter, r *http.Request) {
	pageNo := r.URL.Query().Get("pageNo")
	products := make([]models.ProductModel, 0)
	products, err := helper.FetchProductsHelper(pageNo)
	if err != nil {
		log.Printf("FetchProductsHandler : %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonResponse, err := utilities.JsonData(products)
	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	w.Write(jsonResponse)
}

func FetchProductsCategoryHandler(w http.ResponseWriter, r *http.Request) {
	pageNo := r.URL.Query().Get("pageNo")
	category := r.URL.Query().Get("category")
	products := make([]models.ProductModel, 0)
	products, err := helper.FetchProductsCategoryHelper(pageNo, category)
	if err != nil {
		log.Printf("FetchProductsHandler : %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonResponse, err := utilities.JsonData(products)
	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(jsonResponse)
}
func FetchLatestProductHandler(w http.ResponseWriter, r *http.Request) {
	product, err := helper.FetchLatestProductHelper()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonResponse, err := utilities.JsonData(product)
	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}
func FetchYouMayLike(w http.ResponseWriter, r *http.Request) {
	product, err := helper.FetchYouMayLikeHelper()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonResponse, err := utilities.JsonData(product)
	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}
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
	w.Write(jsonResponse)
}
