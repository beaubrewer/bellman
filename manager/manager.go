package manager

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/beaubrewer/bellmanv2/calendar"
)

// type ScheduledEvents struct {
// 	Events []Event
// }

//Start is the server entry point
func Start() {
	quit := configureSignals()
	//create a calendar
	c := calendar.NewBellmanCalendar()
	c.GetCurrentTheme()
	//a := NewAudioManager()

	for {
		select {
		case <-quit:
			return
		}
	}
}

func configureSignals() <-chan bool {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigs
		fmt.Printf("Shutting down Bellman...")
		//any cleanup needed
		done <- true
	}()
	return done
}
