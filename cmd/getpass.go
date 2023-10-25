package cmd

import (
	"Binalyze-OfflineUnzip/pkg"
	"Binalyze-OfflineUnzip/pkg/ui"
	"Binalyze-OfflineUnzip/pkg/validation"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	"os"
	"strings"
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
		input := cmd.Flags().Lookup("input").Value.String()
		binLic := cmd.Flags().Lookup("key").Value.String()
		binEncPass := cmd.Flags().Lookup("password").Value.String()

		//Get total count of files in folder
		files, err := os.ReadDir(input)
		if err != nil {
			panic(err)
		}

		filesCount := len(files)

		//Build UI Progress bar based on fileCount
		progressTracker := ui.CreateProgress()
		tracker := ui.CreateTracker("Generating Passwords", int64(filesCount), progressTracker)
		ui.InitiateProgress(progressTracker)

		//Start table to write identified passwords to
		t := table.NewWriter()
		//Table Header Names
		t.AppendHeader(table.Row{"Name", "Password"})
		t.SetStyle(table.StyleBold)
		colorBOnW := text.Colors{text.BgWhite, text.FgBlack}
		// set colors
		t.SetColumnConfigs([]table.ColumnConfig{
			{Name: "Name", Colors: text.Colors{text.FgHiBlack}, ColorsHeader: colorBOnW},
			{Name: "Password", Colors: text.Colors{text.FgHiGreen}, ColorsHeader: colorBOnW, ColorsFooter: colorBOnW},
		})

		for _, f := range files {
			if strings.HasSuffix(input+f.Name(), ".zip") {
				uid := local.GetZipUID(input + f.Name())
				pass := local.GenerateZipPass(uid, binLic, binEncPass)

				//Test Zip Passwords and write results to table
				testPassError := validation.TestZipPass(input+f.Name(), pass)
				if testPassError == nil {
					t.AppendRow(table.Row{f.Name(), pass})
				} else {
				}

			}
			tracker.Increment(1)
			//bar.Increment()
		}
		tracker.MarkAsDone()
		ui.FinishProgress(progressTracker)
		fmt.Println(t.Render())
	},
}

func init() {
	rootCmd.AddCommand(getpassCmd)

}
