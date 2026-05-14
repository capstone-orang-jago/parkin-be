package dto

import "time"

type (
	LoginRequestDTO struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	LoginResponseDTO struct {
		UserID       string `json:"user_id"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	RefreshTokenRequestDTO struct {
		RefreshToken string `json:"refresh_token"`
	}
	RefreshTokenResponseDTO struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	RegisterRequestDTO struct {
		Name      string `json:"name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Telephone string `json:"telephone"`
	}
	RegisterResponseDTO struct {
		UserID    string    `json:"user_id"`
		Telephone string    `json:"telephone"`
		OTPSentAt time.Time `json:"otp_sent_at"`
	}
)
