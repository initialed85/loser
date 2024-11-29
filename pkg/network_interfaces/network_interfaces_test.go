package network_interfaces

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetNetworkInterfaces(t *testing.T) {
	networkInterfaces, err := GetNetworkInterfaces()
	require.NoError(t, err)
	require.NotNil(t, networkInterfaces)

	b, err := json.MarshalIndent(networkInterfaces, "", "  ")
	require.NoError(t, err)
	log.Printf("networkInterfaces: %s", string(b))
}
