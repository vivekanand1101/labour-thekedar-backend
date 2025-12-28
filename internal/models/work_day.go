package models

import (
	"time"

	"github.com/google/uuid"
)

// WorkStatus represents the status of a work day
type WorkStatus string

const (
	WorkStatusFullDay WorkStatus = "full_day"
	WorkStatusHalfDay WorkStatus = "half_day"
	WorkStatusAbsent  WorkStatus = "absent"
)

// IsValid checks if the work status is valid
func (ws WorkStatus) IsValid() bool {
	switch ws {
	case WorkStatusFullDay, WorkStatusHalfDay, WorkStatusAbsent:
		return true
	}
	return false
}

// WorkDay represents a work day record for a labour
type WorkDay struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	ProjectID uuid.UUID  `json:"project_id" db:"project_id"`
	LabourID  uuid.UUID  `json:"labour_id" db:"labour_id"`
	WorkDate  time.Time  `json:"work_date" db:"work_date"`
	Status    WorkStatus `json:"status" db:"status"`
	Notes     string     `json:"notes,omitempty" db:"notes"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

// WorkDayWithLabour represents a work day with labour details
type WorkDayWithLabour struct {
	WorkDay
	LabourName string `json:"labour_name"`
}

// CreateWorkDayRequest represents the request to create a work day
type CreateWorkDayRequest struct {
	LabourID uuid.UUID  `json:"labour_id" binding:"required"`
	WorkDate string     `json:"work_date" binding:"required"` // Format: YYYY-MM-DD
	Status   WorkStatus `json:"status" binding:"required"`
	Notes    string     `json:"notes" binding:"max=500"`
}

// UpdateWorkDayRequest represents the request to update a work day
type UpdateWorkDayRequest struct {
	Status WorkStatus `json:"status" binding:"required"`
	Notes  string     `json:"notes" binding:"max=500"`
}

// Validate validates the work day data
func (wd *WorkDay) Validate() error {
	if wd.ProjectID == uuid.Nil {
		return ErrInvalidProject
	}
	if wd.LabourID == uuid.Nil {
		return ErrInvalidLabour
	}
	if wd.WorkDate.IsZero() {
		return ErrInvalidDate
	}
	if !wd.Status.IsValid() {
		return ErrInvalidStatus
	}
	return nil
}
