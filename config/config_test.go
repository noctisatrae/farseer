package config_test

import (
	"testing"

	config "farseer/config"

	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	conf, err := config.Load("../config.toml")
	assert.NoError(t, err)

	postgreConf := conf.Handlers["postgresql"]

	// always-the-same option test
	assert.Equal(t, config.HubParams{
		GossipPort:      2282,
		BootstrapPeers:  []string{
			"/ip4/5.189.129.220/tcp/2282/p2p/12D3KooWASKSaR6jiXSxbr1F4r9BzAh3csFXZkFHSwoCD9k7WSti",
			"/ip4/5.161.126.48/tcp/2282/p2p/12D3KooWLuagvFHo6AWZ9cbQLsDmFNo8E3mJco7ixdu1CxTqjey1",
			"/ip4/144.91.76.58/tcp/2282/p2p/12D3KooWBzjv7Lp37U3qTaEDbcDhfeSQdsPS5aH4iEi5b4gQd2UQ",
		},
		Debug:           true,
		BufferSize:      128,
		ContactInterval: 30,
	}, conf.Hub)

	// dynamic conf
	assert.Equal(t, true, postgreConf.(map[string]interface{})["Enabled"])
}

// can we get the params from a handler's configuration?
func TestParamsFromConf(t *testing.T) {
	conf, err := config.Load("../config.toml")
	assert.NoError(t, err)

	assert.Equal(t, map[string]interface{}{"DbAddress": "postgres://postgres:example@localhost:5432/postgres", "MessageTypesAllowed": []interface{}{int64(1)}}, conf.GetParams("postgresql"))
}
