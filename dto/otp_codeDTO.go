package dto

type OtpCodeRequestDTO struct {
	Telephone string `json:"telephone"`
}

type VerifyOtpCodeRequestDTO struct {
	Telephone string `json:"telephone"`
	OTPCode   string `json:"otp_code"`
}
