package filestorage

type Config struct {
	BasePath string
}

func getDefaultConfig() *Config {
	return &Config{
		BasePath: "./",
	}
}
