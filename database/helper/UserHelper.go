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
	err = database.Aph.Get(&userID, SQL, user.Name, user.Email, pass, user.Role)
	if err != nil {
		log.Printf("SignUpUserHelper Error: %v", err)
		return "", err
	}
	return userID, nil
}
func FetchUserDetailsHelper(userID string) (models.UserModel, error) {
	//language=SQL
	SQL := `SELECT id,name,email,role 
			from users 
			where id=$1 and
		    archived_at IS NULL`
	var user models.UserModel
	err := database.Aph.Get(&user, SQL, userID)
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
func AddRoleHelper(userID, role string) error {
	//language=SQL
	SQL := `UPDATE users 
			SET role=$1 
			WHERE id=$2
			archived_at IS NULL 
			RETURNING id`
	var uid string
	err := database.Aph.Get(&uid, SQL, role, userID)
	if err != nil {
		log.Printf("AddRoleHelper Error: %v", err)
		return err
	}

	return nil
}
func FetchAllUserHelper() ([]models.UserModel, error) {
	//language=SQL
	SQL := `SELECT id,name,email,role
          FROM users
          WHERE archived_at IS NULL`
	user := make([]models.UserModel, 0)
	err := database.Aph.Select(&user, SQL)
	if err != nil {
		log.Printf("FetchAllUserHelper Error: %v", err)
		return user, err
	}
	return user, nil
}
func FetchUserHelper(userID string) (models.UserModel, error) {
	//language=SQL
	SQL := `SELECT id,name,email,role
          FROM users
          WHERE id=$1 and 
          archived_at IS NULL`
	var user models.UserModel
	err := database.Aph.Get(&user, SQL, userID)
	if err != nil {
		log.Printf("FetchUserHelper Error: %v", err)
		return user, err
	}
	return user, nil
}
