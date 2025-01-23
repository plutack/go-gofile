package model

type newFolder struct {
	ParentFolderId string
	FolderName     string
}

func NewFolderPayload(p string, f string) newFolder {
	return newFolder{
		ParentFolderId: p,
		FolderName:     f,
	}
}
