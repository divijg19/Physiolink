package integration

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// uniqueEmail builds a per-run email so repeated runs against the same database
// do not collide on the users.email unique index.
func uniqueEmail(prefix string) string {
	return fmt.Sprintf("%s-%s@example.com", prefix, uuid.NewString())
}

// pgxTextArray adapts a list of strings for a Postgres text[] parameter.
func pgxTextArray(values ...string) interface{} {
	return pq.Array(values)
}
