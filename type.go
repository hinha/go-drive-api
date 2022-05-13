package main

import (
	drive_v2 "google.golang.org/api/drive/v2"
	drive_v3 "google.golang.org/api/drive/v3"
)

type (
	Response struct {
		Data struct {
			Files   []*drive_v3.File           `json:"files,omitempty"`
			Folders []*drive_v2.ChildReference `json:"folders,omitempty"`
			File    *drive_v3.File             `json:"file,omitempty"`
		} `json:"data"`
	}
)
