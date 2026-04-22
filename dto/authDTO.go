package dto

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
		AccessToken string `json:"access_token"`
	}
)
