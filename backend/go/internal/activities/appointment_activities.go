package activities

import (
	"context"
	"fmt"
)

// SendConfirmationEmail sends a confirmation email to the patient
func SendConfirmationEmail(ctx context.Context, appointmentID string) (string, error) {
	// In a real app, this would use an email service (e.g., SendGrid, AWS SES)
	fmt.Printf("Sending confirmation email for appointment: %s\n", appointmentID)
	return "Email sent", nil
}
