package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thejasmeetsingh/EHealth/emails"
	"github.com/thejasmeetsingh/EHealth/internal/database"
	"github.com/thejasmeetsingh/EHealth/models"
	"github.com/thejasmeetsingh/EHealth/utils"
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

func (apiCfg *ApiCfg) ChangePassword(c *gin.Context, dbUser database.User) {
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

	match, err := utils.CheckPassowrdValid(params.CurrentPassword, dbUser.Password)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error caught while validating password: %v", err))
		return
	} else if !match {
		ErrorResponse(c, http.StatusBadRequest, "Current password is incorrect")
		return
	}

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

	hashedPassword, err := utils.GetHashedPassword(params.NewPasswordConfirm)

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error while parsing the request; %v", err))
		return
	}

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

func (apiCfg *ApiCfg) ResetPassword(c *gin.Context) {
	type Parameters struct {
		Email string `json:"email" binding:"required,email"`
	}
	var params Parameters

	if err := c.ShouldBindJSON(&params); err != nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error while parsing the data: %v", err))
		return
	}

	dbUser, err := apiCfg.DB.GetUserByEmail(c, strings.ToLower(params.Email))
	if err != nil {
		ErrorResponse(c, http.StatusForbidden, "Email does not exists")
		return
	}

	_, err = emails.ResetPassword(dbUser.Email, *c.Request)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Error while sending email: %v", err))
		return
	}

	SuccessResponse(c, http.StatusOK, "Reset password email sent successfully", nil)
}

func (apiCfg *ApiCfg) RenderResetPassword(c *gin.Context) {
	c.HTML(http.StatusOK, "reset_password.html", gin.H{
		"token":   c.Param("token"),
		"isValid": false,
	})
}

func (apiCfg *ApiCfg) ValidateResetPassword(c *gin.Context) {
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirm-password")
	token := c.PostForm("token")

	if token == "" {
		c.HTML(http.StatusBadRequest, "reset_password.html", gin.H{
			"token":   token,
			"isValid": false,
			"message": "Invalid Link or Link is expired",
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

	claims, err := utils.VerifyToken(token)
	if err != nil {
		c.HTML(http.StatusBadRequest, "reset_password.html", gin.H{
			"token":   token,
			"isValid": false,
			"message": "Invalid Link or Link is expired",
		})
		return
	}
	fmt.Println(claims.Data)
	err = validators.PasswordValidator(confirmPassword, claims.Data)
	if err != nil {
		c.HTML(http.StatusBadRequest, "reset_password.html", gin.H{
			"token":   token,
			"isValid": false,
			"message": err,
		})
		return
	}

	dbUser, err := apiCfg.DB.GetUserByEmail(c, claims.Data)
	if err != nil {
		c.HTML(http.StatusBadRequest, "reset_password.html", gin.H{
			"token":   token,
			"isValid": false,
			"message": err,
		})
		return
	}

	hashedPassword, err := utils.GetHashedPassword(confirmPassword)

	if err != nil {
		c.HTML(http.StatusBadRequest, "reset_password.html", gin.H{
			"token":   token,
			"isValid": false,
			"message": "Invalid Password",
		})
		return
	}

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
