package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/vivekanand/labour-thekedar-backend/internal/models"
)

// PaymentRepository handles payment database operations
type PaymentRepository struct {
	db *pgxpool.Pool
}

// NewPaymentRepository creates a new PaymentRepository
func NewPaymentRepository(db *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// Create creates a new payment record
func (r *PaymentRepository) Create(ctx context.Context, payment *models.Payment) error {
	query := `
		INSERT INTO payments (project_id, labour_id, amount, payment_date, payment_type, notes)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, query, payment.ProjectID, payment.LabourID,
		payment.Amount, payment.PaymentDate, payment.PaymentType, payment.Notes).
		Scan(&payment.ID, &payment.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

// GetByID retrieves a payment by ID
func (r *PaymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
	query := `
		SELECT id, project_id, labour_id, amount, payment_date, payment_type, notes, created_at
		FROM payments
		WHERE id = $1
	`

	payment := &models.Payment{}
	err := r.db.QueryRow(ctx, query, id).
		Scan(&payment.ID, &payment.ProjectID, &payment.LabourID,
			&payment.Amount, &payment.PaymentDate, &payment.PaymentType,
			&payment.Notes, &payment.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	return payment, nil
}

// GetByProjectID retrieves all payments for a project
func (r *PaymentRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]models.PaymentWithLabour, error) {
	query := `
		SELECT p.id, p.project_id, p.labour_id, p.amount, p.payment_date, p.payment_type, p.notes, p.created_at, l.name
		FROM payments p
		INNER JOIN labours l ON p.labour_id = l.id
		WHERE p.project_id = $1
		ORDER BY p.payment_date DESC, l.name ASC
	`

	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.PaymentWithLabour
	for rows.Next() {
		var p models.PaymentWithLabour
		err := rows.Scan(&p.ID, &p.ProjectID, &p.LabourID,
			&p.Amount, &p.PaymentDate, &p.PaymentType,
			&p.Notes, &p.CreatedAt, &p.LabourName)
		if err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}

	return payments, rows.Err()
}

// GetByLabourID retrieves all payments for a labour
func (r *PaymentRepository) GetByLabourID(ctx context.Context, labourID uuid.UUID) ([]models.Payment, error) {
	query := `
		SELECT id, project_id, labour_id, amount, payment_date, payment_type, notes, created_at
		FROM payments
		WHERE labour_id = $1
		ORDER BY payment_date DESC
	`

	rows, err := r.db.Query(ctx, query, labourID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var p models.Payment
		err := rows.Scan(&p.ID, &p.ProjectID, &p.LabourID,
			&p.Amount, &p.PaymentDate, &p.PaymentType,
			&p.Notes, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}

	return payments, rows.Err()
}

// GetBalance calculates the balance for a labour in a project
// Balance = Total Earned (from work days) - Total Paid
func (r *PaymentRepository) GetBalance(ctx context.Context, projectID, labourID uuid.UUID) (*models.BalanceResponse, error) {
	// Get labour info
	labourQuery := `SELECT name, daily_wage FROM labours WHERE id = $1`
	var labourName string
	var dailyWage decimal.Decimal
	err := r.db.QueryRow(ctx, labourQuery, labourID).Scan(&labourName, &dailyWage)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	// Calculate total earned from work days
	// full_day = 1.0 * daily_wage, half_day = 0.5 * daily_wage, absent = 0
	earnedQuery := `
		SELECT COALESCE(SUM(
			CASE status
				WHEN 'full_day' THEN 1.0
				WHEN 'half_day' THEN 0.5
				ELSE 0
			END
		), 0)
		FROM work_days
		WHERE project_id = $1 AND labour_id = $2
	`
	var workDayMultiplier decimal.Decimal
	err = r.db.QueryRow(ctx, earnedQuery, projectID, labourID).Scan(&workDayMultiplier)
	if err != nil {
		return nil, err
	}
	totalEarned := dailyWage.Mul(workDayMultiplier)

	// Calculate total paid
	paidQuery := `
		SELECT COALESCE(SUM(amount), 0)
		FROM payments
		WHERE project_id = $1 AND labour_id = $2
	`
	var totalPaid decimal.Decimal
	err = r.db.QueryRow(ctx, paidQuery, projectID, labourID).Scan(&totalPaid)
	if err != nil {
		return nil, err
	}

	return &models.BalanceResponse{
		LabourID:    labourID,
		LabourName:  labourName,
		TotalEarned: totalEarned,
		TotalPaid:   totalPaid,
		Balance:     totalEarned.Sub(totalPaid), // Positive = due, Negative = overpaid
	}, nil
}

// Delete deletes a payment record
func (r *PaymentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM payments WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}
