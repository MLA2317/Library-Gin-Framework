package utils

import (
	"library-project/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleError handles AppError and returns appropriate JSON response
func HandleError(c *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok {
		// Log the error
		LogError(appErr.Err, appErr.Message, map[string]interface{}{
			"code":        appErr.Code,
			"status_code": appErr.StatusCode,
			"path":        c.Request.URL.Path,
		})

		c.JSON(appErr.StatusCode, dto.ErrorResponse{
			Error:   appErr.Message,
			Code:    appErr.Code,
			Message: appErr.Error(),
		})
		return
	}

	// Generic error
	LogError(err, "Unhandled error", map[string]interface{}{
		"path": c.Request.URL.Path,
	})

	c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
		Error:   "internal server error",
		Code:    ErrCodeInternalServer,
		Message: err.Error(),
	})
}

// SuccessResponse sends a success response
func SuccessResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, dto.SuccessResponse{
		Message: "success",
		Data:    data,
	})
}

// MessageResponse sends a message-only response
func MessageResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"message": message})
}
