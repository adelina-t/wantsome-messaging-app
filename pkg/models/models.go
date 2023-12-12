package models

type Message struct {
	Message   string
	UserName  string
	Recipient string // Optional: If set, the message is a direct message to this recipient
}
