package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetOffset(c *gin.Context) int32 {
	offsetStr := c.Param("offset")
	if offsetStr == "" {
		offsetStr = "0"
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	return int32(offset)
}
