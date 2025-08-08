package entity

type SellerSummarize struct {
	Id            string `db:"id" json:"id"`
	Img           string `db:"img" json:"img"`
	Sold          string `db:"sold_products" json:"sold_products"`
	AverageRating string `db:"average_rating" json:"average_rating"`
}
