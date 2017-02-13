package shutdown

import (
	"os"
	"os/signal"
	"sort"
	"time"
)

var interrupt chan os.Signal
var halt chan bool
var registry chan bool
var registration chan int

var shutdown map[int]chan bool
var waiting map[int]chan bool
var registers map[int]int

func WaitForTimeout(timeout time.Duration, signals ...os.Signal) {
	Timeout(timeout)
	WaitFor(signals...)
}

func WaitFor(signals ...os.Signal) {
	go register()

	signal.Notify(interrupt, signals...)

	select {
	case <-interrupt:
	case <-halt:
	}

	/* Shutdown the registry */
	registry <- true

	/* Send the shutdown signal to the rest of the workers */
	sendSignals()
}

func Timeout(timeout time.Duration) {
	/* In Duration, register a shutdown */
	go func(d time.Duration) {
		<-time.After(d)
		Now()
	}(timeout)
}

/* Register a shutdown NOW! */
func Now() {
	halt <- true
}

/* Wait for a worker to finish! */
func WaitForMe(register int) *worker {
	registration <- register
	/* Make sure the proper shutdown channel exists */
	if _, exists := shutdown[register]; !exists {
		shutdown[register] = make(chan bool)
	}
	/* Make sure the proper waiting channel exists */
	if _, exists := waiting[register]; !exists {
		waiting[register] = make(chan bool)
	}
	return &worker{priority: register, shutdown: shutdown[register], finished: waiting[register]}
}

/* Start our register */
func register() {
	working := true
	for working {
		select {
		case p := <-registration:
			registers[p]++
		case <-registry:
			working = false
		}
	}
}

/* Send a signal to the channels */
func sendSignals() {
	/* Save for later ... sorting */
	var priorities []int

	/* Get our priority list */
	for priority, _ := range registers {
		priorities = append(priorities, priority)
	}

	/* Make sure we sort them in numerical order */
	sort.Ints(priorities)

	/* Go through each of them and send the shutdown signal */
	for _, priority := range priorities {
		for x := 0; x < registers[priority]; x++ {
			shutdown[priority] <- true
		}
		for x := 0; x < registers[priority]; x++ {
			<-waiting[priority]
		}
	}
}

/* Giddy up! */
func init() {
	halt = make(chan bool)
	interrupt = make(chan os.Signal, 1)
	registration = make(chan int)
	registers = make(map[int]int)
	shutdown = make(map[int]chan bool)
	waiting = make(map[int]chan bool)
	registry = make(chan bool, 1)
}
