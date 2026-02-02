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
	OId         string `db:"order_id" json:"o_id"`
	OUqId       string `db:"order_uq_id" json:"order_uq_id"`
	OPrettyId   string `db:"order_pretty_id" json:"o_order_pretty_id"`
	PayName     string `db:"pay_name" json:"pay_name"`
	OSName      string `db:"os_name" json:"os_name"`
	Rating      bool   `db:"rating" json:"rating"`
	PId         int    `db:"product_id" json:"p_id"`
	PVId        int    `db:"product_variant_id" json:"pv_id"`
	SName       string `db:"s_name" json:"s_name"`
	SImg        string `db:"s_img" json:"s_img"`
	PName       string `db:"p_name" json:"p_name"`
	PImg        string `db:"p_img" json:"p_img"`
	PActive     string `db:"active" json:"active"`
	PVName      string `db:"pv_name" json:"pv_name"`
	PvInterval  int    `db:"interval" json:"interval"`
	IName       string `db:"i_name" json:"i_name"`
	Note        string `db:"note" json:"note"`
	PaymentLink string `db:"payment_link" json:"payment_link"`
	Quantity    int    `db:"quantity" json:"quantity"`
	TotalPrice  int    `db:"total_price" json:"total_price"`
}

type OrderGetResponseBySeller struct {
	OrderID          int    `db:"order_id" json:"order_id"`
	OrderPrettyID    string `db:"order_pretty_id" json:"order_pretty_id"`
	UName            string `db:"u_name" json:"u_name"`
	UImg             string `db:"u_img" json:"u_img"`
	ProductID        int    `db:"product_id" json:"product_id"`
	ProductVariantID int    `db:"product_variant_id" json:"product_variant_id"`
	PName            string `db:"p_name" json:"p_name"`
	PvName           string `db:"pv_name" json:"pv_name"`
	PImg             string `db:"p_img" json:"p_img"`
	Quantity         int    `db:"quantity" json:"quantity"`
	Interval         int    `db:"interval" json:"interval"`
	IName            string `db:"i_name" json:"i_name"`
	PayName          string `db:"pay_name" json:"pay_name"`
	OsName           string `db:"os_name" json:"os_name"`
	TotalPrice       int    `db:"total_price" json:"total_price"`
	CreatedAt        string `db:"created_at" json:"created_at"`
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

type GetCheckoutOrderPaymentResponse struct {
	PId        int    `db:"product_id" json:"p_id"`
	PVId       int    `db:"product_variant_id" json:"pv_id"`
	SName      string `db:"s_name" json:"s_name"`
	SImg       string `db:"s_img" json:"s_img"`
	PName      string `db:"p_name" json:"p_name"`
	PImg       string `db:"p_img" json:"p_img"`
	PVName     string `db:"pv_name" json:"pv_name"`
	Discount   int    `db:"discount" json:"discount"`
	PvInterval int    `db:"interval" json:"interval"`
	IName      string `db:"i_name" json:"i_name"`
	Quantity   int    `db:"quantity" json:"quantity"`
	UnitPrice  int    `db:"unit_price" json:"unit_price"`
	TotalPrice int    `db:"total_price" json:"total_price"`
}

type OrderSubscriptionResponse struct {
	PId          int          `db:"product_id" json:"p_id"`
	PName        string       `db:"product_name" json:"p_name"`
	Quantity     int          `db:"quantity" json:"quantity"`
	UnitPrice    int          `db:"unit_price" json:"unit_price"`
	TotalPrice   int          `db:"total_price" json:"total_price"`
	OrderRequest OrderRequest `json:"order_request"`
}

type GetOrderPaymentResponse struct {
	PId                int    `db:"p_id" json:"p_id"`
	PName              string `db:"p_name" json:"p_name"`
	PCategoryPaymentId int    `db:"category_payment_id" json:"category_payment_id"`
	CPName             string `db:"cp_name" json:"cp_name"`
}

type GetUserProductProductVariant struct {
	UserId           string `db:"user_id"`
	ProductId        int    `db:"product_id"`
	ProductVariantId int    `db:"product_variant_id"`
}

type GetQuantityProductProductVariant struct {
	Quantity         string `db:"quantity"`
	ProductId        int    `db:"product_id"`
	ProductVariantId int    `db:"product_variant_id"`
}

type IntervalProduct struct {
	Id       int `db:"interval_id"`
	Interval int `db:"interval"`
}
