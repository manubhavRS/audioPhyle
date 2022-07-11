package handlers

import (
	"audioPhile/database/helper"
	"audioPhile/middlewares"
	"audioPhile/models"
	"audioPhile/utilities"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func SignUpUserHandler(w http.ResponseWriter, r *http.Request) {
	var addUser models.AddUserModel
	err := json.NewDecoder(r.Body).Decode(&addUser)
	if err != nil {
		log.Printf("SignUpUserHandler: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userID, err := helper.SignUpUserHelper(addUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonData, jsonErr := json.Marshal(userID)
	if jsonErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(jsonData)
}

func SignInUserHandler(w http.ResponseWriter, r *http.Request) {
	var credentials models.SignInUserModel
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		log.Printf("SignInUserHandler: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := helper.FetchUserCredentialsHelper(credentials.Email)
	ok := utilities.CheckPasswordHash(credentials.Password, user.Password)
	if !ok {
		log.Printf("Password Incorrect.")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	tokenString, err := middlewares.GenerateJWT(&user)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["Name"] = "token"
	resp["Value"] = tokenString
	resp["Expires"] = (time.Now().Add(30 * time.Minute)).String()
	jsonResponse, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %v", err)
	}
	w.Write(jsonResponse)
}
