package main

import (
	"github.com/plutack/gofile-api-go-client/cmd/pkg/api"
)

func main() {
	s := api.New(nil)

	_, err := s.GetAvailableServers("eu")
	if err != nil {
		panic(err)
	}
}
