package main

import (
	protos "farseer/protos"
	"os"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"
)

func TestIndividualLoader(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	ll := *log.NewWithOptions(os.Stderr, log.Options{})

	messages := make(chan *protos.GossipMessage)
	err := LoadHandler("rpc", messages, ll)
	assert.NoError(t, err)
}