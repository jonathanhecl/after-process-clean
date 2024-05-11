package main

import (
	"fmt"
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

	for range ticker.C {
		listProcess = process.List()

		if len(listProcess) != len(initProcess) {
			fmt.Println("Process list changed")
			break
		}
	}

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

	fmt.Println("Done")

}
