package entity

type OrderRequest struct {
	PayId      int    `db:"payment_id" json:"pay_id"`
	PId        int    `db:"product_id" json:"p_id"`
	PVId       int    `db:"product_variant_id" json:"pv_id"`
	Note       string `db:"note" json:"note"`
	Quantity   int    `db:"quantity" json:"quantity"`
	UnitPrice  int    `db:"unit_price" json:"unit_price"`
	TotalPrice int    `db:"total_price" json:"total_price"`
}

type OrderGetResponse struct {
	OId        string `db:"order_id" json:"o_id"`
	OPrettyId  string `db:"order_pretty_id" json:"o_order_pretty_id"`
	PayName    string `db:"pay_name" json:"pay_name"`
	OSName     string `db:"os_name" json:"os_name"`
	Rating     bool   `db:"rating" json:"rating"`
	PId        int    `db:"product_id" json:"p_id"`
	PVId       int    `db:"product_variant_id" json:"pv_id"`
	SName      string `db:"s_name" json:"s_name"`
	SImg       string `db:"s_img" json:"s_img"`
	PName      string `db:"p_name" json:"p_name"`
	PImg       string `db:"p_img" json:"p_img"`
	PActive    string `db:"active" json:"active"`
	PVName     string `db:"pv_name" json:"pv_name"`
	PvInterval int    `db:"interval" json:"interval"`
	IName      string `db:"i_name" json:"i_name"`
	Quantity   int    `db:"quantity" json:"quantity"`
	TotalPrice int    `db:"total_price" json:"total_price"`
}

type OrderCheckout struct {
	PId        int `db:"product_id" json:"p_id"`
	PVId       int `db:"product_variant_id" json:"pv_id"`
	Quantity   int `db:"quantity" json:"quantity"`
	UnitPrice  int `db:"unit_price" json:"unit_price"`
	TotalPrice int `db:"total_price" json:"total_price"`
}

type GetCheckoutOrderResponse struct {
	PId        int    `db:"product_id" json:"p_id"`
	PVId       int    `db:"product_variant_id" json:"pv_id"`
	SName      string `db:"s_name" json:"s_name"`
	SImg       string `db:"s_img" json:"s_img"`
	PName      string `db:"p_name" json:"p_name"`
	PImg       string `db:"p_img" json:"p_img"`
	PVName     string `db:"pv_name" json:"pv_name"`
	PvInterval int    `db:"interval" json:"interval"`
	IName      string `db:"i_name" json:"i_name"`
	Quantity   int    `db:"quantity" json:"quantity"`
	TotalPrice int    `db:"total_price" json:"total_price"`
}

type GetOrderPaymentResponse struct {
	PId                int    `db:"p_id" json:"p_id"`
	PName              string `db:"p_name" json:"p_name"`
	PCategoryPaymentId int    `db:"category_payment_id" json:"category_payment_id"`
	CPName             string `db:"cp_name" json:"cp_name"`
}
