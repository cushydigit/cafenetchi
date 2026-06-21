package model

import "time"

type Role string

const (
	RoleCustomer Role = "customer"
	RoleAgent    Role = "agent"
	RoleAdmin    Role = "admin"
)

type User struct {
	ID         uint   `gorm:"primarykey" json:"id"`
	Phone      string `gorm:"unique;not null" json:"phone"`
	FullName   string `json:"full_name"`
	AvatarURL  string `json:"avatar_url"`
	IsVerified bool   `gorm:"default:false" json:"is_verified"`
	Status     string `gorm:"default:'active'" json:"status"`

	Roles []Role `gorm:"-" json:"roles"` // Will be handled via UserProfile

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserProfile struct {
	ID       uint `gorm:"primarykey" json:"id"`
	UserID   uint `gorm:"uniqueIndex" json:"user_id"`
	Role     Role `gorm:"not null" json:"role"`
	IsActive bool `gorm:"default:true" json:"is_active"`

	// Agent fields
	AgentLevel string  `json:"agent_level"`
	HourlyRate float64 `json:"hourly_rate"`
	Rating     float64 `gorm:"default:0" json:"rating"`
	TotalTasks int     `gorm:"default:0" json:"total_tasks"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
