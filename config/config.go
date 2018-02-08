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
	// DefaultLockSpan is used when max_lock_span is not set
	DefaultLockSpan = 10 * time.Second
	// DefaultCleanUpDelay is used when default_cleanup_span is not set
	DefaultCleanUpDelay = 5 * time.Second
)

// Config holds the config.yml values
type Config struct {
	MaxRetries   int           `yaml:"max_retries"`
	RetryDelay   time.Duration `yaml:"retry_delay"`
	MaxLockSpan  time.Duration `yaml:"max_lock_span"`
	CleanUpDelay time.Duration `yaml:"default_cleanup_delay"`
}

// GetConfig returns a new Config instance
func GetConfig() (Config, error) {
	c := Config{
		MaxRetries:   DefaultMaxRetries,
		RetryDelay:   DefaultRetryDelay,
		MaxLockSpan:  DefaultLockSpan,
		CleanUpDelay: DefaultCleanUpDelay,
	}
	configFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return c, err
	}

	err = yaml.Unmarshal(configFile, &c)

	if err != nil {
		return c, err
	}

	return c, nil
}
