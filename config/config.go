package config

import (
	"errors"
	"os"
)

type Config struct {
	MongoURI           string
	Database           string
	FsFilesCollection  string
	FsChunksCollection string
	KafkaBroker        string
	KafkaTopic         string
}

func InitializeConfig() (*Config, error) {
	m := os.Getenv("MONGO_URI")
	if m == "" {
		return nil, errors.New("MONGO_URI cannot be empty")
	}

	d := os.Getenv("DATABASE")
	if d == "" {
		return nil, errors.New("DATABASE cannot be empty")
	}

	f := os.Getenv("FS_FILES_COLLECTION")
	if f == "" {
		return nil, errors.New("FS_FILES_COLLECTION cannot be empty")
	}

	c := os.Getenv("FS_CHUNKS_COLLECTION")
	if c == "" {
		return nil, errors.New("FS_CHUNKS_COLLECTION cannot be empty")
	}

	b := os.Getenv("KAFKA_BROKER")
	if b == "" {
		return nil, errors.New("KAFKA_BROKER is required")
	}

	t := os.Getenv("KAFKA_TOPIC")
	if t == "" {
		return nil, errors.New("KAFKA_TOPIC is required")
	}

	return &Config{
		MongoURI:           m,
		Database:           d,
		FsFilesCollection:  f,
		FsChunksCollection: c,
		KafkaBroker:        b,
		KafkaTopic:         t,
	}, nil
}
