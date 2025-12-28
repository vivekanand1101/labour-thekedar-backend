package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vivekanand/labour-thekedar-backend/internal/models"
	"github.com/vivekanand/labour-thekedar-backend/internal/repository"
)

// WorkDayService handles work day business logic
type WorkDayService struct {
	workDayRepo *repository.WorkDayRepository
	labourRepo  *repository.LabourRepository
}

// NewWorkDayService creates a new WorkDayService
func NewWorkDayService(workDayRepo *repository.WorkDayRepository, labourRepo *repository.LabourRepository) *WorkDayService {
	return &WorkDayService{
		workDayRepo: workDayRepo,
		labourRepo:  labourRepo,
	}
}

// Create creates a new work day record
func (s *WorkDayService) Create(ctx context.Context, projectID uuid.UUID, req *models.CreateWorkDayRequest) (*models.WorkDay, error) {
	// Parse work date
	workDate, err := time.Parse("2006-01-02", req.WorkDate)
	if err != nil {
		return nil, models.ErrInvalidDate
	}

	// Verify labour is assigned to project
	isAssigned, err := s.labourRepo.IsAssignedToProject(ctx, projectID, req.LabourID)
	if err != nil {
		return nil, err
	}
	if !isAssigned {
		return nil, models.ErrInvalidLabour
	}

	workDay := &models.WorkDay{
		ProjectID: projectID,
		LabourID:  req.LabourID,
		WorkDate:  workDate,
		Status:    req.Status,
		Notes:     req.Notes,
	}

	if err := workDay.Validate(); err != nil {
		return nil, err
	}

	if err := s.workDayRepo.Create(ctx, workDay); err != nil {
		return nil, err
	}

	return workDay, nil
}

// GetByID retrieves a work day by ID
func (s *WorkDayService) GetByID(ctx context.Context, id uuid.UUID) (*models.WorkDay, error) {
	return s.workDayRepo.GetByID(ctx, id)
}

// GetByProjectID retrieves all work days for a project
func (s *WorkDayService) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]models.WorkDayWithLabour, error) {
	return s.workDayRepo.GetByProjectID(ctx, projectID)
}

// GetByProjectAndDate retrieves work days for a project on a specific date
func (s *WorkDayService) GetByProjectAndDate(ctx context.Context, projectID uuid.UUID, dateStr string) ([]models.WorkDayWithLabour, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, models.ErrInvalidDate
	}
	return s.workDayRepo.GetByProjectAndDate(ctx, projectID, date)
}

// GetByLabourID retrieves all work days for a labour
func (s *WorkDayService) GetByLabourID(ctx context.Context, labourID uuid.UUID) ([]models.WorkDay, error) {
	return s.workDayRepo.GetByLabourID(ctx, labourID)
}

// Update updates a work day record
func (s *WorkDayService) Update(ctx context.Context, id uuid.UUID, req *models.UpdateWorkDayRequest) (*models.WorkDay, error) {
	workDay, err := s.workDayRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	workDay.Status = req.Status
	workDay.Notes = req.Notes

	if err := workDay.Validate(); err != nil {
		return nil, err
	}

	if err := s.workDayRepo.Update(ctx, workDay); err != nil {
		return nil, err
	}

	return workDay, nil
}

// Delete deletes a work day record
func (s *WorkDayService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.workDayRepo.Delete(ctx, id)
}
