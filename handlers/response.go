package handlers

import (
	"github.com/gin-gonic/gin"
)

func successResponse(c *gin.Context, code int, message string, payload interface{}) {
	type Response struct {
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}

	c.JSON(code, Response{
		Message: message,
		Data:    payload,
	})

}

func errorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"message": message})
}
