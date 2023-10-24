package binzip

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/yeka/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	BinalyzeZipSecret = "binalyze.com/irec"
)

func GenerateZipPass(uid string, binLic string, binEncPass string) string {
	s := []byte(uid + binLic + binEncPass + BinalyzeZipSecret)
	h := sha256.New()
	h.Write([]byte(s))
	ZipHash := hex.EncodeToString(h.Sum(nil))
	return ZipHash
}

func GetZipUID(zipFile string) string {

	archive, err := zip.OpenReader(zipFile)
	if err != nil {
		panic(err)
	}

	for _, f := range archive.File {
		if strings.HasSuffix(f.Name, "Case.ppc") {
			comment := gjson.Get(f.Comment, "uid")
			return comment.Str
		}
	}
	fmt.Println("Could not find \"Case.ppc\" within ZIP.")
	return ""
}

func UnzipFile(zipFile string, zipPassword string, OutputFolder string) {
	archive, err := zip.OpenReader(zipFile)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(OutputFolder, f.Name)

		//File is encrypted
		if f.IsEncrypted() {
			f.SetPassword(zipPassword)
		}

		fmt.Println("unzipping ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(OutputFolder)+string(os.PathSeparator)) {
			fmt.Println("Error: Invalid output path")
		}
		if f.FileInfo().IsDir() {
			fmt.Println("Creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			fmt.Println("Error making directory to export to: " + err.Error())
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			fmt.Println("Error accessing folder: " + err.Error())
		}

		fileInArchive, err := f.Open()
		if err != nil {
			fmt.Println("Error opening archive: " + err.Error())
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			fmt.Println("Error extracting file: " + err.Error())
		}

		dstFile.Close()
		fileInArchive.Close()
	}
}
