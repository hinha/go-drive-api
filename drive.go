package main

import (
	"fmt"
	"google.golang.org/api/drive/v3"
)

type googleDrive struct {
	service *drive.Service
}

type GoogleDrive interface {
	FilesList(pageSize int64) (*drive.FileList, error)
}

func NewDriveService(service *drive.Service) GoogleDrive {
	return &googleDrive{service: service}
}

func (s *googleDrive) FilesList(pageSize int64) (*drive.FileList, error) {
	drv, err := s.service.Files.
		List().
		PageSize(pageSize).
		Fields("nextPageToken, files(id, name, owners)").
		Do()

	if err != nil {
		return nil, fmt.Errorf("unable to retrieve files: %v", err)
	}

	return drv, nil
}
