package login

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func loginWithEnvironmentVariables() (*LoginData, error) {
	return FromEnvironmentVariables()
}

func TestLogin(t *testing.T) {
	credentials, err := loginWithEnvironmentVariables()
	assert.Nil(t, err)
	assert.NotNil(t, credentials)
	fmt.Println(credentials)
}
