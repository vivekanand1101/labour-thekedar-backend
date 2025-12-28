package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vivekanand/labour-thekedar-backend/internal/models"
	"github.com/vivekanand/labour-thekedar-backend/internal/service"
)

// LabourHandler handles labour endpoints
type LabourHandler struct {
	labourService  *service.LabourService
	projectService *service.ProjectService
}

// NewLabourHandler creates a new LabourHandler
func NewLabourHandler(labourService *service.LabourService, projectService *service.ProjectService) *LabourHandler {
	return &LabourHandler{
		labourService:  labourService,
		projectService: projectService,
	}
}

// List handles GET /api/v1/labours
func (h *LabourHandler) List(c *gin.Context) {
	labours, err := h.labourService.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list labours"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"labours": labours})
}

// Create handles POST /api/v1/labours
func (h *LabourHandler) Create(c *gin.Context) {
	var req models.CreateLabourRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	labour, err := h.labourService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create labour"})
		return
	}

	c.JSON(http.StatusCreated, labour)
}

// Get handles GET /api/v1/labours/:id
func (h *LabourHandler) Get(c *gin.Context) {
	labourID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid labour ID"})
		return
	}

	labour, err := h.labourService.GetByID(c.Request.Context(), labourID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "labour not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get labour"})
		return
	}

	c.JSON(http.StatusOK, labour)
}

// Update handles PUT /api/v1/labours/:id
func (h *LabourHandler) Update(c *gin.Context) {
	labourID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid labour ID"})
		return
	}

	var req models.UpdateLabourRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	labour, err := h.labourService.Update(c.Request.Context(), labourID, &req)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "labour not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update labour"})
		return
	}

	c.JSON(http.StatusOK, labour)
}

// Delete handles DELETE /api/v1/labours/:id
func (h *LabourHandler) Delete(c *gin.Context) {
	labourID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid labour ID"})
		return
	}

	if err := h.labourService.Delete(c.Request.Context(), labourID); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "labour not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete labour"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "labour deleted successfully"})
}

// AssignToProject handles POST /api/v1/projects/:id/labours
func (h *LabourHandler) AssignToProject(c *gin.Context) {
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

	var req models.AssignLabourRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.labourService.AssignToProject(c.Request.Context(), projectID, req.LabourID); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "labour not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to assign labour"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "labour assigned successfully"})
}

// RemoveFromProject handles DELETE /api/v1/projects/:id/labours/:labour_id
func (h *LabourHandler) RemoveFromProject(c *gin.Context) {
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

	if err := h.labourService.RemoveFromProject(c.Request.Context(), projectID, labourID); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "labour not assigned to project"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove labour"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "labour removed from project successfully"})
}

// ListByProject handles GET /api/v1/projects/:id/labours
func (h *LabourHandler) ListByProject(c *gin.Context) {
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

	labours, err := h.labourService.GetByProjectID(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list labours"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"labours": labours})
}
