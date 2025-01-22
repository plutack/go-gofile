package model

type AvailableServerResponse struct {
	Status string `json:"servers"`
	Data   struct {
		Servers        []server `json:"servers"`
		ServersAllZone []server `json:"serversAllZone"`
	} `json:"data"`
}

type server struct {
	Name string `json:"name"`
	Zone string `json:"zone"`
}
