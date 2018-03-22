// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"time"

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
		//check calendar every 10 minutes
		//queue sound files to play based on calendar events
		//set theme/files to play when doorbell events are processed
		//process doorbell events via MQTT
		newCalendarWatcher()
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

func newCalendarWatcher() {
	go func() {
		fmt.Printf("API Token from Config: %s\n", config.APIToken())
		if len(config.APIToken()) == 0 {
			fmt.Println("No API Token configured. Please run 'bellman configure'")
		}
		fmt.Println("Waiting for the done signal")
	}()
}
