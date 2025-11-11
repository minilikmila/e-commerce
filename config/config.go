package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config: holds the application configuration values.
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

type AppConfig struct {
	Name        string `mapstructure:"name"`
	Environment string `mapstructure:"environment"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"`
}

type JWTConfig struct {
	Secret          string        `mapstructure:"secret"`
	Issuer          string        `mapstructure:"issuer"`
	AccessTokenTTL  time.Duration `mapstructure:"access_token_ttl"`
	RefreshTokenTTL time.Duration `mapstructure:"refresh_token_ttl"`
}

// Load loads the configuration from the provided path (directory). It falls back to the current working directory.
func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	if path != "" {
		v.AddConfigPath(path)
	}
	v.AddConfigPath(".")

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed reading config file: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	applyFallbacks(&cfg)

	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("app.name", "ecommerce-api")
	v.SetDefault("app.environment", "development")

	v.SetDefault("server.port", 8080)

	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "postgres")
	v.SetDefault("database.name", "ecommerce")
	v.SetDefault("database.sslmode", "disable")

	v.SetDefault("jwt.secret", "change-this-secret")
	v.SetDefault("jwt.issuer", "ecommerce-api")
	v.SetDefault("jwt.access_token_ttl", time.Minute*30)
	v.SetDefault("jwt.refresh_token_ttl", time.Hour*24*7)
}

func applyFallbacks(cfg *Config) {
	if cfg.JWT.AccessTokenTTL == 0 {
		cfg.JWT.AccessTokenTTL = time.Minute * 30
	}

	if cfg.JWT.RefreshTokenTTL == 0 {
		cfg.JWT.RefreshTokenTTL = time.Hour * 24 * 7
	}

	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
}
