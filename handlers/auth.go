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

// SignUp API
func (apiCfg *ApiCfg) Singup(c *gin.Context) {
	type Parameters struct {
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required"`
		IsEndUser bool   `json:"is_end_user" binding:"required"`
	}

	var params Parameters
	err := c.ShouldBindJSON(&params)

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error while parsing the request; %v", err.Error()))
		return
	}

	// Validate Password
	err = validators.PasswordValidator(params.Password, params.Email)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid password format; %v", err))
		return
	}

	// Generate hashed password
	hashedPassword, err := utils.GetHashedPassword(params.Password)

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error while parsing the request; %v", err))
		return
	}

	// Create user account
	user, err := apiCfg.DB.CreateUser(c, database.CreateUserParams{
		ID:         uuid.New(),
		CreatedAt:  time.Now().UTC(),
		ModifiedAt: time.Now().UTC(),
		Email:      strings.ToLower(params.Email),
		Password:   hashedPassword,
		IsEndUser:  params.IsEndUser,
	})

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error while creating a user: %v", err))
		return
	}

	// Generate auth tokens for the user
	tokens, err := utils.GenerateTokens(user.ID.String())

	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Error while generating tokens: %v", err))
		return
	}

	SuccessResponse(c, http.StatusCreated, "Account created successfully!", tokens)
}

// Login API
func (apiCfg *ApiCfg) Login(c *gin.Context) {
	type Parameters struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	var params Parameters
	err := c.ShouldBindJSON(&params)

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error while parsing the request; %v", err.Error()))
		return
	}

	// Check wheather the user exists with the given email or not
	user, err := apiCfg.DB.GetUserByEmail(c, strings.ToLower(params.Email))
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "User does not exists, Please check your credentials")
		return
	}

	// Check the given password with hashed password stored in DB
	match, err := utils.CheckPassowrdValid(params.Password, user.Password)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error caught while validating password: %v", err))
		return
	} else if !match {
		ErrorResponse(c, http.StatusForbidden, "Invalid Credentials")
		return
	}

	// Generate auth tokens for the user
	tokens, err := utils.GenerateTokens(user.ID.String())

	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Error while generating tokens: %v", err))
		return
	}

	SuccessResponse(c, http.StatusCreated, "Logged in Successfully!", tokens)
}

// Refresh Token API
//
// Generate new tokens if the given refresh token is valid
func (apiCfg *ApiCfg) RefreshAccessToken(c *gin.Context) {
	type Parameters struct {
		RefreshToken string `json:"refresh_token"`
	}

	var params Parameters
	err := c.ShouldBindJSON(&params)

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error while parsing the request: %v", err))
		return
	}

	tokens, err := utils.ReIssueAccessToken(params.RefreshToken)
	if err != nil {
		ErrorResponse(c, http.StatusForbidden, fmt.Sprintf("Error while issueing new tokens: %v", err))
		return
	}

	SuccessResponse(c, http.StatusOK, "Tokens re-issued Successfully!", tokens)
}
