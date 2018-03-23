// Copyright © 2018 Beau Brewer <beaubrewer@gmail.com>

package cmd

import (
	"fmt"
	"time"

	bcal "github.com/beaubrewer/bellmanv2/calendar"
	"github.com/beaubrewer/bellmanv2/config"
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
			fmt.Println("Bellman is not configured. Please run 'bellman configure'")
			return
		}
		c := bcal.NewBellmanCalendar()
		c.ListEvents()
		//check calendar every 10 minutes
		//queue sound files to play based on calendar events
		//set theme/files to play when doorbell events are processed
		//process doorbell events via MQTT
		fmt.Println("Created the calendar watcher... now waiting 5 seconds")
		<-time.After(time.Second * 5)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
