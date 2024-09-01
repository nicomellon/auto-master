package services

import (
	"errors"
	"io"

	"github.com/go-audio/audiotools"
	"github.com/google/uuid"
)

type UploadService struct{}

func NewUploadService() *UploadService {
	return &UploadService{}
}

func (s *UploadService) Upload(file io.Reader) (string, error) {
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return "", errors.New("Error reading file")
	}
	format, err := audiotools.HeaderFormat(fileContent)
	if err != nil {
		return "", errors.New("Error reading file format")
	}
	if format == "unknown" {
		return "", errors.New("Invalid file format")
	}
	return uuid.NewString(), nil
}
