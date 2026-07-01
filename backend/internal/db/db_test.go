package db

import (
	"testing"
)

func TestClose_NilReceiver(t *testing.T) {
	var d *DB
	if err := d.Close(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestClose_EmptyDB(t *testing.T) {
	d := &DB{}
	if err := d.Close(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
