package config

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type config struct {
	CalendarID string       `json:"calendar_id,omitempty" mapstructure:"calendar_id"`
	APIToken   oauth2.Token `json:"api-token,omitempty" mapstructure:"api-token"`
}

//APIToken returns the token if available
func APIToken() string {
	var c config
	fmt.Println("Getting api token if it's available")
	if err := viper.Unmarshal(&c); err != nil {
		return ""
	}
	if c.APIToken.AccessToken == "" {
		return ""
	}
	s, err := json.Marshal(c.APIToken)
	if err != nil {
		return ""
	}
	return string(s)
}

//SaveConfig will write the config file back to disk
func SaveConfig() error {
	fmt.Println("--------------------------------")
	fmt.Println("     Configuration Complete!")
	fmt.Println("--------------------------------")
	fmt.Println("Run `bellman server` to start")
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("unable to write the config to disk - %s", err)
	}
	return nil
}
