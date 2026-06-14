package infra

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

func LoadConfig() (*Config, error) {
	v := viper.New()
	if err := loadJSONConfig(v); err != nil {
		return nil, err
	}

	if err := loadDotEnv(); err != nil {
		return nil, err
	}

	configureEnvironment(v)

	cfg := new(Config)
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	applyEnvOverrides(v, cfg)

	return cfg, nil
}

func loadJSONConfig(v *viper.Viper) error {
	v.SetConfigName("config")
	v.SetConfigType("json")
	v.AddConfigPath("./configs")

	if err := v.MergeInConfig(); err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	return nil
}

func loadDotEnv() error {
	if err := gotenv.Load(".env"); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("read .env: %w", err)
		}
	}

	return nil
}

func configureEnvironment(v *viper.Viper) {
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetDefault("corsSettings.allowedOrigins", []string{"*"})
	v.SetDefault("corsSettings.allowedMethods", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"})
	v.SetDefault("corsSettings.allowedHeaders", []string{"*"})
	v.SetDefault("corsSettings.allowCredentials", false)
	v.SetDefault("observability.serviceName", "fin-scheduler-api")
	v.SetDefault("observability.metrics.enabled", true)
	v.SetDefault("observability.metrics.exportEndpoint", "/metrics")
	v.SetDefault("observability.traces.enabled", false)
	v.SetDefault("observability.traces.exportEndpoint", "http://localhost:4318")
	v.SetDefault("observability.traces.rootTraceSamplingRatio", 1.0)
	v.SetDefault("observability.profiling.enabled", false)
	v.SetDefault("observability.profiling.pushURL", "http://localhost:4040")
	bindConfigEnv(v)
}

func bindConfigEnv(v *viper.Viper) {
	bindEnv(v, "env", "ENV")
	bindEnv(v, "serverPort", "SERVER_PORT")
	bindEnv(v, "connectionString", "CONNECTION_STRING")
	bindEnv(v, "corsSettings.allowedOrigins", "CORS_ALLOWED_ORIGINS")
	bindEnv(v, "corsSettings.allowedMethods", "CORS_ALLOWED_METHODS")
	bindEnv(v, "corsSettings.allowedHeaders", "CORS_ALLOWED_HEADERS")
	bindEnv(v, "corsSettings.allowCredentials", "CORS_ALLOW_CREDENTIALS")
	bindEnv(v, "observability.serviceName", "OBSERVABILITY_SERVICE_NAME")
	bindEnv(v, "observability.metrics.enabled", "METRICS_ENABLED")
	bindEnv(v, "observability.metrics.exportEndpoint", "METRICS_EXPORT_ENDPOINT")
	bindEnv(v, "observability.traces.enabled", "TRACES_ENABLED")
	bindEnv(v, "observability.traces.exportEndpoint", "TRACES_EXPORT_ENDPOINT")
	bindEnv(v, "observability.traces.rootTraceSamplingRatio", "TRACES_ROOT_TRACE_SAMPLING_RATIO")
	bindEnv(v, "observability.profiling.enabled", "PROFILING_ENABLED")
	bindEnv(v, "observability.profiling.pushURL", "PROFILING_PUSH_URL")
}

func bindEnv(v *viper.Viper, key string, envNames ...string) {
	args := append([]string{key}, envNames...)
	if err := v.BindEnv(args...); err != nil {
		panic(fmt.Sprintf("bind env %s: %v", key, err))
	}
}

func applyEnvOverrides(v *viper.Viper, cfg *Config) {
	cfg.Env = strings.TrimSpace(v.GetString("env"))

	if serverPort := v.GetInt("serverPort"); serverPort != 0 {
		cfg.ServerPort = serverPort
	}

	cfg.ConnectionString = strings.TrimSpace(v.GetString("connectionString"))
	cfg.CORSSettings.AllowedOrigins = resolveStringList(v, "corsSettings.allowedOrigins", cfg.CORSSettings.AllowedOrigins)
	cfg.CORSSettings.AllowedMethods = resolveStringList(v, "corsSettings.allowedMethods", cfg.CORSSettings.AllowedMethods)
	cfg.CORSSettings.AllowedHeaders = resolveStringList(v, "corsSettings.allowedHeaders", cfg.CORSSettings.AllowedHeaders)
	cfg.CORSSettings.AllowCredentials = v.GetBool("corsSettings.allowCredentials")
	cfg.Observability.ServiceName = strings.TrimSpace(v.GetString("observability.serviceName"))
	cfg.Observability.Metrics.Enabled = v.GetBool("observability.metrics.enabled")
	cfg.Observability.Metrics.ExportEndpoint = normalizeHTTPPath(v.GetString("observability.metrics.exportEndpoint"), cfg.Observability.Metrics.ExportEndpoint)
	cfg.Observability.Traces.Enabled = v.GetBool("observability.traces.enabled")
	cfg.Observability.Traces.ExportEndpoint = strings.TrimSpace(v.GetString("observability.traces.exportEndpoint"))
	cfg.Observability.Traces.RootTraceSamplingRatio = resolveSampleRatio(v.GetFloat64("observability.traces.rootTraceSamplingRatio"), cfg.Observability.Traces.RootTraceSamplingRatio)
	cfg.Observability.Profiling.Enabled = v.GetBool("observability.profiling.enabled")
	cfg.Observability.Profiling.PushURL = strings.TrimSpace(v.GetString("observability.profiling.pushURL"))
}

func resolveStringList(v *viper.Viper, key string, fallback []string) []string {
	rawValue := strings.TrimSpace(v.GetString(key))
	if rawValue == "" {
		return fallback
	}

	parts := strings.Split(rawValue, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmedPart := strings.TrimSpace(part)
		if trimmedPart == "" {
			continue
		}

		values = append(values, trimmedPart)
	}

	if len(values) == 0 {
		return fallback
	}

	return values
}

func normalizeHTTPPath(rawValue string, fallback string) string {
	value := strings.TrimSpace(rawValue)
	if value == "" {
		value = strings.TrimSpace(fallback)
	}
	if value == "" {
		return "/metrics"
	}
	if strings.HasPrefix(value, "/") {
		return value
	}

	return "/" + value
}

func resolveSampleRatio(value float64, fallback float64) float64 {
	if value >= 0 && value <= 1 {
		return value
	}
	if fallback >= 0 && fallback <= 1 {
		return fallback
	}

	return 1
}
