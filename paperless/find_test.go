package paperless

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func configFromEnviroment() *Config {
	return &Config{
		URL:      os.Getenv("PAPERLESS_URL"),
		APIKey:   os.Getenv("PAPERLESS_API_KEY"),
		Username: os.Getenv("PAPERLESS_USERNAME"),
		Password: os.Getenv("PAPERLESS_PASSWORD"),
	}

}

func TestChecksum(t *testing.T) {
	config := configFromEnviroment()
	exists, err := ChecksumExists(config, "5fd53d226dd3107c4289f9c24297cf0a")
	assert.Nil(t, err)
	if exists {
		assert.True(t, exists)
	} else {
		assert.False(t, exists)
	}
}
