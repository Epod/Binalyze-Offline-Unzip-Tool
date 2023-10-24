package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "bin_unzip",
	Short: "Decrypt Binalyze Offline Collector ZIPs without the need to connect to the internet.",
	//TODO: Replace this is relevant text
	Long: `This application allows the interaction with Binalyze Offline Collector ZIPs in air-gaped environments.

This will work entirely offline but will require manually entering information needed to generate ZIP passwords.

Typically all that is needed is the Binalyze License Key 
as this is what the Offline Collector uses when password encrypting ZIP files. 

By default, the program will look for ZIPs within the current folder the program is running from.

Examples:

bin_unzip unzip --key BINALYZE-LICENSE-KEY

bin_unzip getpass --key BINALYZE-LICENSE-KEY`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	//Disables the generate completion function in the Cobra lib that is turned on by default
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.PersistentFlags().StringP("key", "k", "",
		"The license key for the Binalyze instance which generated the Offline-Collector. (Required)")
	rootCmd.PersistentFlags().StringP("password", "p", "",
		"If the Offline Collector was generated with the \"Encrypt Evidence\" setting, provide that here.")
	rootCmd.PersistentFlags().StringP("input", "i", "./",
		"Path to folder containing zips. Defaults to scanning the same directory the program is running from")
	rootCmd.PersistentFlags().StringP("output", "o", "output",
		"Folder name or full path to write results to. Defaults to \"output\" in current directory")

	cobra.OnInitialize(initConfig)

	//Mark required only after Yiper config has had the chance to see if the config file loaded the binalize license key
	rootCmd.MarkPersistentFlagRequired("key")

}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {

		viper.AddConfigPath("$HOME")
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("bin_unzip_config")

	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())

		//Only set the Bin License Key if there is none manually specified in the program flags entered by the user
		if rootCmd.Flags().Lookup("key").Value.String() == "" {
			rootCmd.PersistentFlags().Set("key", viper.GetString("BINALYZE_LICENSE_KEY"))

		}
	}
}
