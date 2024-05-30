package main

import (
	"os"
	"strconv"
	"time"
)

type Package struct {
	Name      string
	Files     []string
	CreatedAt time.Time
}

func (p *Package) totalSize() int64 {
	var totalBytes int64 = 0

	for _, v := range p.Files {
		stat, err := os.Stat(v)

		if err != nil {
			continue
		}

		totalBytes += stat.Size()
	}

	return totalBytes
}

func (p *Package) display() string {
	return "Package: " + p.Name + " " + strconv.Itoa(int(p.totalSize()/(1024*1024))) + " MB"
}
