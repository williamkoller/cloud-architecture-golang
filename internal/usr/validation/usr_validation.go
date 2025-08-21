package validation

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RespondValidationError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": err.Error(),
	})
}
