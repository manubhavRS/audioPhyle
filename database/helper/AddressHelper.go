package helper

import (
	"audioPhile/database"
	"audioPhile/models"
	"log"
)

func AddAddressHelper(address models.AddAddressModel) (string, error) {
	//language=SQL
	SQL := `INSERT INTO addresses(user_id,address,phone_number,landmark) 
			values($1,$2,$3,$4) 
			Returning id`
	var addressID string
	err := database.Aph.Get(&addressID, SQL, address.UserID, address.Address, address.PhoneNumber, address.Landmark)
	if err != nil {
		log.Printf("AddAddressHelper Error: %v", err)
		return "", err
	}
	return addressID, nil
}
func RemoveAddressHelper(address models.RemoveAddressIDModel) error {
	//language=SQL
	SQL := `UPDATE addresses
  		    SET archived_at=CURRENT_TIMESTAMP
  		    WHERE id=$1 and user_id=$2
  		    RETURNING id`
	err := database.Aph.Get(SQL, address.AddressID, address.UserID)
	if err != nil {
		log.Printf("RemoveAddressHelper Error: %v", err)
		return err
	}
	return err
}
func FetchAddressesHelper(userID string) ([]models.AddressModel, error) {
	//language=SQL
	SQL := `SELECT id,address,landmark,phone_number 
		  FROM addresses
		  WHERE user_id=$1 and
		  archived_at IS NULL`
	addresses := make([]models.AddressModel, 0)
	err := database.Aph.Select(&addresses, SQL, userID)
	if err != nil {
		log.Printf("FetchAddressesHelper Error: %v", err)
		return addresses, err
	}
	return addresses, nil
}
