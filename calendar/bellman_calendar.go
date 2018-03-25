package calendar

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/beaubrewer/bellmanv2/config"
	"golang.org/x/oauth2/google"
	gcalendar "google.golang.org/api/calendar/v3"
)

type BellmanEvent struct {
	Type  string              `yaml:"type"`
	Audio map[string][]string `yaml:"audio"`
	Start time.Time
	End   time.Time
}

type BellmanCalendar struct {
	srv *gcalendar.Service
}

func NewBellmanEvent(event *gcalendar.Event) *BellmanEvent {
	var e BellmanEvent
	if err := yaml.Unmarshal([]byte(event.Description), &e); err != nil {
		fmt.Printf("There was a problem unmarshalling the event\n%s\n", err)
	}
	// check to see if this is a full day event
	// full day events do not provide timezone information so we call ParseInLocation
	// and provide the timezone/location data
	if len(event.Start.Date) > 0 {
		z := time.Now().Location()
		e.Start, _ = time.ParseInLocation("2006-02-01", event.Start.Date, z)
		e.End, _ = time.ParseInLocation("2006-02-01", event.End.Date, z)
	} else {
		e.Start, _ = time.Parse(time.RFC3339, event.Start.DateTime)
		e.End, _ = time.Parse(time.RFC3339, event.End.DateTime)
	}
	fmt.Println(e.Start.Unix())
	return &e
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

func (cal *BellmanCalendar) GetEvents(t time.Duration) *gcalendar.Events {
	currentTime := time.Now()
	currentTimeRFC := currentTime.Format(time.RFC3339)
	endOfDay := currentTime.Add(t)
	endOfDayRFC := endOfDay.Format(time.RFC3339)
	events, err := cal.srv.Events.List(config.GetCalendarID()).ShowDeleted(false).
		SingleEvents(true).TimeMin(currentTimeRFC).TimeMax(endOfDayRFC).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve user's events. %v\n", err)
	}
	return events
}

func (cal *BellmanCalendar) GetCurrentTheme() map[string][]string {
	//Find the most recent theme event starting today and going back 1/yr
	tmax := time.Now().Format(time.RFC3339)
	t := time.Now().Add(-365 * 24 * time.Hour).Format(time.RFC3339)
	fmt.Printf("Searching for theme from %s to %s\n", t, tmax)
	events, err := cal.srv.Events.List(config.GetCalendarID()).ShowDeleted(false).
		SingleEvents(true).TimeMin(t).TimeMax(tmax).OrderBy("startTime").Do()
	if err != nil {
		log.Fatal("Unable to retrieve events. %v", err)
	}
	for i := len(events.Items) - 1; i >= 0; i-- {
		e := events.Items[i]
		be := NewBellmanEvent(e)
		fmt.Printf("Event found:\n%s\n", e.Summary)
		fmt.Printf("Start: %s\n", be.Start.Format(time.RFC3339))
		fmt.Printf("End: %s\n", be.End.Format(time.RFC3339))
		if strings.ToLower(be.Type) == "theme" {
			return be.Audio
		}
	}
	fmt.Println("No theme found")
	return map[string][]string{}
}
