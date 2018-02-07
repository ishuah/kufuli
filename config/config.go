package config

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	// DefaultMaxRetries is used when max_retries is not set
	DefaultMaxRetries = 16
	// DefaultRetryDelay is used when retry_delay is not set
	DefaultRetryDelay = 500 * time.Millisecond
)

// Config holds the config.yml values
type Config struct {
	MaxRetries int           `yaml:"max_retries"`
	RetryDelay time.Duration `yaml:"retry_delay"`
}

// GetConfig returns a new Config instance
func GetConfig() (Config, error) {
	configFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return Config{MaxRetries: DefaultMaxRetries,
			RetryDelay: DefaultRetryDelay}, err
	}

	var c Config
	err = yaml.Unmarshal(configFile, &c)

	if err != nil {
		return Config{MaxRetries: DefaultMaxRetries,
			RetryDelay: DefaultRetryDelay}, err
	}

	return c, nil
}
