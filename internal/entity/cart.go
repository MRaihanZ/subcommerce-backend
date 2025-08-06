package entity

type CartProduct struct {
	SName      string `db:"s_name" json:"s_name"`
	PId        int    `db:"p_id" json:"p_id"`
	PName      string `db:"p_name" json:"p_name"`
	PImage     string `db:"img" json:"img"`
	PActive    bool   `db:"active" json:"active"`
	PvId       int    `db:"pv_id" json:"pv_id"`
	PvName     string `db:"pv_name" json:"pv_name"`
	PvInterval int    `db:"interval" json:"interval"`
	PvPrice    int    `db:"price" json:"price"`
	PvMinOrder int    `db:"min_order" json:"min_order"`
	PvStock    int    `db:"stock" json:"stock"`
	IName      string `db:"i_name" json:"i_name"`
	Quantity   int    `db:"quantity" json:"quantity"`
}

type AddProductCart struct {
	PId      int `json:"p_id"`
	PvId     int `json:"pv_id"`
	Quantity int `json:"quantity"`
}

type UpdateProductCart struct {
	PId      int `json:"p_id"`
	PvId     int `json:"pv_id"`
	Quantity int `json:"quantity"`
}

type DeleteProductCart struct {
	PId  int `db:"product_id" json:"p_id"`
	PvId int `db:"product_variant_id" json:"pv_id"`
}
