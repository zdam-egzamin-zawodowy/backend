package fstorage

import "io"

type FileStorage interface {
	Put(file io.Reader, filename string) error
	Remove(filename string) error
}
