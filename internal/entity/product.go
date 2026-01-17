package entity

import (
	"time"
)

type QueryProduct struct {
	PId           int       `db:"p_id" json:"p_id"`
	PName         string    `db:"p_name" json:"p_name"`
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
	PId             int              `db:"p_id" json:"p_id"`
	PName           string           `db:"p_name" json:"p_name"`
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
	PId           int     `db:"p_id" json:"p_id"`
	PName         string  `db:"p_name" json:"p_name"`
	PvId          string  `db:"pv_id" json:"pv_id"`
	PvName        string  `db:"pv_name" json:"pv_name"`
	Image         string  `db:"img" json:"img"`
	Price         int     `db:"price" json:"price"`
	AverageRating float32 `db:"average_rating" json:"average_rating"`
	Sold          int     `db:"sold" json:"sold"`
	Discount      int     `db:"discount" json:"discount"`
	SellerName    string  `db:"seller_name" json:"seller_name"`
}

type QueryProductAdd struct {
	PId           int     `db:"p_id" json:"p_id"`
	PName         string  `db:"p_name" json:"p_name"`
	PiImage       string  `db:"img" json:"pi_image"`
	Description   string  `db:"description" json:"description"`
	Sold          int     `db:"sold" json:"sold"`
	AverageRating float32 `db:"average_rating" json:"average_rating"`
	Active        bool    `db:"active" json:"active"`
	PvId          int     `db:"pv_id" json:"pv_id"`
	PvIsDefault   bool    `db:"is_default" json:"is_default"`
	PvName        string  `db:"pv_name" json:"pv_name"`
	IId           int     `db:"i_id" json:"i_id"`
	IName         string  `db:"i_name" json:"i_name"`
	PvInterval    int     `db:"interval" json:"interval"`
	PvStock       int     `db:"stock" json:"stock"`
	PvSold        int     `db:"pv_sold" json:"pv_sold"`
	PvPrice       int     `db:"price" json:"price"`
	PvDiscount    int     `db:"discount" json:"discount"`
	PvMinOrder    int     `db:"min_order" json:"min_order"`
}

type JsonProductAdd struct {
	PId             int                 `db:"p_id" json:"p_id"`
	PName           string              `db:"p_name" json:"p_name"`
	PiImages        []ProductImage      `json:"images"`
	Description     string              `db:"description" json:"description"`
	Sold            int                 `db:"sold" json:"sold"`
	AverageRating   float32             `db:"average_rating" json:"average_rating"`
	Active          bool                `db:"active" json:"active"`
	ProductVariants []ProductVariantAdd `json:"product_variants"`
}

type ProductVariantAdd struct {
	PvId        int    `db:"pv_id" json:"pv_id"`
	PvIsDefault bool   `db:"is_default" json:"is_default"`
	PvName      string `db:"pv_name" json:"pv_name"`
	IId         int    `db:"i_id" json:"i_id"`
	IName       string `db:"i_name" json:"i_name"`
	PvInterval  int    `db:"interval" json:"interval"`
	PvStock     int    `db:"stock" json:"stock"`
	PvSold      int    `db:"pv_sold" json:"pv_sold"`
	PvPrice     int    `db:"price" json:"price"`
	PvDiscount  int    `db:"discount" json:"discount"`
	PvMinOrder  int    `db:"min_order" json:"min_order"`
}

type ReqProductVariants struct {
	IId       int    `json:"i_id"`
	IsDefault bool   `json:"is_default"`
	Name      string `json:"name"`
	Interval  int    `json:"interval"`
	Stock     int    `json:"stock"`
	Price     int    `json:"price"`
	Discount  int    `json:"discount"`
	MinOrder  int    `json:"min_order"`
}

type ReqProductAdd struct {
	Name        string               `json:"name" binding:"required"`
	Description string               `json:"description"`
	Active      bool                 `json:"active"`
	PVariants   []ReqProductVariants `json:"p_variants"`
}

type ReqProductVariantUpdate struct {
	IId      int    `json:"i_id"`
	Name     string `json:"name"`
	Interval int    `json:"interval"`
	Stock    int    `json:"stock"`
	Price    int    `json:"price"`
	Discount int    `json:"discount"`
	MinOrder int    `json:"min_order"`
}

type ReqProductUpdate struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Stock       int     `json:"stock"`
	Price       float64 `json:"price"`
	Discount    float64 `json:"discount"`
	MinPurchase int     `json:"min_purchase"`
	HasVariant  bool    `json:"has_variant"`
	IsActive    bool    `json:"is_active"`
	Interval    int     `json:"interval"`
	IID         int     `json:"i_id"`
}
