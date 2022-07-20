package handlers

import (
	"audioPhile/database"
	"audioPhile/database/helper"
	"audioPhile/middlewares"
	"audioPhile/models"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

func AddCartHandler(w http.ResponseWriter, r *http.Request) {
	signedUser := middlewares.UserFromContext(r.Context())
	var cart models.AddCartModel
	err := json.NewDecoder(r.Body).Decode(&cart)
	if err != nil {
		log.Printf("AddCartHandler: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cart.UserID = signedUser.ID
	var cartID string

	cart.TotalCost, err = helper.FetchTotalCostHelper(cart.Products)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	txErr := database.Tx(func(tx *sqlx.Tx) error {
		cartID, err := helper.AddCartHelper(cart, tx)
		if err != nil {
			return err
		}

		err = helper.AddCartProductsHelper(cartID, cart.Products, tx)
		return err
	})
	if txErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write([]byte(cartID))
	if err != nil {
		return
	}
}
func UpdateCartProductsHandler(w http.ResponseWriter, r *http.Request) {
	signedUser := middlewares.UserFromContext(r.Context())
	var cart models.UpdateCartModel
	err := json.NewDecoder(r.Body).Decode(&cart)
	if err != nil {
		log.Printf("UpdateCartProductsHandler: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cart.UserID = signedUser.ID

	txErr := database.Tx(func(tx *sqlx.Tx) error {
		err := helper.UpdateCartProductHelper(cart, tx)
		if err != nil {
			return err
		}

		err = helper.RemoveCartProductHelper(tx)
		if err != nil {
			return err
		}

		cart.TotalCost, err = helper.FetchTotalCostHelper(cart.Products)
		if err != nil {
			return err
		}

		err = helper.UpdateCartCostHelper(cart.CartID, cart.TotalCost, tx)

		return err
	})
	if txErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func AddCartProductHandler(w http.ResponseWriter, r *http.Request) {
	signedUser := middlewares.UserFromContext(r.Context())
	var cart models.UpdateCartModel
	err := json.NewDecoder(r.Body).Decode(&cart)
	if err != nil {
		log.Printf("AddCartProductHandler: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cart.UserID = signedUser.ID
	var cartID string

	cart.TotalCost, err = helper.FetchTotalCostHelper(cart.Products)

	txErr := database.Tx(func(tx *sqlx.Tx) error {
		err = helper.AddCartProductsHelper(cartID, cart.Products, tx)
		if err != nil {
			return err
		}

		cart.TotalCost, err = helper.FetchTotalCostHelper(cart.Products)
		if err != nil {
			return err
		}

		err = helper.UpdateCartCostHelper(cart.CartID, cart.TotalCost, tx)
		return err
	})
	if txErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write([]byte(cartID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
