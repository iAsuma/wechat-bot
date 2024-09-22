package model

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
