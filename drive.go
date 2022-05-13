package main

import (
	"fmt"
	drive_v2 "google.golang.org/api/drive/v2"
	drive_v3 "google.golang.org/api/drive/v3"
)

type googleDrive struct {
	service_v3 *drive_v3.Service
	service_v2 *drive_v2.Service
}

type GoogleDrive interface {
	FilesList(pageSize int64) (*drive_v3.FileList, error)
	AllChildren(folderID string) ([]*drive_v2.ChildReference, error)
	FileDetails(fileID string) (*drive_v3.File, error)
}

func NewDriveService(service_v3 *drive_v3.Service, service_v2 *drive_v2.Service) GoogleDrive {
	return &googleDrive{service_v3: service_v3, service_v2: service_v2}
}

func (s *googleDrive) FilesList(pageSize int64) (*drive_v3.FileList, error) {
	drv, err := s.service_v3.Files.
		List().
		PageSize(pageSize).
		Fields("nextPageToken, files(id, name, owners)").
		Do()

	if err != nil {
		return nil, fmt.Errorf("unable to retrieve files: %v", err)
	}

	return drv, nil
}

// AllChildren fetches all the children of a given folder
func (s *googleDrive) AllChildren(folderID string) ([]*drive_v2.ChildReference, error) {
	var cs []*drive_v2.ChildReference
	var pageToken string
	for {
		q := s.service_v2.Children.List(folderID)
		if pageToken != "" {
			q = q.PageToken(pageToken)
		}

		read, err := q.Do()
		if err != nil {
			return cs, err
		}

		cs = append(cs, read.Items...)
		pageToken = read.NextPageToken
		if pageToken == "" {
			break
		}
	}
	return cs, nil
}

// FileDetails fetches the given file
func (s *googleDrive) FileDetails(fileID string) (*drive_v3.File, error) {
	file, err := s.service_v3.Files.Get(fileID).Do()
	if err != nil {
		return nil, err
	}

	return file, nil
}
