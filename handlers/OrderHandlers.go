package handlers

import (
	"audioPhile/database"
	"audioPhile/database/helper"
	"audioPhile/middlewares"
	"audioPhile/models"
	"audioPhile/utilities"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

func AddOrderHandler(w http.ResponseWriter, r *http.Request) {
	var signedUser *models.UserModel
	signedUser = middlewares.UserFromContext(r.Context())
	var order models.AddOrderModel
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		log.Printf("AddOrderHandler: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	order.UserID = signedUser.ID
	cart, err := helper.FetchCartDetailsHelper(order.CartID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		order.DateOfDelivery = utilities.FetchExpectDateOfDelivery()
		order.CartDetails = cart
		err = helper.AddOrderHelper(order, tx)
		if err != nil {
			return err
		}
		err = helper.UpdateProductQuantity(cart, tx) //Bulk update using sqrl and sqlx.IN()
		if err != nil {
			return err
		}
		err = helper.RemoveCartHelper(order.CartID, tx)
		return err
	})
	if txErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
