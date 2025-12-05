package service

import (
	"go.temporal.io/sdk/client"
)

// NewTemporalClient creates a Temporal client to talk to the local server.
// Call Close() on the returned client when shutting down.
func NewTemporalClient() (client.Client, error) {
	return client.Dial(client.Options{
		HostPort: "localhost:7233",
	})
}
