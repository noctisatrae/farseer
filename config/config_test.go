package config_test

import (
	"testing"

	config "github.com/noctisatrae/farseer/config"

	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	conf, err := config.Load("../config.toml")
	assert.NoError(t, err)

	// always-the-same option test
	assert.Equal(t, config.HubParams{
		PublicHubIp:     "92.158.95.48",
		GossipPort:      2282,
		RpcPort:         2283,
		BootstrapPeers:  []string{},
		Debug:           false,
		BufferSize:      128,
		ContactInterval: 3000,
	}, conf.Hub)
}
