package scrape

import (
	"time"
)

const (
	Daily	= iota
	Weekly	= iota
	Monthly	= iota
	RandomWeekDay	= iota
	RandomHour		= iota
)

const (
	epoch
)

type enum int
type Time time.Time

// Process can be scheduled for specific times and woken
type Process struct {
	Start func(chan bool)
	Interval enum
	startTime Time
}

func (P Process) Start() {

}

func timeFrom(interval enum) Time {

}

func sleep(thisLong int) {
	time.Sleep(time.Duration{int64(thisLong)})
}

func StartScheduler() {
	var Q = new(ProcessQueue)

	for {
		/* Wake up and check next Process for start time */	
		if time.Now().After(Q.head()) {
			go func(){
				sayWhenDone := make(chan bool)
				process := Q.DeQ()
				process.Start(done)
				process.startTime.Add(timeFrom(process.Interval))
				<-done
				Q.EnQ(process)
			}
		}
		sleep(epoch)
	}
}

/*
	Every epoch, the scheduler will check if the next-in-line process is ready to go
	Then it starts the process which should be self-contained regarding resources
	After starting the process

	example usage:
		scrape.Process(scrape.City2CityBasic, scheduler.DAILY)
		Schedule(fxn, interval)

*/