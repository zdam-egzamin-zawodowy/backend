package filestorage

import "os"

type Config struct {
	BasePath string
}

func getDefaultConfig() *Config {
	path, _ := os.Getwd()
	return &Config{
		BasePath: path,
	}
}
