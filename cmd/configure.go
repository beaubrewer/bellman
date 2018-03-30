// Copyright © 2018 Beau Brewer <beaubrewer@gmail.com>

package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	calendar "google.golang.org/api/calendar/v3"
)

// configureCmd is used to initialize Bellman's configuration
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
	RunE: func(cmd *cobra.Command, args []string) error {
		// check to make sure there are audio files available
		audioFileNames, err := filepath.Glob("audio/*.mp3")
		if err != nil || len(audioFileNames) == 0 {
			log.Fatalf("Please add MP3 files to Bellmans 'audio' directory before configuring. %v", err)
		}
		// configure google authentication
		b, err := ioutil.ReadFile("config/google-client.json")
		if err != nil {
			log.Fatalf("Unable to read google-client.json file: %v", err)
		}
		oauthConfig, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
		if err != nil {
			log.Fatalf("Unable to parse google-client.json file: %v", err)
		}
		// set the api-token retrieved from google
		tok := getTokenFromWeb(oauthConfig)
		viper.Set("api-token", tok)
		reader := bufio.NewReader(os.Stdin)

		// set the calendar ID
		fmt.Print("Enter the Calendar ID (this is found under the calendar settings): ")
		calendarID, _ := reader.ReadString('\n')
		calendarID = strings.TrimSuffix(calendarID, "\n")
		viper.Set("calendar_id", calendarID)

		// set Bellman defaults
		i := 1
		var loc int = 0
		var inputError error
		defaultAudio := make(map[string]string)
		for {
			fmt.Printf("\n\n---------------\n   Doorbell Setup    \n---------------\n")
			fmt.Printf("Enter a name for door #%d (e.g. front) or 'q' to quit: ", i)
			door, _ := reader.ReadString('\n')
			door = strings.TrimSuffix(door, "\n")
			if door == "q" {
				break
			}
			for pos, name := range audioFileNames {
				fmt.Printf("(%d) %s\n", pos, name)
			}
			fmt.Println("------------------------------")
			fmt.Printf("Enter the file number you want as the default for the %s door: ", door)
			_, inputError = fmt.Scanf("%d", &loc)
			for inputError != nil || len(audioFileNames) < loc {
				fmt.Printf("\nInvalid number. Try again: ")
				_, inputError = fmt.Scanf("%d", &loc)
			}
			doorAudio := audioFileNames[loc]
			defaultAudio[door] = strings.TrimRight(filepath.Base(doorAudio), ".mp3")
			i = i + 1
		}
		viper.Set("default_audio", defaultAudio)

		// set MQTT Host
		fmt.Print("Enter the MQTT Host:port to subscribe to (topic is bellman/doorbell): ")
		mqttHost, _ := reader.ReadString('\n')
		mqttHost = strings.TrimSuffix(mqttHost, "\n")
		viper.Set("mqtt_host", mqttHost)

		fmt.Println("--------------------------------")
		fmt.Println("     Configuration Complete!")
		fmt.Println("--------------------------------")
		fmt.Printf("Run `bellman server` to start\n\n")
		if err := viper.WriteConfigAs("config/config.json"); err != nil {
			return fmt.Errorf("unable to write the config to disk - %s", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)
}

// GetTokenFromWeb uses oauth2 to request an return a Token.
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
