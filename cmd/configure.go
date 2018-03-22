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
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/beaubrewer/bellmanv2/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	calendar "google.golang.org/api/calendar/v3"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure Bellman",
	Long: `The configure command is used to configure Bellman
---------------------------
	Configuration
---------------------------
Bellman requires the following to be configured:
• Authorization to access your Google Calendar (generates api-key)
• Your Calendar ID to set which calendar Bellman uses
• MQTT queue name and username/password

HINT: Your Google Calendar ID can be found via the following steps
1. Go to http://calendar.google.com/ and login to your Google Calendar account.
2. Find and click the Google Calendar Settings.
3. Click the dropdown next to the calendar you would like to use and click 'calendar settings'
4. Your Calendar ID is shown near the bottom (e.g. 1sfvl6ubvv51e4qj67v2brqusk@group.calendar.google.com)
`,
	Run: func(cmd *cobra.Command, args []string) {
		b, err := ioutil.ReadFile("config/google-client.json")
		if err != nil {
			log.Fatalf("Unable to read google-client.json file: %v", err)
		}
		oauthConfig, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
		if err != nil {
			log.Fatalf("Unable to parse client secret file to googleConfig: %v", err)
		}

		tok := getTokenFromWeb(oauthConfig)
		viper.Set("api-key", tok)
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter the Calendar ID (this is found under the calendar settings): ")
		calendarID, _ := reader.ReadString('\n')
		calendarID = strings.TrimSuffix(calendarID, "\n")
		viper.Set("calendar_id", calendarID)
		config.SaveConfig()
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configureCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configureCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// GetTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	fmt.Print("Authorization Code: ")
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}
