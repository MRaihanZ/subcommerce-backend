package entity

type UserSubscription struct {
	Id              string `db:"id" json:"id"`
	OPrettyId       string `db:"order_pretty_id" json:"order_pretty_id"`
	PayName         string `db:"payment_name" json:"payment_name"`
	OSName          string `db:"order_status_name" json:"order_status_name"`
	SName           string `db:"seller_name" json:"seller_name"`
	SImg            string `db:"seller_img" json:"seller_img"`
	ProdId          int    `db:"product_id" json:"product_id"`
	ProdName        string `db:"product_name" json:"product_name"`
	ProdVarId       int    `db:"product_variant_id" json:"product_variant_id"`
	ProdVarName     string `db:"product_variant_name" json:"product_variant_name"`
	PImg            string `db:"product_img" json:"product_img"`
	Quantity        int    `db:"quantity" json:"quantity"`
	Note            string `db:"note" json:"note"`
	ProdVarPrice    int    `db:"price" json:"price"`
	ProdVarDiscount int    `db:"discount" json:"discount"`
	Interval        int    `db:"interval" json:"interval"`
	IName           string `db:"interval_name" json:"interval_name"`
	LastSentAt      string `db:"last_sent_at" json:"last_sent_at"`
	NextSend        string `db:"next_send" json:"next_send"`
	NextWarning     string `db:"next_warning_send" json:"next_warning_send"`
	NextRemove      string `db:"next_remove" json:"next_remove"`
	IsOver          bool   `db:"is_over" json:"is_over"`
}

type UserSubscriptionBySeller struct {
	Id              string `db:"id" json:"id"`
	OPrettyId       string `db:"order_pretty_id" json:"order_pretty_id"`
	PayName         string `db:"payment_name" json:"payment_name"`
	OSName          string `db:"order_status_name" json:"order_status_name"`
	UName           string `db:"user_name" json:"user_name"`
	UImg            string `db:"user_img" json:"user_img"`
	ProdId          int    `db:"product_id" json:"product_id"`
	ProdName        string `db:"product_name" json:"product_name"`
	ProdVarId       int    `db:"product_variant_id" json:"product_variant_id"`
	ProdVarName     string `db:"product_variant_name" json:"product_variant_name"`
	PImg            string `db:"product_img" json:"product_img"`
	Quantity        int    `db:"quantity" json:"quantity"`
	Note            string `db:"note" json:"note"`
	ProdVarPrice    int    `db:"price" json:"price"`
	ProdVarDiscount int    `db:"discount" json:"discount"`
	Interval        int    `db:"interval" json:"interval"`
	IName           string `db:"interval_name" json:"interval_name"`
	LastSentAt      string `db:"last_sent_at" json:"last_sent_at"`
	NextSend        string `db:"next_send" json:"next_send"`
	NextWarning     string `db:"next_warning_send" json:"next_warning_send"`
	NextRemove      string `db:"next_remove" json:"next_remove"`
	IsOver          bool   `db:"is_over" json:"is_over"`
}

type SubscriptionIdPaymentLink struct {
	Id      string `db:"id"`
	PayLInk string `db:"payment_link"`
}
