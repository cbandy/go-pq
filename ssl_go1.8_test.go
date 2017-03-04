// +build go1.8

package pq

import (
	"context"
	"testing"
	"time"
)

func TestIssue570(t *testing.T) {
	maybeSkipSSLTests(t)

	db, err := openTestConnConninfo("host=postgres sslmode=require user=pqgossltest")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	// Should not panic
	db.ExecContext(ctx, "select pg_sleep(1)")
}
