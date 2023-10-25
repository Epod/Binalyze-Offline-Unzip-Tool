package local

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/progress"
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

func UnzipFile(zipFile string, zipPassword string, OutputFolder string, tracker *progress.Tracker) {
	archive, err := zip.OpenReader(zipFile)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	archiveFilesCount := len(archive.File)
	tracker.UpdateTotal(int64(archiveFilesCount))

	for _, f := range archive.File {
		filePath := filepath.Join(OutputFolder, f.Name)

		//File is encrypted
		if f.IsEncrypted() {
			f.SetPassword(zipPassword)
		}

		//fmt.Println("unzipping ", filePath)
		tracker.UpdateMessage("Extracting File: " + f.Name)

		if !strings.HasPrefix(filePath, filepath.Clean(OutputFolder)+string(os.PathSeparator)) {
			tracker.MarkAsErrored()
		}
		//Create directory if the detected file is a dir
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}
		//Failed to make directory
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			tracker.MarkAsErrored()
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		//Could not access destination folder
		if err != nil {
			tracker.MarkAsErrored()
		}

		fileInArchive, err := f.Open()
		//Could not open the zip file
		if err != nil {
			tracker.MarkAsErrored()
		}

		//Error when writing file from zip
		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			tracker.MarkAsErrored()
		}

		dstFile.Close()
		fileInArchive.Close()
		tracker.Increment(1)
	}

	tracker.UpdateMessage("Success: " + zipFile)
	tracker.MarkAsDone()
}
