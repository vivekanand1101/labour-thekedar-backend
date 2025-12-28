package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vivekanand/labour-thekedar-backend/internal/models"
)

// WorkDayRepository handles work day database operations
type WorkDayRepository struct {
	db *pgxpool.Pool
}

// NewWorkDayRepository creates a new WorkDayRepository
func NewWorkDayRepository(db *pgxpool.Pool) *WorkDayRepository {
	return &WorkDayRepository{db: db}
}

// Create creates a new work day record
func (r *WorkDayRepository) Create(ctx context.Context, workDay *models.WorkDay) error {
	query := `
		INSERT INTO work_days (project_id, labour_id, work_date, status, notes)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, query, workDay.ProjectID, workDay.LabourID,
		workDay.WorkDate, workDay.Status, workDay.Notes).
		Scan(&workDay.ID, &workDay.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

// GetByID retrieves a work day by ID
func (r *WorkDayRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.WorkDay, error) {
	query := `
		SELECT id, project_id, labour_id, work_date, status, notes, created_at
		FROM work_days
		WHERE id = $1
	`

	workDay := &models.WorkDay{}
	err := r.db.QueryRow(ctx, query, id).
		Scan(&workDay.ID, &workDay.ProjectID, &workDay.LabourID,
			&workDay.WorkDate, &workDay.Status, &workDay.Notes, &workDay.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	return workDay, nil
}

// GetByProjectID retrieves all work days for a project
func (r *WorkDayRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]models.WorkDayWithLabour, error) {
	query := `
		SELECT wd.id, wd.project_id, wd.labour_id, wd.work_date, wd.status, wd.notes, wd.created_at, l.name
		FROM work_days wd
		INNER JOIN labours l ON wd.labour_id = l.id
		WHERE wd.project_id = $1
		ORDER BY wd.work_date DESC, l.name ASC
	`

	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workDays []models.WorkDayWithLabour
	for rows.Next() {
		var wd models.WorkDayWithLabour
		err := rows.Scan(&wd.ID, &wd.ProjectID, &wd.LabourID,
			&wd.WorkDate, &wd.Status, &wd.Notes, &wd.CreatedAt, &wd.LabourName)
		if err != nil {
			return nil, err
		}
		workDays = append(workDays, wd)
	}

	return workDays, rows.Err()
}

// GetByProjectAndDate retrieves work days for a project on a specific date
func (r *WorkDayRepository) GetByProjectAndDate(ctx context.Context, projectID uuid.UUID, date time.Time) ([]models.WorkDayWithLabour, error) {
	query := `
		SELECT wd.id, wd.project_id, wd.labour_id, wd.work_date, wd.status, wd.notes, wd.created_at, l.name
		FROM work_days wd
		INNER JOIN labours l ON wd.labour_id = l.id
		WHERE wd.project_id = $1 AND wd.work_date = $2
		ORDER BY l.name ASC
	`

	rows, err := r.db.Query(ctx, query, projectID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workDays []models.WorkDayWithLabour
	for rows.Next() {
		var wd models.WorkDayWithLabour
		err := rows.Scan(&wd.ID, &wd.ProjectID, &wd.LabourID,
			&wd.WorkDate, &wd.Status, &wd.Notes, &wd.CreatedAt, &wd.LabourName)
		if err != nil {
			return nil, err
		}
		workDays = append(workDays, wd)
	}

	return workDays, rows.Err()
}

// GetByLabourID retrieves all work days for a labour
func (r *WorkDayRepository) GetByLabourID(ctx context.Context, labourID uuid.UUID) ([]models.WorkDay, error) {
	query := `
		SELECT id, project_id, labour_id, work_date, status, notes, created_at
		FROM work_days
		WHERE labour_id = $1
		ORDER BY work_date DESC
	`

	rows, err := r.db.Query(ctx, query, labourID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workDays []models.WorkDay
	for rows.Next() {
		var wd models.WorkDay
		err := rows.Scan(&wd.ID, &wd.ProjectID, &wd.LabourID,
			&wd.WorkDate, &wd.Status, &wd.Notes, &wd.CreatedAt)
		if err != nil {
			return nil, err
		}
		workDays = append(workDays, wd)
	}

	return workDays, rows.Err()
}

// Update updates a work day record
func (r *WorkDayRepository) Update(ctx context.Context, workDay *models.WorkDay) error {
	query := `
		UPDATE work_days
		SET status = $2, notes = $3
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query, workDay.ID, workDay.Status, workDay.Notes)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}

// Delete deletes a work day record
func (r *WorkDayRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM work_days WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}
