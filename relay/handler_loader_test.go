package main

import (
	"os"
	"testing"

	"farseer/config"
	protos "farseer/protos"

	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"
)

// Can we load an individual handler?
func TestIndividualLoader(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	ll := *log.NewWithOptions(os.Stderr, log.Options{})

	conf, err := config.Load("../config.toml")
	assert.NoError(t, err)

	messages := make(chan *protos.GossipMessage)
	err = LoadHandler("postgresql", messages, ll, conf)
	assert.NoError(t, err)
}

// Can I get the list of the compiled handlers?
func TestListCompiledHandlers(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	handlerList, err := ListCompiledHandlers()
	assert.NoError(t, err)

	assert.Equal(t, []string{"postgresql"}, handlerList)
}

// What handlers do I get from conf?
func TestHandlersFromConf(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	conf, err := config.Load("../config.toml")
	assert.NoError(t, err)

	keyArr := conf.GetHandlers()

	assert.Equal(t, []string{"postgresql"}, keyArr)
}

// Verify if a handler is loaded if 1. it is enabled 2. it is compiled
func TestWhatWillBeLoaded(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	conf, err := config.Load("../config.toml")
	assert.NoError(t, err)

	compiledHandlers, err := ListCompiledHandlers()
	assert.NoError(t, err)

	whatWillBeLoaded := intesectionOfArrays(conf.GetHandlers(), compiledHandlers)

	assert.Equal(t, []string{"postgresql"}, whatWillBeLoaded)
}

// can we load all the handlers without errors?
func TestMultipleLoader(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	ll := *log.NewWithOptions(os.Stderr, log.Options{})

	messages := make(chan *protos.GossipMessage)

	conf, err := config.Load("../config.toml")
	assert.NoError(t, err)

	err = LoadHandlersFromConf(conf, messages, ll)
	assert.NoError(t, err)
}
