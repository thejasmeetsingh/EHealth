package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thejasmeetsingh/EHealth/internal/database"
	"github.com/thejasmeetsingh/EHealth/utils"
	"github.com/thejasmeetsingh/EHealth/validators"
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
		errorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error while parsing the request; %v", err.Error()))
		return
	}

	err = validators.PasswordValidator(params.Password, params.Email)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error while parsing the request; %v", err))
		return
	}

	hashedPassword, err := utils.GetHashedPassword(params.Password)

	if err != nil {
		errorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error while parsing the request; %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(c, database.CreateUserParams{
		ID:         uuid.New(),
		CreatedAt:  time.Now().UTC(),
		ModifiedAt: time.Now().UTC(),
		Email:      strings.ToLower(params.Email),
		Password:   hashedPassword,
		IsEndUser:  params.IsEndUser,
	})

	if err != nil {
		errorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error while creating a user: %v", err))
		return
	}

	tokens, err := utils.GenerateTokens(user.ID)

	if err != nil {
		errorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Error while generating tokens: %v", err))
		return
	}

	successResponse(c, http.StatusCreated, "Account created successfully!", tokens)
}

func (apiCfg *ApiCfg) Login(c *gin.Context) {
	type Parameters struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	var params Parameters
	err := c.ShouldBindJSON(&params)

	if err != nil {
		errorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error while parsing the request; %v", err.Error()))
		return
	}

	user, err := apiCfg.DB.GetUserByEmail(c, strings.ToLower(params.Email))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "User does not exists, Please check your credentials")
		return
	}

	match, err := utils.CheckPassowrdValid(params.Password, user.Password)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error caught while validating password: %v", err))
		return
	} else if !match {
		errorResponse(c, http.StatusForbidden, "Invalid Credentials")
		return
	}

	tokens, err := utils.GenerateTokens(user.ID)

	if err != nil {
		errorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Error while generating tokens: %v", err))
		return
	}

	successResponse(c, http.StatusCreated, "Logged in Successfully!", tokens)
}
