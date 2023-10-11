package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thejasmeetsingh/EHealth/handlers/parameters"
)

func (apiCfg *ApiCfg) Singup(c *gin.Context) {
	var params parameters.AuthParameters
	err := c.ShouldBindJSON(&params)

	if err == nil {
		successResponse(c, http.StatusOK, "Account created successfully!", params)
		return
	}
	errorResponse(c, http.StatusBadRequest, "Error while parsing the request")
}
