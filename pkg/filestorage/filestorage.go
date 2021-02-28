package filestorage

import (
	"io"
	"os"
	"path"

	"github.com/pkg/errors"
)

type fileStorage struct {
	basePath string
}

func New(cfg *Config) FileStorage {
	if cfg == nil {
		cfg = getDefaultConfig()
	}

	return &fileStorage{
		cfg.BasePath,
	}
}

func (storage *fileStorage) Put(file io.Reader, filename string) error {
	fullPath := path.Join(storage.basePath, filename)
	f, err := os.Create(fullPath)
	if err != nil {
		return errors.Wrap(err, "FileStorage.Put")
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		return errors.Wrap(err, "FileStorage.Put")
	}
	return nil
}

func (storage *fileStorage) Delete(filename string) error {
	fullPath := path.Join(storage.basePath, filename)
	err := os.Remove(fullPath)
	if err != nil {
		return errors.Wrap(err, "FileStorage.Delete")
	}
	return nil
}
