package main

import (
	"fmt"
	drive_v2 "google.golang.org/api/drive/v2"
	drive_v3 "google.golang.org/api/drive/v3"
	"io/ioutil"
	"mime/multipart"
)

type googleDrive struct {
	service_v3 *drive_v3.Service
	service_v2 *drive_v2.Service
}

type GoogleDrive interface {
	FilesList(pageSize int64) (*drive_v3.FileList, error)
	AllChildren(folderID string) ([]*drive_v2.ChildReference, error)
	FileDetails(fileID string) (*drive_v3.File, error)
	Upload(file multipart.File, filename string, folderId string) (*drive_v3.File, error)
	Download(fileId string) ([]byte, string, error)
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

func (s *googleDrive) Upload(file multipart.File, filename string, folderId string) (*drive_v3.File, error) {
	driveFile, err := s.service_v3.Files.Create(&drive_v3.File{Name: filename, Parents: []string{folderId}}).Media(file).Do()
	if err != nil {
		return nil, err
	}
	return driveFile, nil
}

// Download of files stored in Google Drive
// fileId provide file location
// data, content-type
func (s *googleDrive) Download(fileId string) ([]byte, string, error) {
	response, err := s.service_v3.Files.Get(fileId).Download()
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, "", err
	}

	return body, response.Header.Get("Content-Type"), nil
}
