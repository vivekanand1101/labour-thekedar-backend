package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/vivekanand/labour-thekedar-backend/internal/models"
	"github.com/vivekanand/labour-thekedar-backend/internal/repository"
)

// LabourService handles labour business logic
type LabourService struct {
	labourRepo *repository.LabourRepository
}

// NewLabourService creates a new LabourService
func NewLabourService(labourRepo *repository.LabourRepository) *LabourService {
	return &LabourService{
		labourRepo: labourRepo,
	}
}

// Create creates a new labour
func (s *LabourService) Create(ctx context.Context, req *models.CreateLabourRequest) (*models.Labour, error) {
	labour := &models.Labour{
		Name:      req.Name,
		Phone:     req.Phone,
		DailyWage: req.DailyWage,
	}

	if err := labour.Validate(); err != nil {
		return nil, err
	}

	if err := s.labourRepo.Create(ctx, labour); err != nil {
		return nil, err
	}

	return labour, nil
}

// GetByID retrieves a labour by ID
func (s *LabourService) GetByID(ctx context.Context, id uuid.UUID) (*models.Labour, error) {
	return s.labourRepo.GetByID(ctx, id)
}

// GetAll retrieves all labours
func (s *LabourService) GetAll(ctx context.Context) ([]models.Labour, error) {
	return s.labourRepo.GetAll(ctx)
}

// GetByProjectID retrieves all labours for a project
func (s *LabourService) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]models.Labour, error) {
	return s.labourRepo.GetByProjectID(ctx, projectID)
}

// Update updates a labour
func (s *LabourService) Update(ctx context.Context, id uuid.UUID, req *models.UpdateLabourRequest) (*models.Labour, error) {
	labour, err := s.labourRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	labour.Name = req.Name
	labour.Phone = req.Phone
	labour.DailyWage = req.DailyWage

	if err := labour.Validate(); err != nil {
		return nil, err
	}

	if err := s.labourRepo.Update(ctx, labour); err != nil {
		return nil, err
	}

	return labour, nil
}

// Delete deletes a labour
func (s *LabourService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.labourRepo.Delete(ctx, id)
}

// AssignToProject assigns a labour to a project
func (s *LabourService) AssignToProject(ctx context.Context, projectID, labourID uuid.UUID) error {
	// Verify labour exists
	_, err := s.labourRepo.GetByID(ctx, labourID)
	if err != nil {
		return err
	}

	return s.labourRepo.AssignToProject(ctx, projectID, labourID)
}

// RemoveFromProject removes a labour from a project
func (s *LabourService) RemoveFromProject(ctx context.Context, projectID, labourID uuid.UUID) error {
	return s.labourRepo.RemoveFromProject(ctx, projectID, labourID)
}

// IsAssignedToProject checks if a labour is assigned to a project
func (s *LabourService) IsAssignedToProject(ctx context.Context, projectID, labourID uuid.UUID) (bool, error) {
	return s.labourRepo.IsAssignedToProject(ctx, projectID, labourID)
}
