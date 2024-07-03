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

	postgreConf := conf.Handlers["postgresql"]

	// always-the-same option test
	assert.Equal(t, config.HubParams{
		PublicHubIp: "92.158.95.48",
		GossipPort:  2282,
		BootstrapPeers: []string{
			"/dns/nemes.farcaster.xyz/tcp/2283/p2p/12D3KooWMQrf6unpGJfLBmTGy3eKTo4cGcXktWRbgMnfbZLXqBbn",
		},
		Debug:           false,
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

	assert.Equal(t, map[string]interface{}{"DbAddress": "postgres://postgres:example@localhost:5432/postgres", "MessageTypesAllowed": []interface{}{int64(1), int64(2)}, "FidsAllowed": []interface{}{int64(10626)}}, conf.GetParams("postgresql"))
}
