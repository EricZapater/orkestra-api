package operators

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type OperatorHandler struct {
	service OperatorService
}

func NewOperatorHandler(service OperatorService) *OperatorHandler {
	return &OperatorHandler{
		service: service,
	}
}

func (h *OperatorHandler) Create(c *gin.Context) {
	var request OperatorRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	operator, err := h.service.Create(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, operator)
}

func (h *OperatorHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var request OperatorRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	operator, err := h.service.Update(c.Request.Context(), id, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, operator)
}

func (h *OperatorHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.service.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *OperatorHandler) FindByID(c *gin.Context) {
	id := c.Param("id")
	operator, err := h.service.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}	
	c.JSON(http.StatusOK, operator)
}

func (h *OperatorHandler) FindAll(c *gin.Context) {
	operators, err := h.service.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, operators)
}

