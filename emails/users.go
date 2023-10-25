package emails

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/thejasmeetsingh/EHealth/utils"
)

func ResetPassword(email string, request http.Request) (bool, error) {
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

	htmlContent = utils.ReplacePlaceholders(htmlContent, map[string]string{
		"{{link}}": fmt.Sprintf("%s://%s", protocol, request.Host),
	})

	go Send([]string{email}, "Reset Password", string(htmlContent))
	return true, nil
}
