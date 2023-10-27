package emails

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/thejasmeetsingh/EHealth/utils"
)

// Send reset password mail to the given email address
func ResetPassword(email string, request http.Request) (bool, error) {
	// Generate reset password template absolute path
	path, err := filepath.Abs(filepath.Join("emails", "templates", "reset_password.html"))
	if err != nil {
		return false, err
	}

	htmlContent, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}

	protocol := "http"
	if request.TLS != nil {
		protocol += "s"
	}

	// Generate a JWT token encoded with user email address and token expiry
	signedToken, err := utils.GetToken(time.Now().Add(time.Hour*1), email)
	if err != nil {
		return false, err
	}

	// Replace the placeholders and load the dynamic content in the template
	htmlContent = utils.ReplacePlaceholders(htmlContent, map[string]string{
		"{{link}}": fmt.Sprintf("%s://%s/%s/%s/", protocol, request.Host, "reset-password", signedToken),
	})

	// This default email is added for development purposes only
	defaultRecipientEmail := os.Getenv("DEFAULT_RECIPIENT_EMAIL")

	go Send([]string{defaultRecipientEmail}, "Reset Password", string(htmlContent))
	return true, nil
}
