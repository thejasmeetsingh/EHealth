package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thejasmeetsingh/EHealth/emails"
	"github.com/thejasmeetsingh/EHealth/internal/database"
	"github.com/thejasmeetsingh/EHealth/models"
	"github.com/thejasmeetsingh/EHealth/utils"
	"github.com/thejasmeetsingh/EHealth/validators"
)

// Fetch user profile details
func (apiCfg *ApiCfg) GetUserProfile(c *gin.Context) {
	dbUser, err := getDBUser(c)
	if err != nil {
		ErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "", models.DatabaseUserToUser(dbUser))
}

// Update user profile details
func (apiCfg *ApiCfg) UpdateUserProfile(c *gin.Context) {
	dbUser, err := getDBUser(c)
	if err != nil {
		ErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

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

	// Validate the given email address
	if !validators.EmailValidator(params.Email) {
		ErrorResponse(c, http.StatusBadRequest, "Invalid email address")
		return
	}

	// Update the profile details
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

// Delete user profile
func (apiCfg *ApiCfg) DeleteUserProfile(c *gin.Context) {
	dbUser, err := getDBUser(c)
	if err != nil {
		ErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	if err := apiCfg.DB.DeleteUser(c, dbUser.ID); err != nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error caught while deleting user profile: %v", err))
		return
	}

	SuccessResponse(c, http.StatusOK, "Profile deleted successfully!", nil)
}

// Change password API for authenticated user
func (apiCfg *ApiCfg) ChangePassword(c *gin.Context) {
	dbUser, err := getDBUser(c)
	if err != nil {
		ErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	type Parameters struct {
		CurrentPassword    string `json:"current_password"`
		NewPassword        string `json:"new_password"`
		NewPasswordConfirm string `json:"new_password_confirm"`
	}
	var params Parameters

	if err := c.ShouldBindJSON(&params); err != nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error while parsing the data: %v", err))
		return
	}

	// Check wheather or not current password is correct or not
	match, err := utils.CheckPassowrdValid(params.CurrentPassword, dbUser.Password)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error caught while validating password: %v", err))
		return
	} else if !match {
		ErrorResponse(c, http.StatusBadRequest, "Current password is incorrect")
		return
	}

	// Validate the new password
	if params.NewPassword != params.NewPasswordConfirm {
		fmt.Println(params.NewPassword, params.NewPasswordConfirm)
		ErrorResponse(c, http.StatusBadRequest, "New password does not match with new password confirm")
		return
	}

	if params.NewPasswordConfirm == params.CurrentPassword {
		ErrorResponse(c, http.StatusBadRequest, "New password should not match the current password")
		return
	}

	if err = validators.PasswordValidator(params.NewPasswordConfirm, dbUser.Email); err != nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid password format; %v", err))
		return
	}

	// Generate the new hashed password
	hashedPassword, err := utils.GetHashedPassword(params.NewPasswordConfirm)

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error while parsing the request; %v", err))
		return
	}

	// Update the password
	_, err = apiCfg.DB.UpdateUserPassword(c, database.UpdateUserPasswordParams{
		Password: hashedPassword,
		ID:       dbUser.ID,
	})

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error while updating password; %v", err))
		return
	}

	SuccessResponse(c, http.StatusOK, "Password changed successfully!", nil)
}

// Reset password API for users who forgot their password
func (apiCfg *ApiCfg) ResetPassword(c *gin.Context) {
	type Parameters struct {
		Email string `json:"email" binding:"required,email"`
	}
	var params Parameters

	if err := c.ShouldBindJSON(&params); err != nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error while parsing the data: %v", err))
		return
	}

	// Fetch the user with the given email address
	dbUser, err := apiCfg.DB.GetUserByEmail(c, strings.ToLower(params.Email))
	if err != nil {
		ErrorResponse(c, http.StatusForbidden, "Email does not exists")
		return
	}

	// Send reset password email
	_, err = emails.ResetPassword(dbUser.Email, *c.Request)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Error while sending email: %v", err))
		return
	}

	SuccessResponse(c, http.StatusOK, "Reset password email sent successfully", nil)
}

// Render the reset password form
func (apiCfg *ApiCfg) RenderResetPassword(c *gin.Context) {
	c.HTML(http.StatusOK, "reset_password.html", gin.H{
		"token":   c.Param("token"),
		"isValid": false,
	})
}

// Validate the new password coming from the reset password form
func (apiCfg *ApiCfg) ValidateResetPassword(c *gin.Context) {
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirm-password")
	token := c.PostForm("token")

	if token == "" {
		c.HTML(http.StatusBadRequest, "reset_password.html", gin.H{
			"token":   token,
			"isValid": false,
			"message": "Invalid Link",
		})
		return
	}

	if password != confirmPassword {
		c.HTML(http.StatusBadRequest, "reset_password.html", gin.H{
			"token":   token,
			"isValid": false,
			"message": "Password do not match with Confirm Password",
		})
		return
	}

	// Verify the reset password token and get encoded token data
	claims, err := utils.VerifyToken(token)
	if err != nil {
		c.HTML(http.StatusBadRequest, "reset_password.html", gin.H{
			"token":   token,
			"isValid": false,
			"message": "Invalid Link",
		})
		return
	}

	// Check the validity of the token
	if !time.Unix(claims.ExpiresAt.Unix(), 0).After(time.Now()) {
		c.HTML(http.StatusBadRequest, "reset_password.html", gin.H{
			"token":   token,
			"isValid": false,
			"message": "Link is expired",
		})
		return
	}

	// Validate the new password
	err = validators.PasswordValidator(confirmPassword, claims.Data)
	if err != nil {
		c.HTML(http.StatusBadRequest, "reset_password.html", gin.H{
			"token":   token,
			"isValid": false,
			"message": err,
		})
		return
	}

	// Fetch user by the email which was encoded in the token
	dbUser, err := apiCfg.DB.GetUserByEmail(c, claims.Data)
	if err != nil {
		c.HTML(http.StatusBadRequest, "reset_password.html", gin.H{
			"token":   token,
			"isValid": false,
			"message": err,
		})
		return
	}

	// Generate the new hashed password
	hashedPassword, err := utils.GetHashedPassword(confirmPassword)

	if err != nil {
		c.HTML(http.StatusBadRequest, "reset_password.html", gin.H{
			"token":   token,
			"isValid": false,
			"message": "Invalid Password",
		})
		return
	}

	// Update the password of the user
	_, err = apiCfg.DB.UpdateUserPassword(c, database.UpdateUserPasswordParams{
		Password: hashedPassword,
		ID:       dbUser.ID,
	})

	if err != nil {
		c.HTML(http.StatusBadRequest, "reset_password.html", gin.H{
			"token":   token,
			"isValid": false,
			"message": "Cannot update the password",
		})
		return
	}

	c.HTML(http.StatusBadRequest, "reset_password.html", gin.H{
		"isValid": true,
		"message": "Password updated successfully!",
	})
}
