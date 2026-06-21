package model

import "time"

type Conversation struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	TaskID     *uint     `json:"task_id"`
	CustomerID uint      `json:"customer_id"`
	AgentID    uint      `json:"agent_id"`
	Type       string    `gorm:"default:'human'" json:"type"` // human or ai
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type Message struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	ConversationID uint      `json:"conversation_id"`
	SenderID       uint      `json:"sender_id"`
	MessageText    string    `json:"message_text"`
	FileURL        string    `json:"file_url,omitempty"`
	FileType       string    `json:"file_type,omitempty"`
	IsRead         bool      `gorm:"default:false" json:"is_read"`
	CreatedAt      time.Time `json:"created_at"`
}
