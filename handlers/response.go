// A base response functions for sending success or error response based on the requirement

package handlers

import (
	"github.com/gin-gonic/gin"
)

func SuccessResponse(c *gin.Context, code int, message string, payload interface{}) {
	type Response struct {
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}

	c.JSON(code, Response{
		Message: message,
		Data:    payload,
	})

}

func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"message": message})
}
