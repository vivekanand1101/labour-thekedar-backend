package otp

import "context"

// Provider defines the interface for OTP providers
type Provider interface {
	// SendOTP sends an OTP to the given phone number
	SendOTP(ctx context.Context, phone string) (string, error)
	// VerifyOTP verifies the OTP for the given phone number
	VerifyOTP(ctx context.Context, phone string, code string) (bool, error)
}
