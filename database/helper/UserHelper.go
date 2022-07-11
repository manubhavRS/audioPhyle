package helper

import (
	"audioPhile/database"
	"audioPhile/models"
	"audioPhile/utilities"
	"log"
)

func SignUpUserHelper(user models.AddUserModel) (string, error) {
	//language=SQL
	SQL := `INSERT INTO users(name,email,password,role) 
			values($1,$2,$3,$4) 
			RETURNING id`
	var userID string
	pass, err := utilities.HashPassword(user.Password)
	if err != nil {
		log.Printf("SignUpUserHelper Error: %v", err)
		return "", err
	}
	err = database.Aph.Get(&userID, SQL, user.Name, pass, user.Email, user.Role)
	if err != nil {
		log.Printf("SignUpUserHelper Error: %v", err)
		return "", err
	}
	return userID, nil
}
func FetchUserDetailsHelper(userID string) (models.UserModel, error) {
	//language=SQL
	SQL := `SELECT id,name,email,role 
			from users`
	var user models.UserModel
	err := database.Aph.Get(&user, SQL, user.ID, user.Name, user.Email, user.Role)
	if err != nil {
		log.Printf("FetchUserDetailsHelper Error: %v", err)
		return user, err
	}
	return user, nil
}
func FetchUserCredentialsHelper(userEmail string) (models.UserModel, error) {
	//language=SQL
	SQL := `SELECT id,name,email,password,role 
			from users
			where email=$1
			AND archived_at IS NULL`
	var user models.UserModel
	err := database.Aph.Get(&user, SQL, userEmail)
	if err != nil {
		log.Printf("FetchUserCredentialsHelper Error: %v", err)
		return user, err
	}
	return user, nil
}
