package workflows

import (
	"time"

	"go.temporal.io/sdk/workflow"

	"github.com/divijg19/physiolink/backend/go/internal/activities"
)

// BookingWorkflowParam is the parameter passed to the workflow
type BookingWorkflowParam struct {
	AppointmentID string
	PatientID     string
	TherapistID   string
}

// BookingWorkflow orchestrates the booking process
func BookingWorkflow(ctx workflow.Context, param BookingWorkflowParam) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("Booking workflow started", "AppointmentID", param.AppointmentID)

	// Execute Activity: Send Confirmation Email
	var result string
	err := workflow.ExecuteActivity(ctx, activities.SendConfirmationEmail, param.AppointmentID).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed", "Error", err)
		return err
	}

	logger.Info("Booking workflow completed", "Result", result)
	return nil
}
