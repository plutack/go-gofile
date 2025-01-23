package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/plutack/gofile-api-go-client/cmd/internal/client"
	"github.com/plutack/gofile-api-go-client/cmd/model"
)

type api struct {
	client *client.Client
}

// Options defines optional configuration for the API client.
type Options struct {
	// APIToken is the authentication token for the GoFile.io API
	APIToken *string
	// RetryCount specifies the number of times to retry failed API requests
	RetryCount *int
	// Timeout specifies the maximum time to wait for an API Request to be resolved
	Timeout *int
}

// New initializes a new API client with optional configuration.
// If opts is nil, default client settings are used.
func New(opts *Options) *api {
	clientConfig := client.NewDefaultClientConfig()
	if opts == nil {
		apiClient := client.NewClient(clientConfig)
		return &api{
			client: apiClient,
		}

	}
	if opts.APIToken != nil {
		clientConfig.APIToken = *opts.APIToken
	}

	if opts.RetryCount != nil {
		clientConfig.RetryCount = *opts.RetryCount
	}

	if opts.Timeout != nil {
		clientConfig.Timeout = time.Duration(*opts.Timeout) * time.Second
	}

	apiClient := client.NewClient(clientConfig)

	return &api{
		client: apiClient,
	}
}

// readResponseBody reads and returns the response body as a byte slice.
func readResponseBody(r *http.Response) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	// Logs the response body for debugging just to pipe to jq
	fmt.Println(string(body))
	return body, nil
}

// GetAvailableServers retrieves available servers, optionally filtered by zone
// Returns a structured response or an error.
func (a *api) GetAvailableServers(zone string) (model.AvailableServerResponse, error) {
	resp, err := a.client.GetAvailableServers(zone)
	if err != nil {
		return model.AvailableServerResponse{}, err

	}
	buf, err := readResponseBody(resp)
	if err != nil {
		return model.AvailableServerResponse{}, err
	}

	var body model.AvailableServerResponse
	json.Unmarshal(buf, &body)
	return body, nil
}

// UploadFile to a specified server
// Returns a structured response or an error.
func (a *api) UploadFile(server string, filePath string, folderID string) (model.UploadFileResponse, error) {
	resp, err := a.client.UploadFile(server, filePath, folderID)
	if err != nil {
		return model.AvailableServerResponse{}, err

	}
	buf, err := readResponseBody(resp)
	if err != nil {
		return model.AvailableServerResponse{}, err
	}

	var body model.AvailableServerResponse
	json.Unmarshal(buf, &body)
	return body, nil
}
