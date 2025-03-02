// package model defines data structures for the API requests and responses
package model

import (
	"encoding/json"
	"strings"
	"time"
)

// deleteContent represents the payload to delete files or folders
type deleteContent struct {
	ContentsID string `json:"contentsId"` // array of ID of contents to be deleted
}

// newFolder represents the payload for creating a new folder.
//
// It contains the ID of the parent folder and the name of the new folder.
type newFolder struct {
	ParentFolderId string `json:"parentFolderId"` // ID of parent folder where a folder will be created
	FolderName     string `json:"folderName"`     // Name of the new folder
}

// updateContent represents the payload to modify a file attribute
type updateContent struct {
	Attribute      string          `json:"attribute"`
	AttributeValue json.RawMessage `json:"attributeValue"`
}

func (u *updateContent) WithName(n string) error {
	u.Attribute = "name"
	value, err := json.Marshal(n)
	if err != nil {
		return err
	}
	u.AttributeValue = value
	return nil
}
func (u *updateContent) WithTags(t string) error {
	u.Attribute = "tags"
	value, err := json.Marshal(t)
	if err != nil {
		return err
	}
	u.AttributeValue = value
	return nil
}

func (u *updateContent) WithDescription(d string) error {
	u.Attribute = "description"
	value, err := json.Marshal(d)
	if err != nil {
		return err
	}
	u.AttributeValue = value
	return nil
}

func (u *updateContent) WithPublic(p bool) error {
	u.Attribute = "public"
	value, err := json.Marshal(p)
	if err != nil {
		return err
	}
	u.AttributeValue = value
	return nil
}
func isValidTimeString(timeStr string) (int64, error) {
	// Try to parse the time string
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

func (u *updateContent) WithExpiry(t string) error {
	u.Attribute = "expiry"
	unixTime, err := isValidTimeString(t)
	if err != nil {
		return err
	}
	value, err := json.Marshal(unixTime)
	if err != nil {
		return err
	}
	u.AttributeValue = value
	return nil
}

func (u *updateContent) WithPassword(p string) error {
	u.Attribute = "password"
	value, err := json.Marshal(p)
	if err != nil {
		return err
	}
	u.AttributeValue = value
	return nil
}

// DeleteContentPayload creates an instance of deleteContent
//
// Returns delete content
func DeleteContentPayload(IDs []string) deleteContent {
	if len(IDs) == 0 {
		panic("at least one ID must be provided")
	}
	s := strings.Join(IDs, ",")
	return deleteContent{
		ContentsID: s,
	}
}

// NewFolderPayload creates an instance of newFolder
//
// Returns newFolder
func NewFolderPayload(p string, f string) newFolder {
	return newFolder{
		ParentFolderId: p,
		FolderName:     f,
	}
}

// NewUpdateContentPayload creates an instance of updateContent
//
// # Attribute can be any of the following with thei expected type for the new value
//
// 1. name = string
//
// 2. type description = string
//
// 3. type tags = []string
//
// 4. type public = bool
//
// 5. type expiry = string
//
// 6. type password = string
//
// Returns updateContent or an error
// deprecated
// func NewUpdateContentPayload(a string, v interface{}) (updateContent, error) {
// 	var marshaledValue []byte
// 	var err error
// 	if a == "tags" {
// 		array, ok := v.([]string)
// 		if !ok {
// 			return updateContent{}, nil
// 		}
// 		tagString := strings.Join(array, ",")
// 		marshaledValue, err = json.Marshal(tagString)
// 	} else {
// 		marshaledValue, err = json.Marshal(v)
// 	}

// 	if err != nil {
// 		return updateContent{}, err
// 	}
// 	return updateContent{
// 		Attribute:      a,
// 		AttributeValue: json.RawMessage(marshaledValue),
// 	}, nil
// }

// Returns updateContent or an error
func NewUpdateContentPayload() *updateContent {
	return &updateContent{}
}
