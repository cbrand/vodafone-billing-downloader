package login

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserInfo(t *testing.T) {
	credentials, err := loginWithEnvironmentVariables()
	assert.Nil(t, err)
	assert.NotNil(t, credentials)
	userInfo, err := GetUserInfo(credentials)
	assert.Nil(t, err)
	assert.NotNil(t, userInfo)
	fmt.Println(userInfo)
}
