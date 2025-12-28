package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vivekanand/labour-thekedar-backend/internal/models"
	"github.com/vivekanand/labour-thekedar-backend/internal/repository"
)

// PaymentService handles payment business logic
type PaymentService struct {
	paymentRepo *repository.PaymentRepository
	labourRepo  *repository.LabourRepository
}

// NewPaymentService creates a new PaymentService
func NewPaymentService(paymentRepo *repository.PaymentRepository, labourRepo *repository.LabourRepository) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		labourRepo:  labourRepo,
	}
}

// Create creates a new payment record
func (s *PaymentService) Create(ctx context.Context, projectID uuid.UUID, req *models.CreatePaymentRequest) (*models.Payment, error) {
	// Parse payment date
	paymentDate, err := time.Parse("2006-01-02", req.PaymentDate)
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

	payment := &models.Payment{
		ProjectID:   projectID,
		LabourID:    req.LabourID,
		Amount:      req.Amount,
		PaymentDate: paymentDate,
		PaymentType: req.PaymentType,
		Notes:       req.Notes,
	}

	if err := payment.Validate(); err != nil {
		return nil, err
	}

	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, err
	}

	return payment, nil
}

// GetByID retrieves a payment by ID
func (s *PaymentService) GetByID(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
	return s.paymentRepo.GetByID(ctx, id)
}

// GetByProjectID retrieves all payments for a project
func (s *PaymentService) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]models.PaymentWithLabour, error) {
	return s.paymentRepo.GetByProjectID(ctx, projectID)
}

// GetByLabourID retrieves all payments for a labour
func (s *PaymentService) GetByLabourID(ctx context.Context, labourID uuid.UUID) ([]models.Payment, error) {
	return s.paymentRepo.GetByLabourID(ctx, labourID)
}

// GetBalance calculates the balance for a labour in a project
func (s *PaymentService) GetBalance(ctx context.Context, projectID, labourID uuid.UUID) (*models.BalanceResponse, error) {
	return s.paymentRepo.GetBalance(ctx, projectID, labourID)
}

// Delete deletes a payment record
func (s *PaymentService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.paymentRepo.Delete(ctx, id)
}
