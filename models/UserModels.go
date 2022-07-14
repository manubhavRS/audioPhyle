package models

import "time"

type SignInUserModel struct {
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}
type UserModel struct {
	ID       string `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
	Role     string `db:"role" json:"role"`
}
type AddUserModel struct {
	Name     string `db:"name" json:"name"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
	Role     string `db:"role" json:"role"`
}
type AddCardModel struct {
	UserID     string    `db:"user_id" json:"userID"`
	CardNumber string    `db:"card_number" json:"cardNumber"`
	ExpireDate string    `db:"expire_date" json:"expireDate"`
	Expiry     time.Time `db:"expire_date", json:"-"`
}
type AddAddressModel struct {
	UserID      string `db:"user_id" json:"userID"`
	Address     string `db:"address" json:"address"`
	Landmark    string `db:"landmark" json:"landmark"`
	PhoneNumber string `db:"phone_number" json:"phoneNumber"`
}
type RemoveAddressIDModel struct {
	UserID    string `db:"user_id" json:"userID"`
	AddressID string `db:"id" json:"addressID"`
}
type RemoveCardIDModel struct {
	UserID string `db:"user_id" json:"userID"`
	CardID string `db:"id" json:"cardID"`
}
