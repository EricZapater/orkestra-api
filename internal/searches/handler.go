package searches

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	SearchService SearchService
}

func NewSearchHandler(service SearchService) *SearchHandler {
	return &SearchHandler{
		SearchService: service,
	}
}

func (h *SearchHandler) GetByText(c *gin.Context){
	var request SearchRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := h.SearchService.GetByText(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}