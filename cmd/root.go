// Copyright © 2018 Beau Brewer <beaubrewer@gmail.com>

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bellman",
	Short: "Bellman - the ultimate doorbell application at your service",
	Long: `Bellman uses any Google Calendar to play doorbell chimes on your schedule.

Sound files for your doorbell chimes must be in MP3 format located in the /chimes directory.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
