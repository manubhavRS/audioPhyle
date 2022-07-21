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
  		    SET archived_at=now()
  		    WHERE id=$1 and user_id=$2
  		    RETURNING id`

	err := database.Aph.Get(SQL, card.CardID, card.UserID)
	if err != nil {
		log.Printf("RemoveCardHelper Error: %v", err)
		return err
	}
	return err
}
func FetchCardsHelper(userID string) ([]models.CardModel, error) {
	//language=SQL
	SQL := `SELECT id,card_number,expire_date 
		  FROM cards
		  WHERE user_id=$1 and
		  archived_at IS NULL`
	cards := make([]models.CardModel, 0)
	err := database.Aph.Select(&cards, SQL, userID)
	if err != nil {
		log.Printf("FetchCardsHelper Error: %v", err)
		return cards, err
	}
	return cards, nil
}
