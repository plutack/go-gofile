package api

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/plutack/gofile-api-go-client/cmd/internal/client"
	"github.com/plutack/gofile-api-go-client/cmd/model"
)

type api struct {
	client *client.Client
}

type Options struct {
	APIToken   *string
	RetryCount *int
	Timeout    *int
}

func New(opts Options) *api {
	clientConfig := client.NewDefaultClientConfig()

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
func readResponseBody(r *http.Response) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	return body, nil
}

func (a *api) getAvailableServers(zone string) (model.AvailableServerResponse, error) {
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
