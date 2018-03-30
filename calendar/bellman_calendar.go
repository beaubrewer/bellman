package calendar

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/beaubrewer/bellman/config"
	"golang.org/x/oauth2/google"
	gcalendar "google.golang.org/api/calendar/v3"
)

// BellmanEvent maps google calendar events to Bellman events
// Start and End times are set correctly for all-day events
type BellmanEvent struct {
	Title string
	Type  string              `yaml:"type"`
	Audio map[string][]string `yaml:"audio"`
	Start time.Time
	End   time.Time
}

// BellmanCalendar provides the Google Calendar service
type BellmanCalendar struct {
	srv *gcalendar.Service
}

// NewBellmanEvent maps Google Calendar Events to BellmanEvents
// The Description (body) of the event is parsed into relevant properties
func NewBellmanEvent(event *gcalendar.Event) *BellmanEvent {
	var e BellmanEvent
	if err := yaml.Unmarshal([]byte(event.Description), &e); err != nil {
		fmt.Printf("There was a problem unmarshalling the event\n%s\n", err)
	}
	e.Title = event.Summary
	// check to see if this is a full day event
	// full day events do not provide timezone information so we call ParseInLocation
	// and provide the timezone/location data
	if len(event.Start.Date) > 0 {
		z := time.Now().Location()
		e.Start, _ = time.ParseInLocation("2006-01-02", event.Start.Date, z)
		e.End, _ = time.ParseInLocation("2006-01-02", event.End.Date, z)
	} else {
		e.Start, _ = time.Parse(time.RFC3339, event.Start.DateTime)
		e.End, _ = time.Parse(time.RFC3339, event.End.DateTime)
	}
	fmt.Printf("Bellman Event Title: %s\nStart: %v\nEnd: %v\n", e.Title, e.Start, e.End)
	return &e
}

// NewBellmanCalendar returns the Google Calendar API client for a particular
// Calendar
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

// GetEvents returns events that fall within the provided duration
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
	// TODO... if bellman starts within an event that is not a full day (between the start/end)
	// it will not set the theme correctly. We could retrieve the last X events and get the latest
	// event if desired
	return events
}
