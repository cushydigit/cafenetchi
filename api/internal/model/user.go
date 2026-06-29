package model

import "time"

type Role string

const (
	RoleCustomer Role = "customer"
	RoleAgent    Role = "agent"
	RoleAdmin    Role = "admin"
)

type UserStatus string

const (
	UserActive  UserStatus = "active"
	UserBlocked UserStatus = "blocked"
	UserDeleted UserStatus = "deleted"
)

type User struct {
	ID        int64
	Phone     string
	FullName  string
	AvatarURL string

	IsVerified bool
	Status     UserStatus

	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserProfile struct {
	ID     int64
	UserID int64
	Role   Role

	IsActive bool

	// Agent-specific fields
	AgentLevel string
	HourlyRate float64
	Rating     float64
	TotalTasks int

	CreatedAt time.Time
	UpdatedAt time.Time
}
