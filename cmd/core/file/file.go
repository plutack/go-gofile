package file

import (
	"io"
	"mime/multipart"
	"os"
)

func Upload(w *multipart.Writer, filePath string, folderId string) error {
	var err error
	errs := make(chan error, 1)
	go func() {
		err = w.WriteField("folderId", folderId)
		if err != nil {
			errs <- err
			return
		}
		f, err := os.Open(filePath)
		if err != nil {
			errs <- err
			return
		}
		defer f.Close()

		part, err := w.CreateFormFile("file", f.Name())
		if err != nil {
			errs <- err
			return
		}
		_, err = io.Copy(part, f)
		if err != nil {
			errs <- err
			return
		}

	}()
	if err := <-errs; err != nil {
		return err
	}
	return nil
}
