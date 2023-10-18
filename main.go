package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/yeka/zip"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	BinalyzeZipSecret = "binalyze.com/irec"
)

func main() {

	//Get info from user needed to decrypt zips
	binLic := flag.String("key", "",
		"The license key for the Binalyze instance which generated the Offline-Collector.\n"+
			"Alternatively, you can create a file named \"key\" in the same folder as this program containing this information.")
	binEncPass := flag.String("password", "",
		"If the Offline Collector was generated with the \"Encrypt Evidence\" setting, provide that here.")
	input := flag.String("input", "./",
		"Path to folder containing zips. Defaults to scanning the same directory the program is running from")
	output := flag.String("output", "output",
		"Folder name or full path to write results to. Defaults to \"output\" in current directory")
	onlyPassList := flag.Bool("passlist", false,
		"If specified, the program will only print out the passwords for the zips. "+
			"Useful if you need to extract the zips through other means")

	flag.Parse()

	//Check for required flags - fail if missing
	if *binLic == "" {

		//Check if a "key" file is present and use the content as the binLic key -
		//Makes things more "double click and forget" friendly
		e, _ := os.Executable()
		exePath, _ := os.ReadDir(path.Dir(e))
		for _, kf := range exePath {
			if kf.Name() == "key" {
				fmt.Println("Saved keyfile found. Using for Binalyze License Key needed for decryption")
				keyfile, err := os.Open(path.Dir(e) + "/key")
				if err != nil {
					log.Fatal(err)
				}
				defer keyfile.Close()

				//Scan for the first line only of the file. Ignore everything else
				scanner := bufio.NewScanner(keyfile)
				var line int
				for scanner.Scan() {
					if line == 0 {
						*binLic = scanner.Text()
						break
					}
				}
				if err := scanner.Err(); err != nil {
					log.Fatalln(err)
				}
				//We have everything we need now - skip the rest of this IF statement
				goto start
			}
		}

		//Could not find saved Key strings using any method
		fmt.Println("Missing --key flag.\n" +
			"Please manually enter your Binalyze license: ")
		fmt.Scanln(binLic)
	}

start:
	//Loop Through All Zip Files
	zips, err := os.ReadDir(*input)
	if err != nil {
		panic(err)
	}

	for _, f := range zips {
		if strings.HasSuffix(*input+f.Name(), ".zip") {
			uid := GetZipUID(*input + f.Name())
			pass := GenerateZipPass(uid, *binLic, *binEncPass)

			//Test Zip Password
			if TestZipPass(*input+f.Name(), pass) == true {
				//Extract Zips when onlypasslist mode is disabled (the default)
				if *onlyPassList == false {
					UnzipFile(*input+f.Name(), pass, *output)
				} else {
					fmt.Println("\n" + "Container Name: " + f.Name())
					fmt.Println("Container Pass: " + pass + "\n--------")
				}
			} else {
				fmt.Println("Error With Archive: " + f.Name())
				fmt.Println("ZIP Password Incorrect. Check that the provided key is correct." + "\n--------")
			}

		}
	}

}

func TestZipPass(zipFile string, zipPass string) bool {
	archive, err := zip.OpenReader(zipFile)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		f.SetPassword(zipPass)
		fileInArchive, err := f.Open()
		if err != nil {
			return false
		}
		fileInArchive.Close()
		//No need to continue the loop if the password worked
		return true
	}

	return false
}

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
			break
		}
	}
	fmt.Println("Could not find a UID...")
	return "unknown"
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
