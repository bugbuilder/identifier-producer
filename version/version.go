package version

import (
	"fmt"
	"runtime"
)

var (
	version   = "unknown"
	gitCommit = "unknown"
	buildDate = "unknown"
)

type info struct {
	Version   string
	GitCommit string
	BuildDate string
	GoVersion string
	Compiler  string
}

func NewInfo() *info {
	return &info{
		Version:   version,
		GitCommit: gitCommit,
		BuildDate: buildDate,
		GoVersion: runtime.Version(),
		Compiler:  runtime.Compiler,
	}
}

func (i info) Print() string {
	return fmt.Sprintf("%+v\n", i)
}
