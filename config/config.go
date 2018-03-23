package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/oauth2"
)

type config struct {
	CalendarID string       `json:"calendar_id,omitempty"`
	APIToken   oauth2.Token `json:"api-token,omitempty"`
}

var c config

//Load and parses the config.json file
func Load() error {
	configFile, err := os.Open("config/config.json")
	if err != nil {
		return err
	}
	defer configFile.Close()
	configBytes, _ := ioutil.ReadAll(configFile)
	fmt.Printf("Loaded: %s\n", string(configBytes))
	err = json.Unmarshal(configBytes, &c)
	return nil
}

//GetAPIToken returns the oauth2 Token if available
func GetAPIToken() *oauth2.Token {
	return &c.APIToken
}

//GetCalendarID returns the configured Google Calendar ID
func GetCalendarID() string {
	return c.CalendarID
}
