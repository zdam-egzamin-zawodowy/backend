package filestorage

import "io"

type FileStorage interface {
	Put(file io.Reader, filename string) error
	Delete(filename string) error
}
