package cmd

import (
	"Binalyze-OfflineUnzip/pkg"
	"Binalyze-OfflineUnzip/pkg/ui"
	"Binalyze-OfflineUnzip/pkg/validation"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"time"
)

// unzipCmd represents the unzip command
var unzipCmd = &cobra.Command{
	Use:   "unzip --key BINALYZE-LICENSE-KEY",
	Short: "Finds and Extracts Binalyze Offline Collector ZIPs",
	Long: `This will find and extract all Binalyze Offline Collector ZIP files. 
By default, this will look for ZIPs within the current folder, and output the results to the folder "output". 
These options are configurable. For example:

unzip --key BINALYZE-LICENSE-KEY --input customzipfolder --output customextractfolder

unzip --key BINALYZE-LICENSE-KEY --input customzipfolder`,
	Run: func(cmd *cobra.Command, args []string) {
		input := cmd.Flags().Lookup("input").Value.String()
		output := cmd.Flags().Lookup("output").Value.String()
		binLic := cmd.Flags().Lookup("key").Value.String()
		binEncPass := cmd.Flags().Lookup("password").Value.String()

		fmt.Println("Extracting Binalyze Offline Collector ZIPs found within " + input + " to " + output)

		//Get count of files in folder to attempt to extract
		files, err := os.ReadDir(input)
		if err != nil {
			panic(err)
		}

		filesCount := len(files)

		//Build UI Progress bar based on fileCount
		progressTracker := ui.CreateProgress()

		ui.InitiateProgress(progressTracker)

		for _, f := range files {
			if strings.HasSuffix(input+f.Name(), ".zip") {
				uid := local.GetZipUID(input + f.Name())
				pass := local.GenerateZipPass(uid, binLic, binEncPass)

				//Test ZIP pass - run function if zip pass is good, error otherwise
				tracker := ui.CreateTracker("Unzipping: "+f.Name(), int64(filesCount), progressTracker)
				testPassError := validation.TestZipPass(input+f.Name(), pass)
				if testPassError == nil {
					local.UnzipFile(input+f.Name(), pass, output, tracker)
				} else {
					tracker.UpdateMessage("Failed: " + f.Name())
					tracker.MarkAsErrored()
					//fmt.Printf("Error: %s\n", testPassError)
					//fmt.Println("Double check the specified Binalyze License Key or Encrypt Evidence Password is correct...\n\n ")
					//Failed for some reason - stop processing this ZIP and skip to the next one.
				}

			}
		}
		time.Sleep(100 * time.Millisecond)
	},
}

func init() {
	rootCmd.AddCommand(unzipCmd)

}
