package model

// currentStats represents information about the user root folder
type currentStats struct {
	FolderCount int `json:"folderCount"` // number of folders in user's root folder
	FileCount   int `json:"fileCount"`   // number of files in user's root folder
	Storage     int `json:"storage"`     //
}

// folderInfo represents the information about a folder
type folderInfo struct {
	ID           string `json:"id"`           // ID of the folder
	Owner        string `json:"owner"`        // ID of the creator of the folder
	Type         string `json:"type"`         // this is always folder
	Name         string `json:"name"`         // name of the folder
	ParentFolder string `json:"parentFolder"` // ID of the parent folder
	CreateTime   string `json:"createTime"`   // date:time the folder was created
	ModTime      string `json:"modTime"`      //
	Code         string `json:"code"`         // short code of the folder?
} // server represents a server with its name and zone
type server struct {
	Name string `json:"name"` // name of the server
	Zone string `json:"zone"` // zone where the server is located
}

// fileInfo represents information about a file
type fileInfo struct {
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

// AccountIDResponse represents the response structure for the user's ID
// Contains status and the user ID
type AccountIDResponse struct {
	Status string `json:"status"`
	Data   struct {
		ID string `json:"id"` //ID of user account
	} `json:"data"`
}

// AccountInformationResponse represent the response structure for a user account information
// Contains status and data about the user account
type AccountInformationResponse struct {
	Status string `json:"status"`
	Data   struct {
		IPTraffic30  int          `json:"ipTraffic30"` //
		ID           string       `json:"id"`          // Id of user account
		CreateTime   int          `json:"createTime"`  // date:time account was created
		Email        string       `json:"email"`       // email address of user
		Tier         string       `json:"tier"`        // tier of user account
		Token        string       `json:"token"`       // bearer token for Authorization header
		RootFolder   string       `json:"rootFolder"`  // ID of user's root folder
		StatsCurrent currentStats `json:"statsCurrent"`
	} `json:"data"`
}

// AvailableServerResponse represents the response structure for available servers.
// Contains status and data about servers in all zones.
type AvailableServerResponse struct {
	Status string `json:"status"`
	Data   struct {
		Servers        []server `json:"servers"`        // servers in the specified zone
		ServersAllZone []server `json:"serversAllZone"` // servers across all zones
	} `json:"data"`
}

// CreateFolderResponse represents the response structure for a successful folder creation
// Contains status and data about the created folder
type CreateFolderResponse struct {
	Status string     `json:"status"`
	Data   folderInfo `json:"data"`
}

// UploadFileResponse represent the response structure for a successful file upload
// Contains status and data about uploaded file
type UploadFileResponse struct {
	Data   fileInfo `json:"data"`
	Status string   `json:"status"`
}
