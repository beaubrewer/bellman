package calendar

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/beaubrewer/bellmanv2/config"
	"github.com/beaubrewer/bellmanv2/sound"
	"golang.org/x/oauth2/google"
	gcalendar "google.golang.org/api/calendar/v3"
)

type BellmanEvent struct {
	Type   string              `yaml:"type"`
	Chimes map[string][]string `yaml:"chimes"`
}

type BellmanCalendar struct {
	Chimes string
	srv    *gcalendar.Service
}

func NewBellmanCalendar() *BellmanCalendar {
	b, err := ioutil.ReadFile("config/google-client.json")
	if err != nil {
		log.Fatalf("Unable to read google-client.json file: %v", err)
	}
	oauthConfig, err := google.ConfigFromJSON(b, gcalendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse google-client.json file: %v", err)
	}
	client := oauthConfig.Client(context.Background(), config.GetAPIToken())
	service, err := gcalendar.New(client)
	if err != nil {
		log.Fatalf("Unable to create a google calendar client")
	}
	return &BellmanCalendar{srv: service}
}

func (cal *BellmanCalendar) ListEvents() {
	t := time.Now().Format(time.RFC3339)
	tmax := time.Now().Add(24 * 365 * time.Hour).Format(time.RFC3339)
	events, err := cal.srv.Events.List(config.GetCalendarID()).ShowDeleted(false).
		SingleEvents(true).TimeMin(t).TimeMax(tmax).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events. %v", err)
	}

	fmt.Println("Upcoming events:")
	if len(events.Items) > 0 {
		for _, i := range events.Items {
			var when string
			// If the DateTime is an empty string the Event is an all-day Event.
			// So only Date is available.
			if i.Start.DateTime != "" {
				when = i.Start.DateTime
			} else {
				when = i.Start.Date
			}
			fmt.Printf("%s (%s)\n", i.Summary, when)
			//trying to parse description as yaml
			var yevent BellmanEvent
			yaml.Unmarshal([]byte(i.Description), &yevent)
			fmt.Printf("\nHere's the type: %s\n", yevent.Type)
			for k := range yevent.Chimes {
				fmt.Printf("%s Door\n", k)
				for _, z := range yevent.Chimes[k] {
					fmt.Println(z)
					sound.Play(z)
				}
			}
		}
	} else {
		fmt.Printf("No upcoming events found.\n")
	}
}
