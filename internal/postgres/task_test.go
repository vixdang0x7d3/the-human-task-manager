package postgres_test

import (
	"context"
	"testing"
)

func TestCreateTask(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := MustOpenDB(t, context.Background())
		defer CloseDB(t, db)
		t.Fatal("no test")
	})
}
