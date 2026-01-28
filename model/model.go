package model

import (
	"time"
)

type Chat struct {
	Id        int       `json:"id"`
	Tittle    string    `json:"tittle"`
	CreatedAt time.Time `json:"created_at"`
	Messages  []Message `json:"messages,omitempty"`
}

type Message struct {
	Id        int       `json:"id"`
	ChatID    int       `json:"chat_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
