package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/vivekanand/labour-thekedar-backend/internal/models"
	"github.com/vivekanand/labour-thekedar-backend/internal/repository"
)

// ProjectService handles project business logic
type ProjectService struct {
	projectRepo *repository.ProjectRepository
	labourRepo  *repository.LabourRepository
}

// NewProjectService creates a new ProjectService
func NewProjectService(projectRepo *repository.ProjectRepository, labourRepo *repository.LabourRepository) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
		labourRepo:  labourRepo,
	}
}

// Create creates a new project
func (s *ProjectService) Create(ctx context.Context, userID uuid.UUID, req *models.CreateProjectRequest) (*models.Project, error) {
	project := &models.Project{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := project.Validate(); err != nil {
		return nil, err
	}

	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

// GetByID retrieves a project by ID
func (s *ProjectService) GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	return s.projectRepo.GetByID(ctx, id)
}

// GetByIDWithLabours retrieves a project with its labours
func (s *ProjectService) GetByIDWithLabours(ctx context.Context, id uuid.UUID) (*models.ProjectWithLabours, error) {
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	labours, err := s.labourRepo.GetByProjectID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.ProjectWithLabours{
		Project: *project,
		Labours: labours,
	}, nil
}

// GetByUserID retrieves all projects for a user
func (s *ProjectService) GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.Project, error) {
	return s.projectRepo.GetByUserID(ctx, userID)
}

// Update updates a project
func (s *ProjectService) Update(ctx context.Context, id uuid.UUID, req *models.UpdateProjectRequest) (*models.Project, error) {
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	project.Name = req.Name
	project.Description = req.Description

	if err := project.Validate(); err != nil {
		return nil, err
	}

	if err := s.projectRepo.Update(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

// Delete deletes a project
func (s *ProjectService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.projectRepo.Delete(ctx, id)
}

// IsOwner checks if a user owns a project
func (s *ProjectService) IsOwner(ctx context.Context, projectID, userID uuid.UUID) (bool, error) {
	return s.projectRepo.IsOwner(ctx, projectID, userID)
}
