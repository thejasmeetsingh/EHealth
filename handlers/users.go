package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thejasmeetsingh/EHealth/internal/database"
	"github.com/thejasmeetsingh/EHealth/models"
	"github.com/thejasmeetsingh/EHealth/validators"
)

func (apiCfg *ApiCfg) GetUserProfile(c *gin.Context, dbUser database.User) {
	SuccessResponse(c, http.StatusOK, "", models.DatabaseUserToUser(dbUser))
}

func (apiCfg *ApiCfg) UpdateUserProfile(c *gin.Context, dbUser database.User) {
	type Parameters struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	var params Parameters

	if err := c.ShouldBindJSON(&params); err != nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error while parsing the data: %v", err))
		return
	}

	if params.Name == "" {
		params.Name = dbUser.Name.String
	}

	if params.Email == "" {
		params.Email = dbUser.Email
	}

	if !validators.EmailValidator(params.Email) {
		ErrorResponse(c, http.StatusBadRequest, "Invalid email address")
		return
	}

	user, err := apiCfg.DB.UpdateUserDetails(c, database.UpdateUserDetailsParams{
		Name: sql.NullString{
			String: params.Name,
			Valid:  true,
		},
		Email: strings.ToLower(params.Email),
		ID:    dbUser.ID,
	})

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error while updating user details: %v", err))
		return
	}

	SuccessResponse(c, http.StatusOK, "", models.DatabaseUserToUser(user))
}

func (apiCfg *ApiCfg) DeleteUserProfile(c *gin.Context, dbUser database.User) {
	if err := apiCfg.DB.DeleteUser(c, dbUser.ID); err != nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error caught while deleting user profile: %v", err))
		return
	}

	SuccessResponse(c, http.StatusOK, "Profile deleted successfully!", nil)
}

func (apiCfg *ApiCfg) ChangePassword(c *gin.Context) {
	SuccessResponse(c, http.StatusOK, "", struct{}{})
}

func (apiCfg *ApiCfg) ResetPassword(c *gin.Context) {
	SuccessResponse(c, http.StatusOK, "", struct{}{})
}
