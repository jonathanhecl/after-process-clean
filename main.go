package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/afterprocessclean/process"
)

const (
	version = "1.0.0"
)

func main() {
	fmt.Println("AfterProcessClean v" + version)

	control.UpdateList(process.List(), true)
	scanning := true

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	go func() {
		<-c

		scanning = false
		fmt.Println("Exiting...")

		ticker.Stop()

		for _, p := range control.AfterList() {
			fmt.Println("New process: ", p.Path, p.CRC32, time.Since(p.RuningSince))
		}

		fmt.Print("Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')

		os.Exit(1)
	}()

	for {
		select {
		case <-ticker.C:
			if scanning {
				control.UpdateList(process.List(), false)
			}
		}
	}
}
