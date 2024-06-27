package main

import (
	"context"
	"errors"
	"farseer/config"
	protos "farseer/protos"
	FcTime "farseer/time"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5"
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

	currentFcTime, err := FcTime.GetFarcasterTime()
	assert.NoError(t, err)

	msgData := protos.MessageData{
		Type:      protos.MessageType_MESSAGE_TYPE_CAST_ADD,
		Fid:       10626,
		Timestamp: uint32(currentFcTime),
		Network:   protos.FarcasterNetwork_FARCASTER_NETWORK_MAINNET,
		Body:      &protos.MessageData_CastAddBody{},
	}

	msgFilter := params["MessageTypesAllowed"]
	log.Debug(msgFilter, "MsgType", msgData.Type.Number())

	err = CheckConfigParams(&msgData, params, []byte{}, func(data *protos.MessageData, hash []byte, params map[string]interface{}) error {
		return errors.New("if this error is raised, the test pass")
	})

	assert.Error(t, err)
}

func TestCastAddHandler(t *testing.T) {
	params := map[string]interface{}{
		"DbAddress": "postgres://postgres:example@localhost:5432/postgres",
	}

	err := InitBehaviour(params)
	assert.NoError(t, err)

	fcTime, err := FcTime.GetFarcasterTime()
	assert.NoError(t, err)

	err = CastAddHandler(&protos.MessageData{
		Type: protos.MessageType_MESSAGE_TYPE_CAST_ADD,
		Fid: 10126,
		Timestamp: uint32(fcTime),
		Network: protos.FarcasterNetwork_FARCASTER_NETWORK_MAINNET,
		Body: &protos.MessageData_CastAddBody{
			CastAddBody: &protos.CastAddBody{
				Text: "Hello",
				Parent: &protos.CastAddBody_ParentCastId{
					ParentCastId: &protos.CastId{
						Fid: 10246,
						Hash: []byte{},
					},
				},
			},
		},
	}, []byte{}, params)
	assert.NoError(t, err)
}

func TestCastRemoveHandler(t *testing.T) {
	params := map[string]interface{}{
		"DbAddress": "postgres://postgres:example@localhost:5432/postgres",
	}

	err := InitBehaviour(params)
	assert.NoError(t, err)

	hdlCtx := params["hdlCtx"].(context.Context)
	conn := params["dbConn"].(*pgx.Conn)

	_, err = conn.Exec(hdlCtx, UpdateCastOnRemove, 2233434445, "")
	assert.NoError(t, err)
}