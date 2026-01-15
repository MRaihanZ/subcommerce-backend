package entity

type ConversationRoom struct {
	Id   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Img  string `db:"img" json:"img"`
}

type Conversation struct {
	Id                 string `db:"id" json:"id"`
	UserId             string `db:"user_id" json:"user_id"`
	SellerId           string `db:"seller_id" json:"seller_id"`
	LastMessageAt      string `db:"last_message_at" json:"last_message_at"`
	LastMessageContent string `db:"last_message_content" json:"last_message_content"`
}
