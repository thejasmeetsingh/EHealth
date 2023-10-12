package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thejasmeetsingh/EHealth/internal/database"
	"github.com/thejasmeetsingh/EHealth/utils"
)

func (apiCfg *ApiCfg) Singup(c *gin.Context) {
	type Parameters struct {
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required"`
		IsEndUser bool   `json:"is_end_user" binding:"required"`
	}

	var params Parameters
	err := c.ShouldBindJSON(&params)

	if err != nil {
		errorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error while parsing the request; %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(c, database.CreateUserParams{
		ID:         uuid.New(),
		CreatedAt:  time.Now().UTC(),
		ModifiedAt: time.Now().UTC(),
		Email:      params.Email,
		Password:   params.Password,
		IsEndUser:  params.IsEndUser,
	})

	if err != nil {
		errorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error while creating a user: %v", err))
		return
	}

	tokens, err := utils.GenerateTokens(user.ID)

	if err != nil {
		errorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Error while creating a user: %v", err))
		return
	}

	successResponse(c, http.StatusCreated, "Account created successfully!", tokens)
}
