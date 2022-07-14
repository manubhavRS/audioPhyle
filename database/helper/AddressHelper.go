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
  		  WHERE id=$1 and user_id=$2`
	_, err := database.Aph.Exec(SQL, address.AddressID, address.UserID)
	if err != nil {
		log.Printf("RemoveAddressHelper Error: %v", err)
		return err
	}
	return err
}
