package util

import (
	"fmt"
	"io"
	"os"

	toml "github.com/pelletier/go-toml/v2"
)

type StorageData map[string]string

// Storage is a structure that holds the filename and a map of key-value pairs.
type Storage struct {
	filename string
	data     StorageData
}

// NewStorage creates a new Storage instance and initializes it with the given filename.
func NewStorage(filename string) (*Storage, error) {
	s := &Storage{
		filename: filename,
		data:     make(map[string]string),
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		err = s.save()
		if err != nil {
			return nil, err
		}
	} else {
		err := s.load()
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

// load reads the contents of the file and loads the key-value pairs into the data map.
func (s *Storage) load() error {
	file, err := os.Open(s.filename)
	if err != nil {
		if os.IsNotExist(err) {
			s.data = make(map[string]string)
			return nil
		}
		return err
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var sd StorageData

	if err := toml.Unmarshal(fileContent, &sd); err != nil {
		return err
	}

	s.data = sd

	return nil
}

// save writes the current key-value pairs in the data map to the file.
func (s *Storage) save() error {
	keyStore, err := toml.Marshal(s.data)
	if err != nil {
		return err
	}

	err = os.WriteFile(s.filename, keyStore, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Get retrieves the value associated with the given key.
func (s *Storage) Get(key string) (string, error) {
	value, ok := s.data[key]
	if !ok {
		return "", fmt.Errorf("key not found: %s", key)
	}
	return value, nil
}

// Set sets a key-value pair in the data map and saves the data to the file.
func (s *Storage) Set(key, value string) error {
	s.data[key] = value
	return s.save()
}

// Delete removes a key-value pair from the data map and saves the data to the file.
func (s *Storage) Delete(key string) error {
	if _, ok := s.data[key]; !ok {
		return fmt.Errorf("key not found: %s", key)
	}
	delete(s.data, key)
	return s.save()
}
