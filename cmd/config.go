package main

import (
	"github.com/cafebazaar/sentry-gateway/metrics"
	"github.com/cafebazaar/sentry-gateway/reverseproxy"
	"github.com/cafebazaar/sentry-gateway/throttle"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	ReverseProxyConfig reverseproxy.Config `mapstructure:"proxy"`
	ThrottleConfig     throttle.Config     `mapstructure:"throttle"`
	MetricsConfig      metrics.Config      `mapstructure:"metrics"`
	ListenAddress      string
	LogLevel           string
}

func LoadConfig() (Config, error) {
	var config Config

	configFilePath, exists := os.LookupEnv("SENTRY_GATEWAY_CONFIG_FILE_PATH")
	if ! exists {
		configFilePath = "./config.yaml"
	}
	// Read Config from file
	viper.SetConfigFile(configFilePath)

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, err
	}

	logLevel, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		return Config{}, err
	}

	logrus.SetLevel(logLevel)
	return config, nil
}
