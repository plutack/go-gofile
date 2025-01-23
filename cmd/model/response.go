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
type UploadFileResponse struct {
	Data   uploadFileInfo `json:"data"`
	Status string         `json:"status"`
}

type uploadFileInfo struct {
	CreateTime       int      `json:"createTime"`       // time the file was uploaded
	DownloadPage     string   `json:"downloadPage"`     // gofile.io download link page for the file
	ID               string   `json:"id"`               // ID of the file on the gofile server
	MD5              string   `json:"md5"`              // MD5 hash of the file
	Mimetype         string   `json:"mimetype"`         // type of the file (eg: "application/zip")
	ModTime          int      `json:"modTime"`          //
	Name             string   `json:"name"`             // name of the file uploaded
	ParentFolder     string   `json:"parentFolder"`     // ID of parent folder
	ParentFolderCode string   `json:"parentFolderCode"` //code of parent folder
	Servers          []string `json:"servers"`          // array of name of servers the uploaded file is on
	Size             int      `json:"size"`             // size of the file in bytes
	Type             string   `json:"type"`             // type of file (eg: "file")
}
