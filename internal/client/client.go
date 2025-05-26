// package client contains low level abstractions to wrap around gofile.io api
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
	"strings"
	"time"

	"github.com/plutack/go-gofile/model"
)

// HTTP request methods for API interactions
const (
	getMethod    = "GET"
	postMethod   = "POST"
	putMethod    = "PUT"
	deleteMethod = "DELETE"
)

// baseUrl is the base  URL used for gofile.io api calls
var baseUrl = "https://api.gofile.io"

// clientConfig contains necessary configuration options to configure a client
type ClientConfig struct {
	APIToken   string        // APIToken is the authentication token for the GoFile.io API
	BaseUrl    string        //BaseUrl is the base url for API request apart from uploadFile API call
	RetryCount int           // RetryCount specifies the number of times to retry failed API requests
	Timeout    time.Duration // Timeout specifies the maximum time to wait for an API Request to be resolved
}

// ProgressCallback represents a function that receives progress updates.
// done is the number of bytes uploaded so far, and total is the total number of bytes.
type ProgressCallback = func(done int64, total int64)

// Client represents an HTTP client for interacting with the GoFile.io API
type Client struct {
	httpClient *http.Client // httpClient is the underlying HTTP client used for API requests
	config     ClientConfig // config holds the configuration settings for the API client
}

// progressReader wraps an io.Reader and reports progress as bytes are read.
// It tracks the total number of bytes read so far and invokes the onRead callback
// with the current progress and total size.
type progressReader struct {
	io.Reader
	total  int64
	size   int64
	onRead func(total int64, size int64)
}

// Read is a custom implementation that wraps the underlying reader's Read method
// and invokes the onRead callback to report progress as data is read.
func (p *progressReader) Read(buf []byte) (int, error) {
	n, err := p.Reader.Read(buf)
	if n > 0 {
		p.total += int64(n)
		p.onRead(p.total, p.size)
	}
	return n, err
}

// Returns the ratio of the tile that has been read in percentages
func (p *progressReader) PercentageCompleted() float64 {
	return (float64(p.total) / float64(p.size)) * 100
}

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
		Timeout:    1 * time.Minute,
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

// Upload creates a multipart/form-data request body for uploading a file.
// Returns a PipeReader that streams the data.
func upload(filePath string, folderId string, contentType *string, onProgress ProgressCallback) *io.PipeReader {
	pr, pw := io.Pipe()
	w := multipart.NewWriter(pw)
	go func() {
		err := w.WriteField("folderId", folderId)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		f, err := os.Open(filePath)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		defer f.Close()
		fi, err := f.Stat()
		progressR := &progressReader{
			Reader: f,
			size:   fi.Size(),
			total:  0,
			onRead: onProgress,
		}
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		part, err := w.CreateFormFile("file", fi.Name())
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		_, err = io.Copy(part, progressR)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		pw.CloseWithError(w.Close())
	}()
	*contentType = w.FormDataContentType()
	return pr
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

// DeleteContent deletes files and folder  the speciifed contentID(s)
// Returns the HTTP response or an error
func (c *Client) DeleteContent(IDs []string) (*http.Response, error) {
	u := c.config.BaseUrl + "/contents"

	payload := model.DeleteContentPayload(IDs)
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(deleteMethod, u, bytes.NewBuffer(data))
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
	u := fmt.Sprintf("%s/accounts/%s", c.config.BaseUrl, id)

	req, err := http.NewRequest(getMethod, u, nil)
	if err != nil {
		return nil, err
	}
	setAuthorizationHeader(req, c.config.APIToken)
	return c.httpClient.Do(req)
}

// UpdateContent changes the attribute of a file or folder.
// Returns the HTTP response or an error
func (c *Client) UpdateContent(contentID string, attribute string, value interface{}) (*http.Response, error) {
	u := fmt.Sprintf("%s/contents/%s/update", c.config.BaseUrl, contentID)

	payload := model.NewUpdateContentPayload()
	var err error

	switch attribute {
	case "name":
		nameStr, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("name must be string, got %T", value)
		}
		err = payload.WithName(nameStr)

	case "description":
		descStr, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("description must be string, got %T", value)
		}
		err = payload.WithDescription(descStr)

	case "tags":
		slice, ok := value.([]string)
		if !ok {
			return nil, fmt.Errorf("description must be string, got %T", value)
		}
		tagString := strings.Join(slice, ",")
		err = payload.WithTags(tagString)

	case "public":
		pubBool, ok := value.(bool)
		if !ok {
			return nil, fmt.Errorf("public must be boolean, got %T", value)
		}
		err = payload.WithPublic(pubBool)

	case "expiry":
		expiryStr, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("expiry must be string in RFC3339 format, got %T", value)
		}
		err = payload.WithExpiry(expiryStr)

	case "password":
		passStr, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("password must be string, got %T", value)
		}
		err = payload.WithPassword(passStr)

	default:
		return nil, fmt.Errorf("unsupported attribute: %s", attribute)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to set attribute %s: %w", attribute, err)
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload failed: %w", err)
	}
	req, err := http.NewRequest(http.MethodPut, u, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	setAuthorizationHeader(req, c.config.APIToken)
	return c.httpClient.Do(req)
}

// UploadFile uploads a file to a specified folder.
// If folderID is empty, a new public folder is created automatically.
// The base URL for the client changes to `https://{server}.gofile.io`
// Returns the HTTP response or an error
func (c *Client) UploadFile(server string, filePath string, folderID string, callbackUpdate ProgressCallback) (*http.Response, error) {
	u := getUploadServerURL(server)
	var ct string // gets the content type from upload function
	pr := upload(filePath, folderID, &ct, callbackUpdate)
	c.httpClient.Timeout = 0
	req, err := http.NewRequest(postMethod, u, pr)
	if err != nil {
		return nil, err
	}
	setAuthorizationHeader(req, c.config.APIToken)
	req.Header.Set("Content-Type", ct)
	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	c.httpClient.Timeout = 1 * time.Minute
	return response, nil
}
