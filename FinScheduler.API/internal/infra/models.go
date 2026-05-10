package infra

type Config struct {
	Env              string
	ServerPort       int
	ConnectionString string
	CORSSettings     CORSSettings
	Observability    ObservabilityConfig
}

type CORSSettings struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
}

type ObservabilityConfig struct {
	ServiceName string
	Metrics     MetricsConfig
	Traces      TracesConfig
	Profiling   ProfilingConfig
}

type MetricsConfig struct {
	Enabled        bool
	ExportEndpoint string
}

type TracesConfig struct {
	Enabled                bool
	ExportEndpoint         string
	RootTraceSamplingRatio float64
}

type ProfilingConfig struct {
	Enabled bool
	PushURL string
}
