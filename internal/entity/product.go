package entity

import "time"

type QueryProduct struct {
	Id            int       `db:"id" json:"id"`
	Name          string    `db:"name" json:"name"`
	PiImage       string    `db:"img" json:"pi_image"`
	Description   string    `db:"description" json:"description"`
	Sold          int       `db:"sold" json:"sold"`
	AverageRating float32   `db:"average_rating" json:"average_rating"`
	RatingCount   int       `db:"rating_count" json:"rating_count"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	ICode         string    `db:"code" json:"code"`
	IName         string    `db:"i_name" json:"i_name"`
	PvId          int       `db:"pv_id" json:"pv_id"`
	PvIsDefault   bool      `db:"is_default" json:"is_default"`
	PvName        string    `db:"pv_name" json:"pv_name"`
	PvInterval    int       `db:"interval" json:"interval"`
	PvStock       int       `db:"stock" json:"stock"`
	PvSold        int       `db:"pv_sold" json:"pv_sold"`
	PvPrice       int       `db:"price" json:"price"`
	PvDiscount    int       `db:"discount" json:"discount"`
	PvMinOrder    int       `db:"min_order" json:"min_order"`
}

type JsonProduct struct {
	Id              int              `db:"id" json:"id"`
	Name            string           `db:"name" json:"name"`
	PiImages        []ProductImage   `json:"images"`
	Description     string           `db:"description" json:"description"`
	Sold            int              `db:"sold" json:"sold"`
	AverageRating   float32          `db:"average_rating" json:"average_rating"`
	RatingCount     int              `db:"rating_count" json:"rating_count"`
	CreatedAt       time.Time        `db:"created_at" json:"created_at"`
	ProductVariants []ProductVariant `json:"product_variants"`
}

type ProductImage struct {
	Image string `db:"img" json:"img"`
}

type ProductVariant struct {
	PvId        int    `db:"pv_id" json:"pv_id"`
	PvIsDefault bool   `db:"is_default" json:"is_default"`
	PvName      string `db:"pv_name" json:"pv_name"`
	ICode       string `db:"code" json:"code"`
	IName       string `db:"i_name" json:"i_name"`
	PvInterval  int    `db:"interval" json:"interval"`
	PvStock     int    `db:"stock" json:"stock"`
	PvSold      int    `db:"pv_sold" json:"pv_sold"`
	PvPrice     int    `db:"price" json:"price"`
	PvDiscount  int    `db:"discount" json:"discount"`
	PvMinOrder  int    `db:"min_order" json:"min_order"`
}

type ProductSummarize struct {
	Id            int     `db:"id" json:"id"`
	Name          string  `db:"name" json:"name"`
	PvName        string  `db:"pv_name" json:"pv_name"`
	Image         string  `db:"img" json:"img"`
	Price         int     `db:"price" json:"price"`
	AverageRating float32 `db:"average_rating" json:"average_rating"`
	Sold          int     `db:"sold" json:"sold"`
	SellerName    string  `db:"seller_name" json:"seller_name"`
}
