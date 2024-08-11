package domain

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID           uuid.UUID `json:"-"`
	Email        string    `json:"email" db:"email"`
	Password     string    `json:"password" db:"password"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	Ip           string    `json:"ip" db:"ip"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
}
