package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/plutack/gofile-api-go-client/cmd/core/file"
	"github.com/plutack/gofile-api-go-client/cmd/model"
)

var baseUrl = "https://api.gofile.io"

type ClientConfig struct {
	apiToken   string
	baseUrl    string
	retryCount int
	timeout    time.Duration
}
type Client struct {
	httpClient *http.Client
	config     ClientConfig
}

var getMethod = "GET"
var postMethod = "POST"

func NewDefaultClientConfig() ClientConfig {
	return ClientConfig{
		apiToken:   os.Getenv("gofile_api_key"),
		baseUrl:    baseUrl,
		retryCount: 3,
		timeout:    10 * time.Second,
	}
}

func newClient(c ClientConfig) *Client {
	return &Client{
		config: c,
		httpClient: &http.Client{
			Timeout: c.timeout,
		},
	}
}

func setAuthorizationHeader(r *http.Request, t string) {
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t))
}

func (c *Client) GetAvailableServers(zone string) (*http.Response, error) {
	u, err := url.Parse(c.config.baseUrl + "/servers")
	if err != nil {
		panic(err)
	}

	q := u.Query()

	if zone != "" {
		q.Add("zone", zone)
	}
	u.RawQuery = q.Encode()
	req, err := http.NewRequest(getMethod, u.String(), nil)
	if err != nil {
		return nil, err
	}
	setAuthorizationHeader(req, c.config.apiToken)
	return c.httpClient.Do(req)
}

func (c *Client) createFolder(parentFolderId string, name string) (*http.Response, error) {
	u, err := url.Parse(c.config.baseUrl + "/contents/createFolder")
	if err != nil {
		panic(err)
	}

	payload := model.NewFolderPayload(parentFolderId, name)
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(postMethod, u.String(), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	setAuthorizationHeader(req, c.config.apiToken)
	req.Header.Set("Content-Type", "application/json")
	return c.httpClient.Do(req)
}

func (c *Client) getAccountId() (*http.Response, error) {
	u, err := url.Parse(c.config.baseUrl + "/accounts/getid")
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(getMethod, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.config.apiToken)
	return c.httpClient.Do(req)
}

func (c *Client) getAccountInformation(id string) (*http.Response, error) {
	u, err := url.Parse(c.config.baseUrl + fmt.Sprintf("/accounts/%s", id))
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(getMethod, u.String(), nil)
	if err != nil {
		return nil, err
	}
	setAuthorizationHeader(req, c.config.apiToken)
	return c.httpClient.Do(req)
}

func (c *Client) uploadFile(filePath string, folderId string) (*http.Response, error) {
	u, err := url.Parse(c.config.baseUrl + "/accounts/getid")
	if err != nil {
		panic(err)
	}

	pr, pw := io.Pipe()
	w := multipart.NewWriter(pw)

	err = file.Upload(w, filePath, folderId)
	if err != nil {
		pw.CloseWithError(err)
		return nil, err
	}
	pw.CloseWithError(w.Close())

	req, err := http.NewRequest(postMethod, u.String(), pr)
	if err != nil {
		return nil, err
	}
	setAuthorizationHeader(req, c.config.apiToken)
	req.Header.Set("Content-Type", w.FormDataContentType())
	return c.httpClient.Do(req)
	// i still need to change the timeout here
}
