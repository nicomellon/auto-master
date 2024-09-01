package services

import (
	"io"

	"github.com/nicomellon/auto-master/domain"
)

type UploadService struct{}

func NewUploadService() *UploadService {
	return &UploadService{}
}

func (s *UploadService) Upload(file io.Reader) (*domain.Track, error) {
	track, err := domain.NewTrack(file)
	if err != nil {
		return nil, err
	}
	return track, nil
}
