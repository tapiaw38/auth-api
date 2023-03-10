package models

import (
	"time"
)

// User is the model for the user table
type User struct {
	Id            string    `json:"id"`
	FirstName     string    `json:"first_name,omitempty"`
	LastName      string    `json:"last_name,omitempty"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	Password      string    `json:"password"`
	PhoneNumber   string    `json:"phone_number,omitempty"`
	Picture       string    `json:"picture,omitempty"`
	Address       string    `json:"address,omitempty"`
	IsActive      bool      `json:"is_active,omitempty"`
	VerifiedEmail bool      `json:"verified_email,omitempty"`
	Token         string    `json:"token,omitempty"`
	TokenExpiry   time.Time `json:"token_expiry,omitempty"`
	Roles         []Role    `json:"roles,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}

// UserResponse is the model for the user table
type UserResponse struct {
	Id            string    `json:"id"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Picture       string    `json:"picture"`
	PhoneNumber   string    `json:"phone_number"`
	Address       string    `json:"address"`
	IsActive      bool      `json:"is_active,omitempty"`
	VerifiedEmail bool      `json:"verified_email,omitempty"`
	Token         string    `json:"token,omitempty"`
	TokenExpiry   time.Time `json:"token_expiry,omitempty"`
	Roles         []Role    `json:"roles"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
