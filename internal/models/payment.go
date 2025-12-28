package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// PaymentType represents the type of payment
type PaymentType string

const (
	PaymentTypeAdvance   PaymentType = "advance"
	PaymentTypeDailyWage PaymentType = "daily_wage"
	PaymentTypeBonus     PaymentType = "bonus"
)

// IsValid checks if the payment type is valid
func (pt PaymentType) IsValid() bool {
	switch pt {
	case PaymentTypeAdvance, PaymentTypeDailyWage, PaymentTypeBonus:
		return true
	}
	return false
}

// Payment represents a payment made to a labour
type Payment struct {
	ID          uuid.UUID       `json:"id" db:"id"`
	ProjectID   uuid.UUID       `json:"project_id" db:"project_id"`
	LabourID    uuid.UUID       `json:"labour_id" db:"labour_id"`
	Amount      decimal.Decimal `json:"amount" db:"amount"`
	PaymentDate time.Time       `json:"payment_date" db:"payment_date"`
	PaymentType PaymentType     `json:"payment_type" db:"payment_type"`
	Notes       string          `json:"notes,omitempty" db:"notes"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
}

// PaymentWithLabour represents a payment with labour details
type PaymentWithLabour struct {
	Payment
	LabourName string `json:"labour_name"`
}

// CreatePaymentRequest represents the request to create a payment
type CreatePaymentRequest struct {
	LabourID    uuid.UUID       `json:"labour_id" binding:"required"`
	Amount      decimal.Decimal `json:"amount" binding:"required"`
	PaymentDate string          `json:"payment_date" binding:"required"` // Format: YYYY-MM-DD
	PaymentType PaymentType     `json:"payment_type" binding:"required"`
	Notes       string          `json:"notes" binding:"max=500"`
}

// BalanceResponse represents the balance for a labour
type BalanceResponse struct {
	LabourID    uuid.UUID       `json:"labour_id"`
	LabourName  string          `json:"labour_name"`
	TotalEarned decimal.Decimal `json:"total_earned"`
	TotalPaid   decimal.Decimal `json:"total_paid"`
	Balance     decimal.Decimal `json:"balance"` // TotalEarned - TotalPaid (positive = due, negative = overpaid)
}

// Validate validates the payment data
func (p *Payment) Validate() error {
	if p.ProjectID == uuid.Nil {
		return ErrInvalidProject
	}
	if p.LabourID == uuid.Nil {
		return ErrInvalidLabour
	}
	if p.Amount.IsNegative() || p.Amount.IsZero() {
		return ErrInvalidAmount
	}
	if p.PaymentDate.IsZero() {
		return ErrInvalidDate
	}
	if !p.PaymentType.IsValid() {
		return ErrInvalidPaymentType
	}
	return nil
}
