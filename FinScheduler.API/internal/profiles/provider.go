package profiles

import (
	"fmt"
	"strings"

	"finscheduler/internal/infra"

	"github.com/grafana/pyroscope-go"
)

func InitProfiler(cfg *infra.Config) (*pyroscope.Profiler, error) {
	if !cfg.Observability.Profiling.Enabled {
		return nil, nil
	}

	pushURL := strings.TrimSpace(cfg.Observability.Profiling.PushURL)
	if pushURL == "" {
		return nil, fmt.Errorf("profiling enabled but profiling push URL is empty")
	}

	prof, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: cfg.Observability.ServiceName,
		ServerAddress:   pushURL,
		Tags: map[string]string{
			"env": cfg.Env,
		},
	})

	return prof, err
}
