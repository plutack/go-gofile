package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/plutack/gofile-api-go-client/cmd/internal/client"
	"github.com/plutack/gofile-api-go-client/cmd/model"
)

type api struct {
	client *client.Client
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
