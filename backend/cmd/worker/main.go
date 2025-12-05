package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/divijg19/physiolink/backend/internal/activities"
	"github.com/divijg19/physiolink/backend/internal/workflows"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort: "localhost:7233",
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "appointment-task-queue", worker.Options{})

	w.RegisterWorkflow(workflows.BookingWorkflow)
	w.RegisterActivity(activities.SendConfirmationEmail)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
