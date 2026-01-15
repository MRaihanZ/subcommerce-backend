package entity

type ConversationRoom struct {
	Id   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Img  string `db:"img" json:"img"`
}

type Conversation struct {
	Id                 string  `db:"id" json:"id"`
	UserId             string  `db:"user_id" json:"user_id"`
	SellerId           string  `db:"seller_id" json:"seller_id"`
	Name               string  `db:"name" json:"name"`
	Img                string  `db:"img" json:"img"`
	LastMessageAt      *string `db:"last_message_at" json:"last_message_at"`
	LastMessageContent *string `db:"last_message_content" json:"last_message_content"`
}

type Message struct {
	Id      string `db:"id" json:"id"`
	IsUser  bool   `db:"is_user" json:"is_user"`
	Content string `db:"content" json:"content"`
	SentAt  string `db:"sent_at" json:"sent_at"`
}
