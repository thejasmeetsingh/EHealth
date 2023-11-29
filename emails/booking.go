package emails

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/thejasmeetsingh/EHealth/internal/database"
	"github.com/thejasmeetsingh/EHealth/utils"
)

func SendBookingCreationEmail(user database.User, booking database.Booking, request http.Request) (bool, error) {
	// Generate booking create template absolute path
	path, err := filepath.Abs(filepath.Join("emails", "templates", "booking_create.html"))
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

	// Replace the placeholders and load the dynamic content in the template
	htmlContent = utils.ReplacePlaceholders(htmlContent, map[string]string{
		"{{user_name}}":  user.Name.String,
		"{{user_email}}": user.Email,
		"{{start_dt}}":   booking.StartDatetime.Format(time.RFC822),
		"{{end_dt}}":     booking.EndDatetime.Format(time.RFC822),
	})

	// This default email is added for development purposes only
	defaultRecipientEmail := os.Getenv("DEFAULT_RECIPIENT_EMAIL")

	go Send([]string{defaultRecipientEmail}, "New Booking Request", string(htmlContent))

	return true, nil
}

func SendBookingAcceptedEmail(payload map[string]string, request http.Request) (bool, error) {
	// Generate booking create template absolute path
	path, err := filepath.Abs(filepath.Join("emails", "templates", "booking_accepted.html"))
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

	// Replace the placeholders and load the dynamic content in the template
	htmlContent = utils.ReplacePlaceholders(htmlContent, map[string]string{
		"{{name}}":     payload["name"],
		"{{address}}":  payload["address"],
		"{{start_dt}}": payload["start_dt"],
		"{{end_dt}}":   payload["end_dt"],
	})

	// This default email is added for development purposes only
	defaultRecipientEmail := os.Getenv("DEFAULT_RECIPIENT_EMAIL")

	go Send([]string{defaultRecipientEmail}, "Booking request accepted!", string(htmlContent))

	return true, nil
}

func SendBookingRejectedEmail(payload map[string]string, request http.Request) (bool, error) {
	// Generate booking create template absolute path
	path, err := filepath.Abs(filepath.Join("emails", "templates", "booking_rejected.html"))
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

	// Replace the placeholders and load the dynamic content in the template
	htmlContent = utils.ReplacePlaceholders(htmlContent, map[string]string{
		"{{name}}":     payload["name"],
		"{{address}}":  payload["address"],
		"{{start_dt}}": payload["start_dt"],
		"{{end_dt}}":   payload["end_dt"],
	})

	// This default email is added for development purposes only
	defaultRecipientEmail := os.Getenv("DEFAULT_RECIPIENT_EMAIL")

	go Send([]string{defaultRecipientEmail}, "Booking request rejected", string(htmlContent))

	return true, nil
}
