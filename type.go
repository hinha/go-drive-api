package main

import "google.golang.org/api/drive/v3"

type (
	Response struct {
		Data struct {
			Files []*drive.File `json:"files"`
		} `json:"data"`
	}
)
