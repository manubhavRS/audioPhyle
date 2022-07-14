package handlers

import (
	"audioPhile/database/helper"
	"audioPhile/middlewares"
	"audioPhile/models"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func AddCardHandler(w http.ResponseWriter, r *http.Request) {
	var signedUser *models.UserModel
	signedUser = middlewares.UserFromContext(r.Context())
	var card models.AddCardModel
	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		log.Printf("AddCardHandler: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	card.UserID = signedUser.ID
	if len(card.CardNumber) > 12 {
		log.Printf("AddCardHandler: Invalid Card Number")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	expDate, err := time.Parse("01-06", card.ExpireDate)
	if err != nil {
		log.Printf("AddCardHandler: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if time.Now().Year() >= expDate.Year() && time.Now().Month() >= expDate.Month() {
		log.Printf("AddCardHandler: Card Expired")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	card.Expiry = expDate
	cardID, err := helper.AddCardHelper(card)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(cardID))
}
func RemoveCardHandler(w http.ResponseWriter, r *http.Request) {
	var signedUser *models.UserModel
	signedUser = middlewares.UserFromContext(r.Context())
	var card models.RemoveCardIDModel
	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		log.Printf("AddAddressHandler: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	card.UserID = signedUser.ID
	err = helper.RemoveCardHelper(card)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
