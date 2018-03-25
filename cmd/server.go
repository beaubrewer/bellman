// Copyright © 2018 Beau Brewer <beaubrewer@gmail.com>

package cmd

import (
	"fmt"
	"time"

	"github.com/beaubrewer/bellmanv2/config"
	"github.com/beaubrewer/bellmanv2/manager"
	"github.com/spf13/cobra"
)

var done chan time.Time

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start Bellman",
	Long:  `The server command will start Bellman`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.Load(); err != nil {
			fmt.Printf("%s\nBellman is not configured. Please run 'bellman configure'", err)
			return
		}
		manager.Start()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
