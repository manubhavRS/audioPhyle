package handlers

import (
	"audioPhile/database/helper"
	"audioPhile/middlewares"
	"audioPhile/models"
	"audioPhile/utilities"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

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

//func FetchUserHandler(w http.ResponseWriter, r *http.Request) {
//	var user models.UserModel
//	err := json.NewDecoder(r.Body).Decode(&user)
//	if err != nil {
//		log.Printf("FetchUserHandler: %v", err)
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	err = helper.FetchUserHelper(user.ID)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//}
