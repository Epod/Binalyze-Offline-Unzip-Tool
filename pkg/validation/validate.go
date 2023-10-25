package validation

import (
	"errors"
	"github.com/yeka/zip"
)

func TestZipPass(zipFile string, zipPass string) error {
	archive, err := zip.OpenReader(zipFile)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		f.SetPassword(zipPass)
		fileInArchive, err := f.Open()
		if err != nil {
			return errors.New("Password Incorrect For: " + f.Name)
		}
		fileInArchive.Close()
		//Stop the loop if we have a successful password test
		return nil
	}

	return nil
}
