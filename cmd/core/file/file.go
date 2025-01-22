package file

import (
	"io"
	"mime/multipart"
	"os"
)

func Upload(w *multipart.Writer, filePath string, folderId string) error {
	var err error

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	part, err := w.CreateFormFile("file", f.Name())
	if err != nil {
		return err
	}
	err = w.WriteField("folderId", folderId)
	if err != nil {
		return err
	}
	_, err = io.Copy(part, f)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}
