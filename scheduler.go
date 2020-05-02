package scrape

import (
	"math/rand"
	"time"
)

// Define Intervals for convenient usage.
// Accuracy of time is not the intention here, but best approximation given
const (
	Daily         = iota
	Weekly        = iota
	Monthly       = iota
	RandomWeekDay = iota
	RandomHour    = iota
)

const (
	epoch = iota /* Consider what this needs to be later, an hour? 10 minutes? 1? */
)

type enum int

// Process can be scheduled for specific times and woken
type Process struct {
	Start     func()
	Interval  enum
	startTime time.Time
	runCount  int64
}

func schedule(this *Process) {
	this.startTime = time.Now().Add(timeFrom(this.Interval))
}

// Run a process and signal when it has completed one iteration
func Run(P *Process) (done chan bool) {
	done = make(chan bool, 1)

	go func() {
		P.runCount++
		P.Start()
		done <- true
	}()

	return done
}

func timeFrom(interval enum) time.Duration {
	day := time.Hour * 24

	/* There are several methods for creating a Duration, be cautioned against odd styles */
	switch interval {
	case Daily:
		return time.Duration(day)
	case Weekly:
		return time.Duration(day * 7)
	case Monthly:
		return time.Duration(day * 30)
	case RandomWeekDay:
		return time.Duration(rand.Intn(7)+1) * time.Hour * 24
	case RandomHour:
		return time.Duration(rand.Intn(24)) * time.Hour
	default:
		return time.Duration(5) * time.Minute /* A better alternative may be to set very long duration? */
	}
}

func sleep(thisLong int) {
	time.Sleep(time.Duration(thisLong) * time.Minute)
}

// StartScheduler may not be needed! But it fits the model of the world I expect, so run it
func StartScheduler() {
	/*
		Every epoch, the scheduler will check if the next-in-line process is ready to go
		Then it starts the process which should be self-contained regarding resources
		After starting the process

		example usage:
			scrape.Process(scrape.City2CityBasic, scheduler.DAILY)
			Schedule(fxn, interval)
	*/

	var Q = new(ProcessQueue)

	for {
		/* Wake up and check next Process for start time */
		もうそろそろになた := time.Now().After(Q.head().startTime)

		if もうそろそろになた {
			go func() {
				process := Q.deQ()
				Run(&process)

				/* Put back in line */
				process.startTime.Add(timeFrom(process.Interval))
				Q.enQ(process)
			}()
		}
		sleep(epoch)
	}
}

// ProcessQueue is a half-baked Queue that will for now just hold processes in case iteration is desired later.
// Ideally improve by enforcing priority by time (Heap seems like good fit)
// May NOT need to use Queue for scheduling purposes if each go routine can manage their own sleep
type ProcessQueue struct {
	members []Process
}

func (Q ProcessQueue) empty() bool {
	return len(Q.members) == 0
}

func (Q ProcessQueue) head() Process {
	if Q.empty() {
		NoOperation := Process{}
		return NoOperation
	}
	return Q.members[0]
}

func (Q ProcessQueue) enQ(P Process) {
	Q.members = append(Q.members, P)
}

func (Q ProcessQueue) deQ() Process {
	if Q.empty() {
		NoOperation := Process{}
		return NoOperation
	}
	nextProc := Q.pop(0)

	return nextProc
}

func (Q ProcessQueue) pop(index int) Process {
	popcorn := Q.members[index]
	Q.members = append(Q.members[:index], Q.members[index+1:]...)

	return popcorn
}
