package models

import "errors"

// Validation errors
var (
	ErrInvalidPhone       = errors.New("invalid phone number")
	ErrInvalidName        = errors.New("invalid name")
	ErrInvalidAmount      = errors.New("invalid amount")
	ErrInvalidDate        = errors.New("invalid date")
	ErrInvalidStatus      = errors.New("invalid status")
	ErrInvalidPaymentType = errors.New("invalid payment type")
	ErrInvalidProject     = errors.New("invalid project")
	ErrInvalidLabour      = errors.New("invalid labour")
	ErrInvalidOTP         = errors.New("invalid or expired OTP")
)

// Database errors
var (
	ErrNotFound      = errors.New("record not found")
	ErrAlreadyExists = errors.New("record already exists")
	ErrUnauthorized  = errors.New("unauthorized access")
	ErrForbidden     = errors.New("forbidden access")
)
