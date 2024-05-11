package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/afterprocessclean/process"
)

const (
	version = "0.1.0"
)

func main() {
	fmt.Println("AfterProcessClean v" + version)

	initProcess := process.List()
	listProcess := []process.ProcessStruct{}

	ticker := time.NewTicker(5 * time.Second)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c

		// differences
		for _, p := range listProcess {
			isNew := true
			for _, i := range initProcess {
				if p.PID == i.PID {
					isNew = false
					break
				}
			}

			if isNew {
				fmt.Println("New process: ", p.PID, p.Path, p.Filename)
			}
		}

		os.Exit(1)
	}()

	for range ticker.C {
		listProcess = process.List()
	}

	fmt.Println("Done")

}
