package config

import (
	"encoding/base64"

	"github.com/kelseyhightower/envconfig"
)

// Config represents env values
type Config struct {
	DogStatsD     string `envconfig:"DOGSTATSD"`
	RsaPrivateKey string `split_words:"true"`
}

// New initializes and returns a new config
func New() Config {
	var cfg Config
	envconfig.MustProcess("", &cfg)
	return cfg
}

// ParseRsaPrivateKeyHex parses the RsaPrivateKey hex and return bytes slice
func ParseRsaPrivateKeyHex(k string) ([]byte, error) {
	rsaBytes, err := base64.StdEncoding.DecodeString(k)
	if err != nil {
		return nil, err
	}

	return rsaBytes, nil
}
