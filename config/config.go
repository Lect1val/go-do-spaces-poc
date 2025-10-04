package config

import "os"

type Config struct {
	DOKey      string
	DOSecret   string
	DORegion   string
	DOEndpoint string
	DOBucket   string
}

func LoadConfig() *Config {
	return &Config{
		DOKey:      os.Getenv("DO_SPACES_KEY"),
		DOSecret:   os.Getenv("DO_SPACES_SECRET"),
		DORegion:   os.Getenv("DO_SPACES_REGION"),
		DOEndpoint: os.Getenv("DO_SPACES_ENDPOINT"),
		DOBucket:   os.Getenv("DO_SPACES_BUCKET"),
	}
}
