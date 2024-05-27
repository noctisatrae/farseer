package main

import (
	"os"
	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	GossipPort uint
	BootstrapPeers []string
	Debug bool
	BufferSize uint
	ContactInterval uint
}

func Load(path string) (Config, error) {
	var config Config
	
	fileByte, err := os.ReadFile(path)
	if err != nil {
		return Config{
			GossipPort: 2282,
			BootstrapPeers: []string{"/dns/nemes.farcaster.xyz/tcp/2282/p2p/12D3KooWMQrf6unpGJfLBmTGy3eKTo4cGcXktWRbgMnfbZLXqBbn"},
			Debug: false,
			BufferSize: 128,
			ContactInterval: 30,	
		}, err
	}

	err = toml.Unmarshal(fileByte, &config)
	if err != nil {
		return Config{
			GossipPort: 2282,
			BootstrapPeers: []string{"/dns/nemes.farcaster.xyz/tcp/2282/p2p/12D3KooWMQrf6unpGJfLBmTGy3eKTo4cGcXktWRbgMnfbZLXqBbn"},
			Debug: false,
			BufferSize: 128,
			ContactInterval: 30,	
		}, err
	}

	return config, nil
}