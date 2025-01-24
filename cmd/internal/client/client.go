// package client contains low level abstractions to wrap around gofile.io api
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/plutack/gofile-api-go-client/cmd/core/file"
	"github.com/plutack/gofile-api-go-client/cmd/model"
)

// baseUrl is the base  URL used for gofile.io api calls
var baseUrl = "https://api.gofile.io"

// clientConfig contains necessary configuration options to configure a client
type ClientConfig struct {
	// APIToken is the authentication token for the GoFile.io API
	APIToken string
	//BaseUrl is the base url for API request apart from uploadFile API call
	BaseUrl string
	// RetryCount specifies the number of times to retry failed API requests
	RetryCount int
	// Timeout specifies the maximum time to wait for an API Request to be resolved
	Timeout time.Duration
}

// Client represents an HTTP client for interacting with the GoFile.io API
type Client struct {
	// httpClient is the underlying HTTP client used for API requests
	httpClient *http.Client
	// config holds the configuration settings for the API client
	config ClientConfig
}

// HTTP request methods for API interactions
const (
	getMethod  = "GET"
	postMethod = "POST"
)

// NewDefaultClientConfig creates a default ClientConfig with preset values
// - API token from environment variable
// - Default base URL
// - 3 retry attempts
// - 10-second timeout
func NewDefaultClientConfig() ClientConfig {
	return ClientConfig{
		APIToken:   os.Getenv("gofile_api_key"),
		BaseUrl:    baseUrl,
		RetryCount: 3,
		Timeout:    10 * time.Second,
	}
}

// NewClient creates a new Client with the provided configuration
// It initializes an HTTP client with the specified timeout
func NewClient(c ClientConfig) *Client {
	return &Client{
		config: c,
		httpClient: &http.Client{
			Timeout: c.Timeout,
		},
	}
}

// setAuthorizationHeader adds a bearer token to the request's Authorization header
func setAuthorizationHeader(r *http.Request, t string) {
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t))
}

func getUploadServerURL(server string) string {
	return fmt.Sprintf("https://%s.gofile.io/contents/uploadfile", server)
}

// GetAvailableServers retrieves available servers, optionally filtered by zone
// If a zone is provided, it's added as a query parameter
// Returns the HTTP response or an error
func (c *Client) GetAvailableServers(zone string) (*http.Response, error) {
	u, err := url.Parse(c.config.BaseUrl + "/servers")
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
	setAuthorizationHeader(req, c.config.APIToken)
	return c.httpClient.Do(req)
}

// CreateFolder creates a folder in a folder with the speciifed parentFolderId
// If name is not specified, a name is auto-generated
// Returns the HTTP response or an error
func (c *Client) CreateFolder(parentFolderID string, name string) (*http.Response, error) {
	u := c.config.BaseUrl + "/contents/createFolder"

	payload := model.NewFolderPayload(parentFolderID, name)
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(postMethod, u, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	setAuthorizationHeader(req, c.config.APIToken)
	req.Header.Set("Content-Type", "application/json")
	return c.httpClient.Do(req)
}

// GetAccountId  gets the user ID
// Returns the HTTP response or an error
func (c *Client) GetAccountId() (*http.Response, error) {
	u := c.config.BaseUrl + "/accounts/getid"

	req, err := http.NewRequest(getMethod, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.config.APIToken)
	return c.httpClient.Do(req)
}

// GetAccountInformation gets the account information of the specifed user ID
// Returns the HTTP response or an error
func (c *Client) GetAccountInformation(id string) (*http.Response, error) {
	u := c.config.BaseUrl + fmt.Sprintf("/accounts/%s", id)

	req, err := http.NewRequest(getMethod, u, nil)
	if err != nil {
		return nil, err
	}
	setAuthorizationHeader(req, c.config.APIToken)
	return c.httpClient.Do(req)
}

// UploadFile uploads a file to a specified folder.
// If folderID is empty, a new public folder is created automatically.
// The base URL for the client changes to `https://{server}.gofile.io`
// Returns the HTTP response or an error
func (c *Client) UploadFile(server string, filePath string, folderID string) (*http.Response, error) {
	u := getUploadServerURL(server)
	var ct string // gets the content type from upload function
	pr := file.Upload(filePath, folderID, &ct)

	req, err := http.NewRequest(postMethod, u, pr)
	if err != nil {
		return nil, err
	}
	setAuthorizationHeader(req, c.config.APIToken)
	req.Header.Set("Content-Type", ct)
	return c.httpClient.Do(req)
	// i still need to change the timeout here
}
