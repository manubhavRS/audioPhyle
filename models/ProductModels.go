package models

import (
	_ "github.com/go-playground/validator"
	"github.com/volatiletech/null"
)

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
type ProductListModel struct {
	TotalCount int            `db:"total_count" json:"totalCount"`
	Products   []ProductModel `db:"products" json:"products"`
}

type ProductModel struct {
	ID       null.String `db:"id" json:"id"`
	Name     null.String `db:"name" json:"name"`
	Category null.String `db:"category" json:"category"`
	Price    null.String `db:"price" json:"price"`
	About    null.String `db:"about" json:"about"`
	Quantity null.String `db:"quantity" json:"quantity"`
}
type ProductAssetModel struct {
	ID        string `db:"id" json:"id"`
	AssetName string `db:"name" json:"name"`
}
type ImageStructure struct {
	ImageName string
	URL       string
}
type ProductID struct {
	ID string `db:"id" json:"id"`
}
type ProductSearchModel struct {
	PageNo        int
	Category      string
	Search        string
	Limit         int
	Latest        string
	YML           string
	OrderBy       string
	CheckCategory string
}
type ProductListSearchModel struct {
	PageNo        int
	Category      []string
	Search        string
	Limit         int
	Latest        bool
	YML           bool
	OrderBy       string
	CheckCategory bool
}
type CategoryModel struct {
	Category string `db:"category" json:"category"`
}
type FetchCategoryModel struct {
	ID       null.String `db:"id" json:"id"`
	Category string      `db:"category" json:"category"`
}
type ProductModelWithCount struct {
	ID         null.String `db:"id" json:"id"`
	Name       null.String `db:"name" json:"name"`
	Category   null.String `db:"category" json:"category"`
	Price      null.String `db:"price" json:"price"`
	About      null.String `db:"about" json:"about"`
	Quantity   null.String `db:"quantity" json:"quantity"`
	TotalCount int         `db:"full_count" json:"full_count"`
}
