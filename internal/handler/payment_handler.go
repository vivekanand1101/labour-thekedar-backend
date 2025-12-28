package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vivekanand/labour-thekedar-backend/internal/models"
	"github.com/vivekanand/labour-thekedar-backend/internal/service"
)

// PaymentHandler handles payment endpoints
type PaymentHandler struct {
	paymentService *service.PaymentService
	projectService *service.ProjectService
}

// NewPaymentHandler creates a new PaymentHandler
func NewPaymentHandler(paymentService *service.PaymentService, projectService *service.ProjectService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		projectService: projectService,
	}
}

// ListByProject handles GET /api/v1/projects/:id/payments
func (h *PaymentHandler) ListByProject(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project ID"})
		return
	}

	// Verify project ownership
	isOwner, err := h.projectService.IsOwner(c.Request.Context(), projectID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify ownership"})
		return
	}
	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	payments, err := h.paymentService.GetByProjectID(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list payments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payments": payments})
}

// Create handles POST /api/v1/projects/:id/payments
func (h *PaymentHandler) Create(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project ID"})
		return
	}

	// Verify project ownership
	isOwner, err := h.projectService.IsOwner(c.Request.Context(), projectID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify ownership"})
		return
	}
	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	var req models.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.paymentService.Create(c.Request.Context(), projectID, &req)
	if err != nil {
		if errors.Is(err, models.ErrInvalidDate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
			return
		}
		if errors.Is(err, models.ErrInvalidLabour) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "labour not assigned to this project"})
			return
		}
		if errors.Is(err, models.ErrInvalidAmount) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount"})
			return
		}
		if errors.Is(err, models.ErrInvalidPaymentType) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment type, use advance, daily_wage, or bonus"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create payment"})
		return
	}

	c.JSON(http.StatusCreated, payment)
}

// ListByLabour handles GET /api/v1/labours/:id/payments
func (h *PaymentHandler) ListByLabour(c *gin.Context) {
	labourID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid labour ID"})
		return
	}

	payments, err := h.paymentService.GetByLabourID(c.Request.Context(), labourID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list payments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payments": payments})
}

// GetBalance handles GET /api/v1/projects/:id/labours/:labour_id/balance
func (h *PaymentHandler) GetBalance(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project ID"})
		return
	}
	labourID, err := uuid.Parse(c.Param("labour_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid labour ID"})
		return
	}

	// Verify project ownership
	isOwner, err := h.projectService.IsOwner(c.Request.Context(), projectID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify ownership"})
		return
	}
	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	balance, err := h.paymentService.GetBalance(c.Request.Context(), projectID, labourID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "labour not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get balance"})
		return
	}

	c.JSON(http.StatusOK, balance)
}

// Delete handles DELETE /api/v1/payments/:id
func (h *PaymentHandler) Delete(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	paymentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment ID"})
		return
	}

	// Get payment to verify project ownership
	payment, err := h.paymentService.GetByID(c.Request.Context(), paymentID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get payment"})
		return
	}

	// Verify project ownership
	isOwner, err := h.projectService.IsOwner(c.Request.Context(), payment.ProjectID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify ownership"})
		return
	}
	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	if err := h.paymentService.Delete(c.Request.Context(), paymentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "payment deleted successfully"})
}
