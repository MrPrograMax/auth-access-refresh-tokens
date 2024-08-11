package domain

import "time"

type Session struct {
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	IpAddress    string    `json:"ip_address" db:"ip_address"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
}
