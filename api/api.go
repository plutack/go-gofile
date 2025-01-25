// package api exposes is exposed for user's reuse or dependency sharing
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/plutack/go-gofile/internal/client"
	"github.com/plutack/go-gofile/model"
)

type Api struct {
	client *client.Client
}

// Options defines optional configuration for the API client.
type Options struct {
	APIToken   *string // APIToken is the authentication token for the GoFile.io API
	RetryCount *int    // RetryCount specifies the number of times to retry failed API requests
	Timeout    *int    // Timeout specifies the maximum time to wait for an API Request to be resolved

}

// New initializes a new API client with optional configuration.
//
// If opts is nil, default client settings are used.
func New(opts *Options) *Api {
	clientConfig := client.NewDefaultClientConfig()
	if opts == nil {
		apiClient := client.NewClient(clientConfig)
		return &Api{
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

	return &Api{
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
//
// zone can either be "eu" or "na"
// Returns a structured response or an error.
func (a *Api) GetAvailableServers(zone string) (model.AvailableServerResponse, error) {
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

// DeleteContent delete files and folders uploaded or created by user
//
// If contentID arguement is not supplied code panics.
//
// Returns a structured response or an error.
func (a *Api) DeleteContent(contentID ...string) (model.DeleteContentResponse, error) {
	resp, err := a.client.DeleteContent(contentID)
	if err != nil {
		return model.DeleteContentResponse{}, err

	}
	buf, err := readResponseBody(resp)
	if err != nil {
		return model.DeleteContentResponse{}, err
	}

	var body model.DeleteContentResponse
	json.Unmarshal(buf, &body)
	return body, nil
}

// UpdateContent changes the attribute of a file or a folder.
//
// # Attribute can be any of the following with their expected type for the new value
//
// 1. name = string
//
// 2. type description = string
//
// 3. type tags = []string
//
// 4. type public = bool
//
// 5. type expiry = string
//
// 6. type password = string
//
// Returns a structured response or an error.
func (a *Api) UpdateContent(contentID string, attribute string, newAttributeValue interface{}) (model.UpdateContentResponse, error) {
	resp, err := a.client.UpdateContent(contentID, attribute, newAttributeValue)
	if err != nil {
		return model.UpdateContentResponse{}, err
	}
	buf, err := readResponseBody(resp)
	if err != nil {
		return model.UpdateContentResponse{}, err
	}

	var body model.UpdateContentResponse
	json.Unmarshal(buf, &body)
	return body, nil
}

// UploadFile saves a file on a specified server
//
// Returns a structured response or an error.
func (a *Api) UploadFile(server string, filePath string, folderID string) (model.UploadFileResponse, error) {
	resp, err := a.client.UploadFile(server, filePath, folderID)
	if err != nil {
		return model.UploadFileResponse{}, err

	}
	buf, err := readResponseBody(resp)
	if err != nil {
		return model.UploadFileResponse{}, err
	}

	var body model.UploadFileResponse
	json.Unmarshal(buf, &body)
	return body, nil
}

// CreateFolder makes a new folder at the root of the specified parent folder id
//
//	See model.CreateFolderResponse for struct structure
//
// Returns a structured response or an error.
func (a *Api) CreateFolder(parentFolderID string, name string) (model.CreateFolderResponse, error) {
	resp, err := a.client.CreateFolder(parentFolderID, name)
	if err != nil {
		return model.CreateFolderResponse{}, err

	}
	buf, err := readResponseBody(resp)
	if err != nil {
		return model.CreateFolderResponse{}, err
	}

	var body model.CreateFolderResponse
	json.Unmarshal(buf, &body)
	return body, nil
}

// GetAccountID returns a struct containing the user account ID.
//
//	See model.AccountIDResponse for struct structure
//
// Returns a structured response or an error.
func (a *Api) GetAccountID() (model.AccountIDResponse, error) {
	resp, err := a.client.GetAccountId()
	if err != nil {
		return model.AccountIDResponse{}, err

	}
	buf, err := readResponseBody(resp)
	if err != nil {
		return model.AccountIDResponse{}, err
	}

	var body model.AccountIDResponse
	json.Unmarshal(buf, &body)
	return body, nil
}

// GetAccountInformation returns a struct containing the user account information.
//
// NOTE: this is where the root folder ID can be gotten from
//
//	See model.AccountInformationResponse for struct structure
//
// Returns a structured response or an error.
func (a *Api) GetAccountInformation(accountId string) (model.AccountInformationResponse, error) {
	resp, err := a.client.GetAccountInformation(accountId)
	if err != nil {
		return model.AccountInformationResponse{}, err

	}
	buf, err := readResponseBody(resp)
	if err != nil {
		return model.AccountInformationResponse{}, err
	}

	var body model.AccountInformationResponse
	json.Unmarshal(buf, &body)
	return body, nil
}

// features to be implemented
// func (a *api) ResetToken() {}
// premium features to  be implemented
// func (a *api) CreateDirectLink()       {}
// func (a *api) UpdateDirectLinkConfig() {}
// func (a *api) DeleteDirectLink()       {}
// func (a *api) CopyContent()            {}
// func (a *api) MoveContent()            {}
// func (a *api) ImportContent()          {}
