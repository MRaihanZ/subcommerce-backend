package entity

import "time"

type Seller struct {
	Id                         string     `db:"id" json:"id"`
	Name                       string     `db:"name" json:"name"`
	Img                        string     `db:"img" json:"img"`
	Address                    string     `db:"address" json:"address"`
	SoldProducts               int        `db:"total_sold_products" json:"total_sold_products"`
	AverageRating              float32    `db:"average_rating" json:"average_rating"`
	RatingTotal                int        `db:"rating_total" json:"rating_total"`
	RatingCount                int        `db:"rating_count" json:"rating_count"`
	CurrentMonthSales          int        `db:"current_month_sales" json:"current_month_sales"`
	PreviousMonthSales         *int       `db:"previous_month_sales" json:"previous_month_sales"`
	CurrentMonthCancellations  int        `db:"current_month_cancellations" json:"current_month_cancellations"`
	PreviousMonthCancellations *int       `db:"previous_month_cancellations" json:"previous_month_cancellations"`
	CurrentMonthRevenue        int64      `db:"current_month_revenue" json:"current_month_revenue"`
	PreviousMonthRevenue       *int64     `db:"previous_month_revenue" json:"previous_month_revenue"`
	CurrentMonth               time.Time  `db:"current_month" json:"current_month"`
	PreviousMonth              *time.Time `db:"previous_month" json:"previous_month"`
	CreatedAt                  time.Time  `db:"created_at" json:"created_at"`
}

type SellerSummarize struct {
	Id            string  `db:"id" json:"id"`
	Name          string  `db:"name" json:"name"`
	Img           string  `db:"img" json:"img"`
	Address       string  `db:"address" json:"address"`
	Sold          int     `db:"sold_products" json:"sold_products"`
	AverageRating float32 `db:"average_rating" json:"average_rating"`
}
