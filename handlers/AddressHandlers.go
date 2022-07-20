package handlers

import (
	"audioPhile/database/helper"
	"audioPhile/middlewares"
	"audioPhile/models"
	"encoding/json"
	"log"
	"net/http"
)

func AddAddressHandler(w http.ResponseWriter, r *http.Request) {
	signedUser := middlewares.UserFromContext(r.Context())
	var address models.AddAddressModel
	err := json.NewDecoder(r.Body).Decode(&address)
	if err != nil {
		log.Printf("AddAddressHandler: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	address.UserID = signedUser.ID

	cardID, err := helper.AddAddressHelper(address)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(cardID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func RemoveAddressHandler(w http.ResponseWriter, r *http.Request) {
	signedUser := middlewares.UserFromContext(r.Context())
	var address models.RemoveAddressIDModel
	err := json.NewDecoder(r.Body).Decode(&address)
	if err != nil {
		log.Printf("AddAddressHandler: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	address.UserID = signedUser.ID

	err = helper.RemoveAddressHelper(address)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
