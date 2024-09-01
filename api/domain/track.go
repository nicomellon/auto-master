package domain

import (
	"errors"
	"io"

	"github.com/go-audio/audiotools"
	"github.com/google/uuid"
)

type Track struct {
	ID uuid.UUID
}

func NewTrack(file io.Reader) (*Track, error) {
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.New("Error reading file")
	}
	format, err := audiotools.HeaderFormat(fileContent)
	if err != nil {
		return nil, errors.New("Error reading file format")
	}
	if format == "unknown" {
		return nil, errors.New("Invalid file format")
	}
	return &Track{uuid.New()}, nil
}
