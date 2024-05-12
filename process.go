package main

import (
	"time"

	"github.com/afterprocessclean/process"
)

type myProcess struct {
	Path        string
	CRC32       string
	RuningSince time.Time
	Before      bool
}

var (
	processes []myProcess
)

func saveInit([]process.ProcessStruct) {
	processes = nil
	for _, p := range processes {
		processes = append(processes, myProcess{
			Path:        p.Path,
			CRC32:       "",
			RuningSince: time.Now(),
			Before:      true,
		})
	}
}
