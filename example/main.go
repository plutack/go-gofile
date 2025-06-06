// package main just illustrate a typical workflow  of using this package
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/plutack/go-gofile/api"
)

// What will be carried out in this example
// 1. Get account id
// 2. Get root folder id
// 3. create a folder called `test folder`
// 4. upload two files named ' `testfile1`, `testfile2` to EU servers
// 5. rename `testfile1` to `testfile1_renamed`
// 6. delete `testfile2`
// 7. rename `test folder` to `test folder renamed`

// CreateTextFile creates a text file at the specified location with the given name and content.
func CreateTextFile(filename, content, location string) error {
	fullPath := filepath.Join(location, filename)

	err := os.MkdirAll(location, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	fmt.Println("File created successfully:", fullPath)
	return nil
}

func progressCallbackfunc(done int64, total int64) {
	percentageCompleted := float64(done) / float64(total) * 100
	fmt.Printf("%.2f%% completed\n", percentageCompleted)
}

func main() {
	location, err := os.UserHomeDir()

	if err != nil {
		panic(err)
	}
	c := api.New(nil)
	s, err := c.GetAvailableServers("eu")
	if err != nil {
		panic(err)
	}
	// pick one of the random eu server returned
	euServer := s.Data.Servers[0].Name //this will be used to upload files

	accIdresp, err := c.GetAccountID() // this has the  account id nested in it
	if err != nil {
		panic(err)
	}
	accInfoResp, err := c.GetAccountInformation(accIdresp.Data.ID) // accountId gotten
	if err != nil {
		panic(err)
	}

	rootFolderId := accInfoResp.Data.RootFolder
	log.Printf("root folder id: %s\n", rootFolderId)

	folderInfoResp, err := c.CreateFolder(rootFolderId, "test folder")
	if err != nil {
		panic(err)
	}
	folderId := folderInfoResp.Data.ID
	log.Printf("test folder id: %s\n", folderId)

	CreateTextFile("testfile1", "hello world", location)
	CreateTextFile("testfile2", "hello world again", location)

	uploadFileResp1, err := c.UploadFile(euServer, filepath.Join(location, "testfile1"), rootFolderId, progressCallbackfunc)
	if err != nil {
		panic(err)
	}
	log.Printf("--------------\ntest folder 1 info\nname: %s\nID: %s\n--------------\n", uploadFileResp1.Data.Name, uploadFileResp1.Data.ID)
	uploadFileResp2, err := c.UploadFile(euServer, filepath.Join(location, "testfile2"), rootFolderId, progressCallbackfunc)
	if err != nil {
		panic(err)
	}
	log.Printf("--------------\ntest folder 2 info\nname: %s\nID: %s\n--------------\n", uploadFileResp2.Data.Name, uploadFileResp2.Data.ID)
	_, err = c.UpdateContent(uploadFileResp1.Data.ID, "name", "testfile1_renamed")
	if err != nil {
		panic(err)
	}
	_, err = c.UpdateContent(folderId, "name", "testfolder renamed")
	if err != nil {
		panic(err)
	}
	log.Println("completed")
}
