package util

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestInitConfig(t *testing.T) {
	os.Clearenv()
	os.Setenv("HD_USERNAME", "myusername")
	os.Setenv("HD_PASSWORD", "mypassword")
	os.Setenv("HD_DATAPATH", "mydatapath")
	os.Setenv("HD_JWTSECRET", "myjwtsecret")

	if assert.NoError(t, InitConfig()) {
		assert.Equal(t, AppConfig.Username, "myusername")
		assert.Equal(t, AppConfig.Password, "mypassword")
		assert.Equal(t, AppConfig.DataPath, "mydatapath")
		assert.Equal(t, AppConfig.JwtSecret, "myjwtsecret")
	}
}
