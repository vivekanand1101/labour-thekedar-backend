package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vivekanand/labour-thekedar-backend/internal/models"
	"github.com/vivekanand/labour-thekedar-backend/internal/service"
)

// WorkDayHandler handles work day (attendance) endpoints
type WorkDayHandler struct {
	workDayService *service.WorkDayService
	projectService *service.ProjectService
}

// NewWorkDayHandler creates a new WorkDayHandler
func NewWorkDayHandler(workDayService *service.WorkDayService, projectService *service.ProjectService) *WorkDayHandler {
	return &WorkDayHandler{
		workDayService: workDayService,
		projectService: projectService,
	}
}

// List handles GET /api/v1/projects/:id/attendance
func (h *WorkDayHandler) List(c *gin.Context) {
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

	// Check if filtering by date
	dateStr := c.Query("date")
	if dateStr != "" {
		workDays, err := h.workDayService.GetByProjectAndDate(c.Request.Context(), projectID, dateStr)
		if err != nil {
			if errors.Is(err, models.ErrInvalidDate) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list attendance"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"attendance": workDays})
		return
	}

	workDays, err := h.workDayService.GetByProjectID(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list attendance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"attendance": workDays})
}

// Create handles POST /api/v1/projects/:id/attendance
func (h *WorkDayHandler) Create(c *gin.Context) {
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

	var req models.CreateWorkDayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workDay, err := h.workDayService.Create(c.Request.Context(), projectID, &req)
	if err != nil {
		if errors.Is(err, models.ErrInvalidDate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
			return
		}
		if errors.Is(err, models.ErrInvalidLabour) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "labour not assigned to this project"})
			return
		}
		if errors.Is(err, models.ErrInvalidStatus) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status, use full_day, half_day, or absent"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create attendance record"})
		return
	}

	c.JSON(http.StatusCreated, workDay)
}

// Update handles PUT /api/v1/attendance/:id
func (h *WorkDayHandler) Update(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	workDayID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid attendance ID"})
		return
	}

	// Get work day to verify project ownership
	workDay, err := h.workDayService.GetByID(c.Request.Context(), workDayID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "attendance record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get attendance record"})
		return
	}

	// Verify project ownership
	isOwner, err := h.projectService.IsOwner(c.Request.Context(), workDay.ProjectID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify ownership"})
		return
	}
	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	var req models.UpdateWorkDayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedWorkDay, err := h.workDayService.Update(c.Request.Context(), workDayID, &req)
	if err != nil {
		if errors.Is(err, models.ErrInvalidStatus) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update attendance record"})
		return
	}

	c.JSON(http.StatusOK, updatedWorkDay)
}

// Delete handles DELETE /api/v1/attendance/:id
func (h *WorkDayHandler) Delete(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	workDayID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid attendance ID"})
		return
	}

	// Get work day to verify project ownership
	workDay, err := h.workDayService.GetByID(c.Request.Context(), workDayID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "attendance record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get attendance record"})
		return
	}

	// Verify project ownership
	isOwner, err := h.projectService.IsOwner(c.Request.Context(), workDay.ProjectID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify ownership"})
		return
	}
	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	if err := h.workDayService.Delete(c.Request.Context(), workDayID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete attendance record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "attendance record deleted successfully"})
}
