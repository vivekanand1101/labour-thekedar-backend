package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vivekanand/labour-thekedar-backend/internal/models"
)

// LabourRepository handles labour database operations
type LabourRepository struct {
	db *pgxpool.Pool
}

// NewLabourRepository creates a new LabourRepository
func NewLabourRepository(db *pgxpool.Pool) *LabourRepository {
	return &LabourRepository{db: db}
}

// Create creates a new labour
func (r *LabourRepository) Create(ctx context.Context, labour *models.Labour) error {
	query := `
		INSERT INTO labours (name, phone, daily_wage)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query, labour.Name, labour.Phone, labour.DailyWage).
		Scan(&labour.ID, &labour.CreatedAt, &labour.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// GetByID retrieves a labour by ID
func (r *LabourRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Labour, error) {
	query := `
		SELECT id, name, phone, daily_wage, created_at, updated_at
		FROM labours
		WHERE id = $1
	`

	labour := &models.Labour{}
	err := r.db.QueryRow(ctx, query, id).
		Scan(&labour.ID, &labour.Name, &labour.Phone, &labour.DailyWage,
			&labour.CreatedAt, &labour.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	return labour, nil
}

// GetAll retrieves all labours
func (r *LabourRepository) GetAll(ctx context.Context) ([]models.Labour, error) {
	query := `
		SELECT id, name, phone, daily_wage, created_at, updated_at
		FROM labours
		ORDER BY name ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var labours []models.Labour
	for rows.Next() {
		var l models.Labour
		err := rows.Scan(&l.ID, &l.Name, &l.Phone, &l.DailyWage,
			&l.CreatedAt, &l.UpdatedAt)
		if err != nil {
			return nil, err
		}
		labours = append(labours, l)
	}

	return labours, rows.Err()
}

// GetByProjectID retrieves all labours assigned to a project
func (r *LabourRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]models.Labour, error) {
	query := `
		SELECT l.id, l.name, l.phone, l.daily_wage, l.created_at, l.updated_at
		FROM labours l
		INNER JOIN project_labours pl ON l.id = pl.labour_id
		WHERE pl.project_id = $1
		ORDER BY l.name ASC
	`

	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var labours []models.Labour
	for rows.Next() {
		var l models.Labour
		err := rows.Scan(&l.ID, &l.Name, &l.Phone, &l.DailyWage,
			&l.CreatedAt, &l.UpdatedAt)
		if err != nil {
			return nil, err
		}
		labours = append(labours, l)
	}

	return labours, rows.Err()
}

// Update updates a labour
func (r *LabourRepository) Update(ctx context.Context, labour *models.Labour) error {
	query := `
		UPDATE labours
		SET name = $2, phone = $3, daily_wage = $4, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query, labour.ID, labour.Name, labour.Phone, labour.DailyWage).
		Scan(&labour.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ErrNotFound
		}
		return err
	}

	return nil
}

// Delete deletes a labour
func (r *LabourRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM labours WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}

// AssignToProject assigns a labour to a project
func (r *LabourRepository) AssignToProject(ctx context.Context, projectID, labourID uuid.UUID) error {
	query := `
		INSERT INTO project_labours (project_id, labour_id)
		VALUES ($1, $2)
		ON CONFLICT (project_id, labour_id) DO NOTHING
	`

	_, err := r.db.Exec(ctx, query, projectID, labourID)
	return err
}

// RemoveFromProject removes a labour from a project
func (r *LabourRepository) RemoveFromProject(ctx context.Context, projectID, labourID uuid.UUID) error {
	query := `DELETE FROM project_labours WHERE project_id = $1 AND labour_id = $2`

	result, err := r.db.Exec(ctx, query, projectID, labourID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}

// IsAssignedToProject checks if a labour is assigned to a project
func (r *LabourRepository) IsAssignedToProject(ctx context.Context, projectID, labourID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM project_labours WHERE project_id = $1 AND labour_id = $2)`

	var exists bool
	err := r.db.QueryRow(ctx, query, projectID, labourID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
