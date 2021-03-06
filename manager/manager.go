// Copyright © 2018 Beau Brewer <beaubrewer@gmail.com>

package manager

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/beaubrewer/bellman/calendar"
	"github.com/beaubrewer/bellman/manager/jobs"
	"github.com/beaubrewer/bellman/manager/mqtt"
	"github.com/beaubrewer/bellman/manager/timer"
)

var getBellmanEventsJob *jobs.GetBellmanEventsJob

//Start is the server entry point
func Start() {
	fmt.Println("Starting Bellman...")
	quit := configureSignals()
	calendarEvents := make(chan []*calendar.BellmanEvent)

	//start the MQTT consumer
	mqtt.Start()

	//start the job to fetch events from Google and process
	//bellman events
	getBellmanEventsJob = &jobs.GetBellmanEventsJob{
		Events:       calendarEvents,
		AudioUpdater: timer.NewAudioUpdater(),
	}
	s := timer.New()
	s.Every(1 * time.Hour).Do(getBellmanEventsJob)
	s.Start()

	//audioupdater updates the audio catalog on the events scheduled date/time
	getBellmanEventsJob.AudioUpdater.Start()

	fmt.Println("...at your service")
	<-quit
}

func configureSignals() <-chan bool {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigs
		fmt.Printf("Shutting down Bellman...\n")
		//any cleanup needed
		mqtt.Stop()
		getBellmanEventsJob.AudioUpdater.Stop()
		done <- true
	}()
	return done
}
