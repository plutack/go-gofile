package model

// AvailableServerResponse represents the response structure for available servers.
// Contains status and data about servers in all zones.
type AvailableServerResponse struct {
	Status string `json:"servers"`
	Data   struct {
		Servers        []server `json:"servers"`        // servers in the specified zone
		ServersAllZone []server `json:"serversAllZone"` // servers across all zones
	} `json:"data"`
}

// server represents a server with its name and zone
type server struct {
	Name string `json:"name"` // name of the server
	Zone string `json:"zone"` // zone where the server is located
}
