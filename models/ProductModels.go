package models

type AddProductModel struct {
	Name     string             `db:"name" json:"name"`
	Category string             `db:"category" json:"category"`
	Price    string             `db:"price" json:"price"`
	Feature  string             `db:"feature" json:"feature"`
	InTheBox []AddInTheBoxModel `db:"in_the_box" json:"inTheBox"`
	About    string             `db:"about" json:"about"`
	Quantity string             `db:"quantity" json:"quantity"`
}

type AddInTheBoxModel struct {
	Name     string `db:"name" json:"name"`
	Quantity string `db:"quantity" json:"quantity"`
}

type ProductModel struct {
	ID       string `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Category string `db:"category" json:"category"`
	Price    string `db:"price" json:"price"`
	About    string `db:"about" json:"about"`
	Quantity string `db:"quantity" json:"quantity"`
}
