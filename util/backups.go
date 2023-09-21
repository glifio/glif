package util

import "time"

type BackupsStorage struct {
	*Storage
}

var backupsStore *BackupsStorage

func BackupsStore() *BackupsStorage {
	return backupsStore
}

func NewBackupsStore(filename string) error {
	backupsDefault := map[string]string{
		"confirmed-exists": "false",
		"modified-at":      "",
	}

	s, err := NewStorage(filename, backupsDefault, true)
	if err != nil {
		return err
	}

	backupsStore = &BackupsStorage{s}

	return nil
}

func (a *BackupsStorage) Invalidate() {
	v, _ := time.Now().UTC().MarshalText()
	a.Set("modified-at", string(v))
	a.Set("confirmed-exists", "false")
}
