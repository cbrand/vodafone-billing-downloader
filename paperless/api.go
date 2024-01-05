package paperless

import (
	"encoding/base64"
	"fmt"
)

type Config struct {
	URL      string
	APIKey   string
	Username string
	Password string
}

func (config *Config) GetAuthorizationHeader() string {
	if len(config.APIKey) > 0 {
		return fmt.Sprintf("Token %s", config.APIKey)
	} else {
		usernamePassword := fmt.Sprintf("%s:%s", config.Username, config.Password)
		return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(usernamePassword)))
	}
}
