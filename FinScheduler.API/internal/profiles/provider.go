package profiles

import (
	"github.com/grafana/pyroscope-go"
)

func InitProfiler() (*pyroscope.Profiler, error) {
	prof, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: "finscheduler",
		//ServerAddress: "http://localhost:4040",
	})

	return prof, err
}
