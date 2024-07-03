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
		RpcPort: 2283,
		BootstrapPeers: []string{},
		Debug:           false,
		BufferSize:      128,
		ContactInterval: 3000,
	}, conf.Hub)

	// dynamic conf
	assert.Equal(t, true, postgreConf.(map[string]interface{})["Enabled"])
}

// can we get the params from a handler's configuration?
func TestParamsFromConf(t *testing.T) {
	conf, err := config.Load("../config.toml")
	assert.NoError(t, err)

	assert.Equal(t, map[string]interface{}{"DbAddress": "postgres://postgres:example@db:5432/postgres", "FidsAllowed": []interface{}{int64(10626)}}, conf.GetParams("postgresql"))
}
