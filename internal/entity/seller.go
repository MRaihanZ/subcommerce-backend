package entity

import "time"

type SellerSummarize struct {
	Id            string  `db:"id" json:"id"`
	Name          string  `db:"name" json:"name"`
	Img           string  `db:"img" json:"img"`
	Sold          int     `db:"sold_products" json:"sold_products"`
	AverageRating float32 `db:"average_rating" json:"average_rating"`
}

type Seller struct {
	Id            string    `db:"id" json:"id"`
	Name          string    `db:"name" json:"name"`
	Img           string    `db:"img" json:"img"`
	Address       string    `db:"address" json:"address"`
	SoldProducts  int       `db:"sold_products" json:"sold_products"`
	AverageRating float32   `db:"average_rating" json:"average_rating"`
	RatingTotal   int       `db:"rating_total" json:"rating_total"`
	RatingCount   int       `db:"rating_count" json:"rating_count"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}
