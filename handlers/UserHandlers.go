package handlers

import (
	"audioPhile/database/helper"
	"audioPhile/middlewares"
	"audioPhile/models"
	"audioPhile/utilities"
	"database/sql"
	"encoding/json"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
)

func AdminSignUpUserHandler(w http.ResponseWriter, r *http.Request) {
	var addUser models.AddUserModel
	err := json.NewDecoder(r.Body).Decode(&addUser)
	if err != nil {
		log.Printf("AdminSignUpUserHandler: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	addUser.Role = "admin"
	userID, err := helper.SignUpUserHelper(addUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var user models.UserModel
	user.ID = userID
	tokenString, err := middlewares.GenerateJWT(&user)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["Name"] = "token"
	resp["Value"] = tokenString
	resp["Expires"] = utilities.FetchExpireTime().String()
	jsonResponse, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %v", err)
	}
	w.Write(jsonResponse)
}
func SignUpUserHandler(w http.ResponseWriter, r *http.Request) {
	var addUser models.AddUserModel
	err := json.NewDecoder(r.Body).Decode(&addUser)
	if err != nil {
		log.Printf("SignUpUserHandler: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	addUser.Role = "user"
	userID, err := helper.SignUpUserHelper(addUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var user models.UserModel
	user.ID = userID
	tokenString, err := middlewares.GenerateJWT(&user)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["Name"] = "token"
	resp["Value"] = tokenString
	resp["Expires"] = utilities.FetchExpireTime().String()
	jsonResponse, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %v", err)
	}
	w.Write(jsonResponse)
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
	if err != nil && err == sql.ErrNoRows {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
	resp["Expires"] = utilities.FetchExpireTime().String()
	jsonResponse, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %v", err)
	}
	w.Write(jsonResponse)
}
func AddAdminRoleHandler(w http.ResponseWriter, r *http.Request) {
	role := "admin"
	var user models.UserModel
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("AddAdminRoleHandler: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = helper.AddRoleHelper(user.ID, role)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func FetchAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	user, err := helper.FetchAllUserHelper()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonResponse, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %v", err)
	}
	w.Write(jsonResponse)
}

func FetchUserHandler(w http.ResponseWriter, r *http.Request) {
	var signedUser *models.UserModel
	signedUser = middlewares.UserFromContext(r.Context())
	var userDetails models.UserDetail
	egp := new(errgroup.Group)
	addresses := make([]models.AddressModel, 0)
	cards := make([]models.CardModel, 0)
	var err error
	egp.Go(func() error {
		addresses, err = helper.FetchAddressesHelper(signedUser.ID)
		return err
	})
	egp.Go(func() error {
		cards, err = helper.FetchCardsHelper(signedUser.ID)
		return err
	})
	userDetails.ID = signedUser.ID
	userDetails.Name = signedUser.Name
	userDetails.Email = signedUser.Email
	userDetails.Role = signedUser.Role
	txErr := egp.Wait() //  waiting for secondary information
	if txErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userDetails.Address = addresses
	userDetails.Card = cards
	jsonResponse, err := json.Marshal(userDetails)
	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %v", err)
	}
	w.Write(jsonResponse)
}
