package main

import (
	"github.com/pelletier/go-toml/v2"
	"os"
)

type HubParams struct {
	GossipPort      uint
	BootstrapPeers  []string
	Debug           bool
	BufferSize      uint
	ContactInterval uint
}

type Config struct {
	Hub       HubParams
	Servers   map[string]interface{} `toml:"servers"`
	Databases map[string]interface{} `toml:"databases"`
}

func Load(path string) (Config, error) {
	var config Config

	fileByte, err := os.ReadFile(path)
	if err != nil {
		return Config{
			Hub: HubParams{
				GossipPort:      2282,
				BootstrapPeers:  []string{"/dns/nemes.farcaster.xyz/tcp/2282/p2p/12D3KooWMQrf6unpGJfLBmTGy3eKTo4cGcXktWRbgMnfbZLXqBbn"},
				Debug:           false,
				BufferSize:      128,
				ContactInterval: 30,
			},
		}, err
	}

	err = toml.Unmarshal(fileByte, &config)
	if err != nil {
		return Config{
			Hub: HubParams{
				GossipPort:      2282,
				BootstrapPeers:  []string{"/dns/nemes.farcaster.xyz/tcp/2282/p2p/12D3KooWMQrf6unpGJfLBmTGy3eKTo4cGcXktWRbgMnfbZLXqBbn"},
				Debug:           false,
				BufferSize:      128,
				ContactInterval: 30,
			},
		}, err
	}

	var serverData map[string]interface{}
	if err = toml.Unmarshal(fileByte, &serverData); err == nil {
		if servers, ok := serverData["servers"].(map[string]interface{}); ok {
			config.Servers = servers
		}
	}

	var databasesData map[string]interface{}
	if err = toml.Unmarshal(fileByte, &serverData); err == nil {
		if databases, ok := databasesData["databases"].(map[string]interface{}); ok {
			config.Databases = databases
		}
	}

	return config, nil
}
