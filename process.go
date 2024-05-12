package main

import (
	"fmt"
	"time"

	"github.com/afterprocessclean/process"
)

type processStruct struct {
	Path        string
	CRC32       string
	RuningSince time.Time
	Before      bool
}

type controlStruct struct {
	processes []*processStruct
}

var control controlStruct

func init() {
	control.processes = []*processStruct{}
}

func (c *controlStruct) addProcess(path string, crc32 string, before bool) {
	c.processes = append(c.processes, &processStruct{
		Path:        path,
		CRC32:       crc32,
		RuningSince: time.Now(),
		Before:      before,
	})

	if !before {
		fmt.Println(time.Now(), "Added process: ", path)
	}
}

func (c *controlStruct) removeProcess(path string) {
	for i, p := range c.processes {
		if p.Path == path {
			c.processes = append(c.processes[:i], c.processes[i+1:]...)
			break
		}
	}

	fmt.Println(time.Now(), "Removed process: ", path)
}

func (c *controlStruct) getProcess(path string) *processStruct {
	for _, p := range c.processes {
		if p.Path == path {
			return p
		}
	}

	return nil
}

func (c *controlStruct) updateProcess(path string, crc32 string) {
	for _, p := range c.processes {
		if p.Path == path {
			p.CRC32 = crc32
			break
		}
	}
}

func (c *controlStruct) AfterList() (list []*processStruct) {
	for _, p := range c.processes {
		if p.Before {
			continue
		}

		list = append(list, &processStruct{
			Path:        p.Path,
			CRC32:       p.CRC32,
			RuningSince: p.RuningSince,
		})
	}

	return
}

func (c *controlStruct) UpdateList(list []process.ProcessStruct, before bool) {
	for _, p := range list {
		if p.Path == "" {
			continue
		}

		if c.getProcess(p.Path) == nil {
			c.addProcess(p.Path, "", before)
		} else {
			c.updateProcess(p.Path, "")
		}
	}

	for _, p := range c.processes {
		found := false
		for _, l := range list {
			if l.Path == p.Path {
				found = true
				break
			}
		}

		if !found {
			c.removeProcess(p.Path)
		}
	}
}
