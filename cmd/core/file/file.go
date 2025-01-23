// Package file contains implementation related to unit file
package file

import (
	"io"
	"mime/multipart"
	"os"
)

// Upload creates a multipart/form-data request body for uploading a file.
// Returns a PipeReader that streams the data.
func Upload(filePath string, folderId string, contentType *string) *io.PipeReader {
	pr, pw := io.Pipe()
	w := multipart.NewWriter(pw)
	go func() {
		err := w.WriteField("folderId", folderId)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		f, err := os.Open(filePath)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		defer f.Close()

		part, err := w.CreateFormFile("file", f.Name())
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		_, err = io.Copy(part, f)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		pw.CloseWithError(w.Close())
	}()
	*contentType = w.FormDataContentType()
	return pr
}
