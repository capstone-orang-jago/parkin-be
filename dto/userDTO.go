package dto

import "time"

type UserRequestDTO struct {
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UserResponseDTO struct {
	ID         string    `json:"id"`
	Name       string    `json:"name,omitempty"`
	Telephone  string    `json:"telephone,omitempty"`
	Email      string    `json:"email,omitempty"`
	Password   string    `json:"password,omitempty"`
	Role       string    `json:"role,omitempty"`
	Saldo      int64     `json:"saldo,omitempty"`
	IsVerified bool      `json:"is_verified,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}
