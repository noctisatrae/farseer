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

	rpcConf := conf.Servers["rpc"]

	// always-the-same option test
	assert.Equal(t, HubParams{
		GossipPort:      2282,
		BootstrapPeers:  []string{"/dns/nemes.farcaster.xyz/tcp/2282/p2p/12D3KooWMQrf6unpGJfLBmTGy3eKTo4cGcXktWRbgMnfbZLXqBbn"},
		Debug:           false,
		BufferSize:      128,
		ContactInterval: 30,
	}, conf.Hub)

	// dynamic conf
	assert.Equal(t, true, rpcConf.(map[string]interface{})["enabled"])
}
