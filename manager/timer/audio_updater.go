package timer

import (
	"fmt"
	"sync"
	"time"

	"github.com/beaubrewer/bellmanv2/calendar"
	"github.com/beaubrewer/bellmanv2/manager/audio"
)

// AudioUpdater holds the current events
type AudioUpdater struct {
	stop          chan struct{}
	running       bool
	currentEvent  *calendar.BellmanEvent
	bellmanEvents []*calendar.BellmanEvent
	mutex         sync.Mutex
}

// NewAudioUpdater returns an AudioUpdater
func NewAudioUpdater() *AudioUpdater {
	return &AudioUpdater{
		stop:          make(chan struct{}),
		running:       false,
		currentEvent:  nil,
		bellmanEvents: nil,
	}
}

// UpdateBellmanEvents replaces the bellman events collection
func (a *AudioUpdater) UpdateBellmanEvents(events []*calendar.BellmanEvent) {
	a.mutex.Lock()
	a.bellmanEvents = events
	a.mutex.Unlock()
}

func (a *AudioUpdater) run() {
	ticker := time.NewTicker(100 * time.Millisecond)
	go func() {
		for {
			select {
			case <-ticker.C:
				a.checkForUpdates()
				continue
			case <-a.stop:
				ticker.Stop()
				return
			}
		}
	}()
}

func (a *AudioUpdater) checkForUpdates() {
	//fmt.Print(".")
	//This handles most of the business logic for updating the audio catalog
	//It checks the start/end dates and determines if it needs to update
	//NOTE: The events are sorted by date currently. This happens in the get_calendar_events job.
	//This allows us to only check the first event in the collection.
	a.mutex.Lock()
	now := time.Now()
	if len(a.bellmanEvents) != 0 {
		event := a.bellmanEvents[0]
		if now.After(event.Start) && now.Before(event.End) {
			//set the audio theme
			fmt.Printf("Setting the audio theme to %s\n", event.Title)
			audio.UpdateCatalog(event.Audio)
			a.currentEvent = event
			a.bellmanEvents = a.bellmanEvents[1:]
		}
	}
	if a.currentEvent != nil && now.After(a.currentEvent.End) {
		//remove the theme
		fmt.Println("Removing the audio theme")
		audio.UpdateCatalog(map[string][]string{})
		a.currentEvent = nil
	}
	a.mutex.Unlock()
}

// Start the AudioUpdater in its own go-routine, or no-op if already started.
func (a *AudioUpdater) Start() {
	if a.running {
		return
	}
	a.running = true
	go a.run()
}

// Stop the AudioUpdater if it is running.
func (a *AudioUpdater) Stop() {
	if !a.running {
		return
	}
	a.stop <- struct{}{}
	a.running = false
}
