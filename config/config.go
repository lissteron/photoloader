package config

import (
	"errors"

	"github.com/spf13/viper"

	"github.com/lissteron/photoloader/pkg/ehttp"
)

const (
	DefaultLogLevel = "info"

	_defaultSoftMemoryLimitMB = 512
	_megabyte                 = 1024 * 1024
)

var ErrEmptyPhotoDir = errors.New("empty photo dir")

type Config struct {
	Base        BaseConfig
	HTTPServer  *ehttp.ServerConfig
	PhotoConfig *PhotoConfig
}

func initConfig() {
	viper.AutomaticEnv()

	// Base config
	viper.SetDefault("LOG_LEVEL", DefaultLogLevel)
	viper.SetDefault("ENV", "development")
	viper.SetDefault("SOFT_MEMORY_LIMIT_MB", _defaultSoftMemoryLimitMB)
	viper.SetDefault("HTTP_LISTEN_ADDR", ":8080")

	viper.SetDefault("PHOTO_PATH", "/photos")
}

func NewConfig() *Config {
	initConfig()

	return &Config{
		Base:        NewBaseConfig(),
		HTTPServer:  NewHTTPConfig(),
		PhotoConfig: NewPhotoConfig(),
	}
}

type BaseConfig struct {
	ENV             string
	LogLevel        string
	SoftMemoryLimit int64
	WithMigrations  bool
}

func NewBaseConfig() BaseConfig {
	return BaseConfig{
		ENV:             viper.GetString("ENV"),
		LogLevel:        viper.GetString("LOG_LEVEL"),
		SoftMemoryLimit: viper.GetInt64("SOFT_MEMORY_LIMIT_MB") * _megabyte,
		WithMigrations:  viper.GetBool("WITH_MIGRATIONS"),
	}
}

func NewHTTPConfig() *ehttp.ServerConfig {
	return &ehttp.ServerConfig{
		ListenAddr: viper.GetString("HTTP_LISTEN_ADDR"),
	}
}

type PhotoConfig struct {
	Path string // for server
	Dir  string // for storage
}

func NewPhotoConfig() *PhotoConfig {
	return &PhotoConfig{
		Path: viper.GetString("PHOTO_PATH"),
		Dir:  viper.GetString("PHOTO_DIR"),
	}
}

func (c *Config) Validate() error {
	if c.PhotoConfig.Dir == "" {
		return ErrEmptyPhotoDir
	}

	return nil
}
