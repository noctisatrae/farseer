package main

import (
	"errors"
	"farseer/config"
	protos "farseer/protos"
	"farseer/time"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"
)

func TestInitConn(t *testing.T) {
	err := InitBehaviour(map[string]interface{}{
		"DbAddress": "postgres://postgres:example@localhost:5432/postgres",
	})

	assert.NoError(t, err)
}

func TestParamsFromConf(t *testing.T) {
	conf, err := config.Load("../config.toml")
	assert.NoError(t, err)

	err = InitBehaviour(conf.GetParams("postgresql"))
	assert.NoError(t, err)
}

func TestParamsCheck(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	conf, err := config.Load("../config.toml")
	assert.NoError(t, err)

	params := conf.GetParams("postgresql")

	currentFcTime, err := time.GetFarcasterTime()
	assert.NoError(t, err)

	msgData := protos.MessageData{
		Type:      protos.MessageType_MESSAGE_TYPE_CAST_ADD,
		Fid:       2,
		Timestamp: uint32(currentFcTime),
		Network:   protos.FarcasterNetwork_FARCASTER_NETWORK_MAINNET,
		Body:      &protos.MessageData_CastAddBody{},
	}

	msgFilter := params["MessageTypesAllowed"]
	log.Debug(msgFilter, "MsgType", msgData.Type.Number())

	err = CheckConfigParams(&msgData, params, func(data *protos.MessageData, params map[string]interface{}) error {
		return errors.New("if this error is raised, the test pass")
	})

	assert.Error(t, err)
}
