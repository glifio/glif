package util_test

import (
	"os"
	"testing"

	"github.com/glif-confidential/cli/util"
)

const testFilename = "test_data.toml"

func TestNewStorage(t *testing.T) {
	_, err := util.NewStorage(testFilename)
	if err != nil {
		t.Errorf("NewStorage() error: %v", err)
	}

	// Cleanup
	os.Remove(testFilename)
}

func TestGetSetDelete(t *testing.T) {
	store, err := util.NewStorage(testFilename)
	if err != nil {
		t.Fatalf("NewStorage() error: %v", err)
	}

	// Test Set()
	err = store.Set("key1", "value1")
	if err != nil {
		t.Errorf("Set() error: %v", err)
	}

	// Test Get()
	value, err := store.Get("key1")
	if err != nil {
		t.Errorf("Get() error: %v", err)
	}

	if value != "value1" {
		t.Errorf("Get() expected 'value1', got '%s'", value)
	}

	// Test Delete()
	err = store.Delete("key1")
	if err != nil {
		t.Errorf("Delete() error: %v", err)
	}

	_, err = store.Get("key1")
	if err == nil {
		t.Errorf("Get() expected error, got nil")
	}

	// Cleanup
	os.Remove(testFilename)
}

func TestNonexistentKey(t *testing.T) {
	store, err := util.NewStorage(testFilename)
	if err != nil {
		t.Fatalf("NewStorage() error: %v", err)
	}

	_, err = store.Get("nonexistent_key")
	if err == nil {
		t.Errorf("Get() expected error, got nil")
	}

	err = store.Delete("nonexistent_key")
	if err == nil {
		t.Errorf("Delete() expected error, got nil")
	}

	// Cleanup
	os.Remove(testFilename)
}
