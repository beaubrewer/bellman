// Copyright © 2018 Beau Brewer <beaubrewer@gmail.com>

package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

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
	RunE: func(cmd *cobra.Command, args []string) error {
		b, err := ioutil.ReadFile("config/google-client.json")
		if err != nil {
			log.Fatalf("Unable to read google-client.json file: %v", err)
		}
		oauthConfig, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
		if err != nil {
			log.Fatalf("Unable to parse google-client.json file: %v", err)
		}

		tok := getTokenFromWeb(oauthConfig)
		viper.Set("api-token", tok)
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter the Calendar ID (this is found under the calendar settings): ")
		calendarID, _ := reader.ReadString('\n')
		calendarID = strings.TrimSuffix(calendarID, "\n")
		viper.Set("calendar_id", calendarID)
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
