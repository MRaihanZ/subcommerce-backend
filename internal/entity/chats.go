package entity

type ConversationRoom struct {
	Id   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Img  string `db:"img" json:"img"`
}

type Conversation struct {
	Id                 string `db:"id" json:"id"`
	LastMessageAt      string `db:"last_message_at" json:"last_message_at"`
	LastMessageContent string `db:"last_message_content" json:"last_message_content"`
}
