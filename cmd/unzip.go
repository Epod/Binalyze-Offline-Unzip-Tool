/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	binzip "Binalyze-OfflineUnzip/pkg"
	"Binalyze-OfflineUnzip/pkg/validation"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
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

		zips, err := os.ReadDir(input)
		if err != nil {
			panic(err)
		}

		for _, f := range zips {
			if strings.HasSuffix(input+f.Name(), ".zip") {
				uid := binzip.GetZipUID(input + f.Name())
				pass := binzip.GenerateZipPass(uid, binLic, binEncPass)

				//Test ZIP pass - run function if zip pass is good, error otherwise
				testPassError := validation.TestZipPass(input+f.Name(), pass)
				if testPassError == nil {
					binzip.UnzipFile(input+f.Name(), pass, output)
				} else {
					fmt.Printf("Error: %s\n", testPassError)
					fmt.Println("Double check the specified Binalyze License Key or Encrypt Evidence Password is correct...\n\n ")
					//Failed for some reason - stop processing this ZIP and skip to the next one.
					continue
				}

			}
		}
	},
}

func init() {
	rootCmd.AddCommand(unzipCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// unzipCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// unzipCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
