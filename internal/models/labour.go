package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Labour represents a labourer in the system
type Labour struct {
	ID        uuid.UUID       `json:"id" db:"id"`
	Name      string          `json:"name" db:"name"`
	Phone     string          `json:"phone,omitempty" db:"phone"`
	DailyWage decimal.Decimal `json:"daily_wage" db:"daily_wage"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

// ProjectLabour represents the association between a project and a labour
type ProjectLabour struct {
	ProjectID  uuid.UUID `json:"project_id" db:"project_id"`
	LabourID   uuid.UUID `json:"labour_id" db:"labour_id"`
	AssignedAt time.Time `json:"assigned_at" db:"assigned_at"`
}

// LabourWithBalance represents a labour with their payment balance
type LabourWithBalance struct {
	Labour
	TotalEarned decimal.Decimal `json:"total_earned"`
	TotalPaid   decimal.Decimal `json:"total_paid"`
	Balance     decimal.Decimal `json:"balance"` // TotalEarned - TotalPaid
}

// CreateLabourRequest represents the request to create a labour
type CreateLabourRequest struct {
	Name      string          `json:"name" binding:"required,max=255"`
	Phone     string          `json:"phone" binding:"max=20"`
	DailyWage decimal.Decimal `json:"daily_wage" binding:"required"`
}

// UpdateLabourRequest represents the request to update a labour
type UpdateLabourRequest struct {
	Name      string          `json:"name" binding:"required,max=255"`
	Phone     string          `json:"phone" binding:"max=20"`
	DailyWage decimal.Decimal `json:"daily_wage" binding:"required"`
}

// AssignLabourRequest represents the request to assign a labour to a project
type AssignLabourRequest struct {
	LabourID uuid.UUID `json:"labour_id" binding:"required"`
}

// Validate validates the labour data
func (l *Labour) Validate() error {
	if l.Name == "" {
		return ErrInvalidName
	}
	if len(l.Name) > 255 {
		return ErrInvalidName
	}
	if l.DailyWage.IsNegative() {
		return ErrInvalidAmount
	}
	return nil
}
