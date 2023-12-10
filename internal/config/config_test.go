package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func setEnvs(envVars map[string]string) {
	for key, value := range envVars {
		os.Setenv(key, value)
	}
}

func TestGetServerConfig(t *testing.T) {
	// set the environment variables
	setEnvs(map[string]string{
		"SERVER_ADDR": "",
		"SERVER_PORT": "",
	})

	// call tested function
	cfg := GetServerConfig()

	// assertions
	require.Equal(t, "localhost", cfg.Address)
	require.Equal(t, "8080", cfg.Port)
}

func TestGetClientConfig(t *testing.T) {
	setEnvs(map[string]string{
		"CLIENT_ID":      "test",
		"CLIENT_ROOM_ID": "",
	})
	cfg, err := GetClientConfig()

	require.NoError(t, err)
	require.Equal(t, "test", cfg.ClientID)
	require.Equal(t, "default", cfg.RoomID)
}

func TestGetClientConfigMissingClientID(t *testing.T) {
	setEnvs(map[string]string{
		"CLIENT_ID":      "",
		"CLIENT_ROOM_ID": "default",
	})

	_, err := GetClientConfig()

	require.Error(t, err)
	require.Equal(t, "CLIENT_ID env variable is empty", err.Error())
}
