package main

import (
	"testing"

	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	conf, err := Load("../config.toml")
	assert.NoError(t, err)

	postgreConf := conf.Handlers["postgresql"]

	// always-the-same option test
	assert.Equal(t, HubParams{
		GossipPort:      2282,
		BootstrapPeers:  []string{"/dns/nemes.farcaster.xyz/tcp/2282/p2p/12D3KooWMQrf6unpGJfLBmTGy3eKTo4cGcXktWRbgMnfbZLXqBbn"},
		Debug:           true,
		BufferSize:      128,
		ContactInterval: 30,
	}, conf.Hub)

	// dynamic conf
	assert.Equal(t, true, postgreConf.(map[string]interface{})["Enabled"])
}

// can we get the params from a handler's configuration?
func TestParamsFromConf(t *testing.T) {
	conf, err := Load("../config.toml")
	assert.NoError(t, err)

	assert.Equal(t, map[string]interface{}{"DbAddress": "REDACTED"}, conf.GetParams("postgresql"))
}
