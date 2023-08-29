package util_test

import (
	"os"
	"testing"

	"github.com/glifio/cli/util"
)

const testFilename = "test_data.toml"

var defaultMap = map[string]string{
	"key1": "",
	"key2": "",
}

func TestNewStorage(t *testing.T) {
	_, err := util.NewStorage(testFilename, defaultMap, true)
	if err != nil {
		t.Errorf("NewStorage() error: %v", err)
	}

	// Cleanup
	os.Remove(testFilename)
}

func TestSetNonExistentKey(t *testing.T) {
	store, err := util.NewStorage(testFilename, defaultMap, true)
	if err != nil {
		t.Fatalf("NewStorage() error: %v", err)
	}

	err = store.Set("key3", "value3")
	if err != nil {
		t.Errorf("Set() error: %v", err)
		t.FailNow()
	}

	v3, err := store.Get("key3")
	if err != nil {
		t.Errorf("Get() error: %v", err)
		t.FailNow()
	}
	if v3 != "value3" {
		t.Errorf("Get returned the wrong value")
		t.FailNow()
	}
}

func TestGetSetDelete(t *testing.T) {
	store, err := util.NewStorage(testFilename, defaultMap, true)
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
	store, err := util.NewStorage(testFilename, defaultMap, true)
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
