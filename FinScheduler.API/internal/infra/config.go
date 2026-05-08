package infra

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Env              string
	ServerPort       int
	ConnectionString string
	ServiceName      string
	CORSSettings     CORSSettings
}

type CORSSettings struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
}

func LoadConfig() (*Config, error) {
	v := viper.New()
	if err := loadJSONConfig(v); err != nil {
		return nil, err
	}

	if err := loadEnvFile(v); err != nil {
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

func loadEnvFile(v *viper.Viper) error {
	v.SetConfigFile(".env")
	v.SetConfigType("env")

	if err := v.MergeInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return fmt.Errorf("read .env: %w", err)
		}
	}

	return nil
}

func configureEnvironment(v *viper.Viper) {
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetDefault("corsSettings.allowedOrigins", []string{"*"})
	v.SetDefault("corsSettings.allowedMethods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	v.SetDefault("corsSettings.allowedHeaders", []string{"*"})
	v.SetDefault("corsSettings.allowCredentials", false)
	bindConfigEnv(v)
}

func bindConfigEnv(v *viper.Viper) {
	bindEnv(v, "env", "ENV")
	bindEnv(v, "serverPort", "SERVER_PORT")
	bindEnv(v, "connectionString", "CONNECTION_STRING")
	bindEnv(v, "serviceName", "SERVICE_NAME")
	bindEnv(v, "corsSettings.allowedOrigins", "CORS_ALLOWED_ORIGINS")
	bindEnv(v, "corsSettings.allowedMethods", "CORS_ALLOWED_METHODS")
	bindEnv(v, "corsSettings.allowedHeaders", "CORS_ALLOWED_HEADERS")
	bindEnv(v, "corsSettings.allowCredentials", "CORS_ALLOW_CREDENTIALS")
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
	cfg.ServiceName = strings.TrimSpace(v.GetString("serviceName"))
	cfg.CORSSettings.AllowedOrigins = resolveStringList(v, "corsSettings.allowedOrigins", cfg.CORSSettings.AllowedOrigins)
	cfg.CORSSettings.AllowedMethods = resolveStringList(v, "corsSettings.allowedMethods", cfg.CORSSettings.AllowedMethods)
	cfg.CORSSettings.AllowedHeaders = resolveStringList(v, "corsSettings.allowedHeaders", cfg.CORSSettings.AllowedHeaders)
	cfg.CORSSettings.AllowCredentials = v.GetBool("corsSettings.allowCredentials")
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
