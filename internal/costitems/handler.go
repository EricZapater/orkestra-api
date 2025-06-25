package costitems

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CostItemHandler struct {
	service CostItemService
}

func NewCostItemHandler(service CostItemService) *CostItemHandler {
	return &CostItemHandler{
		service: service,
	}
}

func (h *CostItemHandler) Create(c *gin.Context){
	var request CostItemRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.service.Create(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *CostItemHandler) Update (c *gin.Context){
	id := c.Param("id")
	var request CostItemRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.service.Update(c.Request.Context(),id,  &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *CostItemHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.service.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func(h *CostItemHandler)GetByID(c *gin.Context){
	id := c.Param("id")
	data, err := h.service.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func(h *CostItemHandler)GetByProjectID(c *gin.Context){
	id := c.Param("projectid")
	data, err := h.service.FindByProjectID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func(h *CostItemHandler)GetAll(c *gin.Context){	
	data, err := h.service.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}