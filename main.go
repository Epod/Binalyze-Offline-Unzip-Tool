package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/yeka/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	OutputFolder      = "output"
	BinalyzeZipSecret = "binalyze.com/irec"
)

func main() {

	//Get info from user needed to decrypt zips
	binLic := flag.String("license", "",
		"The license key for the Binalyze instance which generated the Offline-Collector")

	flag.Parse()

	//Check for required flags - fail if missing
	if *binLic == "" {
		fmt.Println("Missing --license flag.\n" +
			"Please manually enter your Binalyze license: ")
		fmt.Scanln(binLic)
	}

	//Loop Through All Zip Files
	zips, err := os.ReadDir(".")
	if err != nil {
		panic(err)
	}

	for _, f := range zips {
		if strings.HasSuffix(f.Name(), ".zip") {
			uid := GetZipUID(f.Name())
			pass := GenerateZipPass(uid, *binLic)
			UnzipFile(f.Name(), pass)
		}
	}

}

func GenerateZipPass(uid string, binLic string) string {
	s := []byte(uid + binLic + BinalyzeZipSecret)
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
			break
		}
	}
	fmt.Println("Could not find a UID...")
	return "unknown"
}

func UnzipFile(zipFile string, zipPassword string) {
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
			fmt.Println("invalid file path")
			return
		}
		if f.FileInfo().IsDir() {
			fmt.Println("creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}
}
