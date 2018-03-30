// Copyright © 2018 Beau Brewer <beaubrewer@gmail.com>

package jobs

import (
	"time"

	"github.com/beaubrewer/bellman/calendar"
	"github.com/beaubrewer/bellman/manager/timer"
)

var client *calendar.BellmanCalendar
var bellmanEvents []*calendar.BellmanEvent

// GetBellmanEventsJob will get 3 days of events
type GetBellmanEventsJob struct {
	Events       chan []*calendar.BellmanEvent
	AudioUpdater *timer.AudioUpdater
}

// Run will return 3 days of events
func (j GetBellmanEventsJob) Run() {
	if client == nil {
		client = calendar.NewBellmanCalendar()
	}
	bellmanEvents = nil
	e := client.GetEvents(time.Hour * 72)
	for _, ev := range e.Items {
		be := calendar.NewBellmanEvent(ev)
		bellmanEvents = append(bellmanEvents, be)
	}
	j.AudioUpdater.UpdateBellmanEvents(bellmanEvents)
}
