package api

import (
	"io"
	"net/http"

	"github.com/plutack/go-files-api/cmd/internal/client"
)

type api struct {
	client *client.Client
}

func readResponseBody(r *http.Response) (string, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	return string(body), nil
}

func (a *api) getAvailableServers(zone string) (string, error) {
	resp, err := a.client.GetAvailableServers(zone)
	if err != nil {
		return "", err

	}
	body, err := readResponseBody(resp)
	if err != nil {
		return "", err
	}

	return body, nil

}
