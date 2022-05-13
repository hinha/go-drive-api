package main

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	drive_v2 "google.golang.org/api/drive/v2"
	drive_v3 "google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

type authClient struct {
	cfg        *oauth2.Config
	httpClient *http.Client
}

type AuthDrive interface {
	CodeURL() string
	GetToken(authCode string) (token *oauth2.Token, err error)
	SaveToken(token *oauth2.Token)
	DriveService(ctx context.Context) (GoogleDrive, error)
}

func NewOAuthClient() AuthDrive {
	dir, _ := os.Getwd()
	b, err := ioutil.ReadFile(path.Join(dir, "client_secret.json"))
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive_v3.DriveMetadataReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	return &authClient{cfg: config}
}

func (c *authClient) CodeURL() string {
	return c.cfg.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}

func (c *authClient) GetToken(authCode string) (token *oauth2.Token, err error) {
	token, err = tokenFromFile("token.json")
	if err != nil {
		token, err = c.cfg.Exchange(context.TODO(), authCode)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve token from web %v", err)
		}
	}
	return
}

// SaveToken token to a file path.
func (c *authClient) SaveToken(token *oauth2.Token) {
	path := "token.json"
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// GetClient Retrieve a token, saves the token, then returns the generated client.
func (c *authClient) getClient(ctx context.Context) (*http.Client, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		return nil, fmt.Errorf("not found token")
	}
	c.httpClient = c.cfg.Client(ctx, tok)
	return c.httpClient, nil
}

func (c *authClient) DriveService(ctx context.Context) (GoogleDrive, error) {
	httpClient, err := c.getClient(ctx)
	if err != nil {
		return nil, err
	}
	srv2, err := drive_v2.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}
	srv3, err := drive_v3.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}
	return NewDriveService(srv3, srv2), nil
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}
