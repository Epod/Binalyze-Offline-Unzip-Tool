package cmd

import (
	binzip "Binalyze-OfflineUnzip/pkg"
	"Binalyze-OfflineUnzip/pkg/validation"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// getpassCmd represents the getpass command
var getpassCmd = &cobra.Command{
	Use:   "getpass --key BINALYZE-LICENSE-KEY",
	Short: "Only export the list of Binalyze ZIP Passwords",
	Long: `This will only export the passwords for the ZIPs provided. 
This is useful if you are looking to unzip the Binalyze Offline Collector ZIPs using another method outside this program. 
For example:

bin_unzip getpass --key BINALYZE-LICENSE-KEY

bin_unzip getpass --key BINALYZE-LICENSE-KEY --input zips --output extracted`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Generating Binalyze Offline Collector ZIP Passwords")

		input := cmd.Flags().Lookup("input").Value.String()
		binLic := cmd.Flags().Lookup("key").Value.String()
		binEncPass := cmd.Flags().Lookup("password").Value.String()

		zips, err := os.ReadDir(input)
		if err != nil {
			panic(err)
		}

		for _, f := range zips {
			if strings.HasSuffix(input+f.Name(), ".zip") {
				uid := binzip.GetZipUID(input + f.Name())
				pass := binzip.GenerateZipPass(uid, binLic, binEncPass)

				//Test Zip Password
				testPassError := validation.TestZipPass(input+f.Name(), pass)
				if testPassError == nil {
					fmt.Println("\n" + "Container Name: " + f.Name())
					fmt.Println("Container Pass: " + pass + "\n--------")
				} else {
					fmt.Printf("Error: %s\n", testPassError)
					fmt.Println("Double check the specified Binalyze License Key or Encrypt Evidence Password is correct...\n\n ")
				}

			}
		}
	},
}

func init() {
	rootCmd.AddCommand(getpassCmd)

}
