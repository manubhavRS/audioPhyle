package helper

import (
	"audioPhile/database"
	"audioPhile/models"
	"log"
)

func AddCardHelper(card models.AddCardModel) (string, error) {
	//language=SQL
	SQL := `INSERT INTO cards(user_id,card_number,expire_date) 
			values($1,$2,$3::DATE) 
			Returning id`
	var cardID string
	err := database.Aph.Get(&cardID, SQL, card.UserID, card.CardNumber, card.Expiry)
	if err != nil {
		log.Printf("AddCardHelper Error: %v", err)
		return "", err
	}
	return cardID, nil
}
func RemoveCardHelper(card models.RemoveCardIDModel) error {
	//language=SQL
	SQL := `UPDATE cards
  		  SET archived_at=CURRENT_TIMESTAMP
  		  WHERE id=$1 and user_id=$2`
	_, err := database.Aph.Exec(SQL, card.CardID, card.UserID)
	if err != nil {
		log.Printf("RemoveCardHelper Error: %v", err)
		return err
	}
	return err
}
