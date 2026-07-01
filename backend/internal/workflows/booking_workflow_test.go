package workflows

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"

	"github.com/divijg19/physiolink/backend/internal/activities"
)

func TestBookingWorkflow_Success(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	env.OnActivity(activities.SendConfirmationEmail, mock.Anything, mock.Anything).Return("Email sent", nil)

	param := BookingWorkflowParam{
		AppointmentID: "appt-123",
		PatientID:     "pat-456",
		TherapistID:   "th-789",
	}

	env.ExecuteWorkflow(BookingWorkflow, param)

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
}

func TestBookingWorkflow_ActivityFails(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	errActivity := errors.New("email service unavailable")
	env.OnActivity(activities.SendConfirmationEmail, mock.Anything, mock.Anything).Return("", errActivity)

	param := BookingWorkflowParam{
		AppointmentID: "appt-456",
		PatientID:     "pat-789",
		TherapistID:   "th-012",
	}

	env.ExecuteWorkflow(BookingWorkflow, param)

	require.True(t, env.IsWorkflowCompleted())
	require.Error(t, env.GetWorkflowError())
}
