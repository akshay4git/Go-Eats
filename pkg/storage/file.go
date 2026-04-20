package storage

import (
	"io"
	"os"
)

type ImageStorage interface {
	Upload(fileName string, file io.Reader) (string, error)
}

func CreateImageStorage(storageType string) ImageStorage {
	switch storageType {
	case "local":
		basePath := os.Getenv("LOCAL_STORAGE_PATH")
		if basePath == "" {
			basePath = "/app/uploads"
		}
		localFileStore := &LocalFileStorage{BasePath: basePath}
		createUploadDirectory(localFileStore.BasePath)
		return localFileStore
	default:
		panic("Unsupported storage type: " + storageType)
	}
}