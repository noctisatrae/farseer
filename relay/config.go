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
	Hub      HubParams
	Handlers map[string]interface{} `toml:"handlers"`
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

	return config, nil
}

func (conf Config) GetHandlers() []string {
	keys := []string{}
	for k := range conf.Handlers {
		isKEnabled := conf.Handlers[k].(map[string]interface{})["Enabled"]
		if isKEnabled == true {
			keys = append(keys, k)
		} else if isKEnabled == nil {
			return keys
		}
	}
	return keys
}

func (conf Config) GetParams(handler string) map[string]interface{} {
	// Check if the handler exists and has a "params" field
	if handlerParams, ok := conf.Handlers[handler]; ok && handlerParams != nil {
		if params, ok := handlerParams.(map[string]interface{})["Params"].(map[string]interface{}); ok {
			return params
		}
	}

	// Return an empty map if the handler does not exist or does not have a "params" field
	return map[string]interface{}{}
}
