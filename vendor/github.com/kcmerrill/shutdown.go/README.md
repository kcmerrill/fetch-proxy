#Shutdown.go
A really simple library to help handle the registring of subroutines in go(in priority order)


#Example Code
```
package main

import (
	"fmt"
	"github.com/kcmerrill/shutdown.go"
	"syscall"
	"time"
)

func main() {
    /* Start our application */
	fmt.Println("Listening for shutdown ...")

	/* simulate a shutdown.Now() */
	go DoSomethingToTriggerAShutdown()
	for x := 10000; x >= 0; x-- {
		go worker(x)
	}

	/* Simulate a timeout */
	//shutdown.Timeout("10m")
	//shutdown.WaitFor(syscall.SIGINT, syscall.SIGTERM)
	shutdown.WaitForTimeout("10m", syscall.SIGINT, syscall.SIGTERM) // this is the same as the above two lines
	fmt.Println("\nShutting down ...")
}

func worker(id int) {
	/* Wait for this worker to finish, with a high priority(the lower the number, the faster the shutdown */
	worker := shutdown.WaitForMe(id)

	/* You'd typically put this into your select {} */
	<-worker.Stop()

	/* Tell the world we are ready to shutdown */
	worker.Finished()
}


func DoSomethingToTriggerAShutdown() {
	<-time.After(30 * time.Second)
	fmt.Println("Shutdown() triggered")
	shutdown.Now()
}
```
