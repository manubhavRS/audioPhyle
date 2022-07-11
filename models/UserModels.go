package models

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
