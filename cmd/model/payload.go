package model

// newFolder represents the payload for creating a new folder.
// It contains the ID of the parent folder and the name of the new folder.
type newFolder struct {
	ParentFolderId string // ID of parent folder where a folder will be created
	FolderName     string // Name of the new folder
}

// NewFolderPayload creates an instance of newFolder
// Returns newFolder
func NewFolderPayload(p string, f string) newFolder {
	return newFolder{
		ParentFolderId: p,
		FolderName:     f,
	}
}
