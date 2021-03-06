// Copyright © 2018 Beau Brewer <beaubrewer@gmail.com>

package cmd

import (
	"fmt"
	"time"

	"github.com/beaubrewer/bellman/config"
	"github.com/beaubrewer/bellman/manager"
	"github.com/spf13/cobra"
)

var done chan time.Time

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start Bellman",
	Long:  `The server command will start Bellman`,
	Run: func(cmd *cobra.Command, args []string) {
		manager.Start()
	},
}

func init() {
	if err := config.Load(); err != nil {
		fmt.Printf("%s\nBellman is not configured. Please run 'bellman configure'", err)
		return
	}
	rootCmd.AddCommand(serverCmd)
}
