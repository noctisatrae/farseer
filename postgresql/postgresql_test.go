package main

import (
	"errors"
	"testing"

	"github.com/noctisatrae/farseer/config"
	protos "github.com/noctisatrae/farseer/protos"
	FcTime "github.com/noctisatrae/farseer/time"
	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"
)

func TestInitConn(t *testing.T) {
	err := InitBehaviour(map[string]interface{}{
		"DbAddress": "postgres://postgres:example@localhost:5432/postgres",
	})

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
		Type:      protos.MessageType_MESSAGE_TYPE_CAST_ADD,
		Fid:       10126,
		Timestamp: uint32(fcTime),
		Network:   protos.FarcasterNetwork_FARCASTER_NETWORK_MAINNET,
		Body: &protos.MessageData_CastAddBody{
			CastAddBody: &protos.CastAddBody{
				Text: "Hello",
			},
		},
	}, []byte{3, 4, 5, 6}, params)
	assert.NoError(t, err)
}

func TestCastRemoveHandler(t *testing.T) {
	params := map[string]interface{}{
		"DbAddress": "postgres://postgres:example@localhost:5432/postgres",
	}

	err := InitBehaviour(params)
	assert.NoError(t, err)

	fcTime, err := FcTime.GetFarcasterTime()
	assert.NoError(t, err)

	err = CastRemoveHandler(&protos.MessageData{
		Type:      protos.MessageType_MESSAGE_TYPE_CAST_REMOVE,
		Fid:       10126,
		Timestamp: uint32(fcTime),
		Network:   protos.FarcasterNetwork_FARCASTER_NETWORK_MAINNET,
		Body: &protos.MessageData_CastRemoveBody{
			CastRemoveBody: &protos.CastRemoveBody{
				TargetHash: []byte{3, 4, 5, 6},
			},
		},
	}, []byte{}, params)

	assert.NoError(t, err)
}

func TestLinkAdd(t *testing.T) {
	params := map[string]interface{}{
		"DbAddress": "postgres://postgres:example@localhost:5432/postgres",
	}

	err := InitBehaviour(params)
	assert.NoError(t, err)

	fcTime, err := FcTime.GetFarcasterTime()
	assert.NoError(t, err)

	err = LinkAddHandler(&protos.MessageData{
		Type:      protos.MessageType_MESSAGE_TYPE_LINK_ADD,
		Fid:       10126,
		Timestamp: uint32(fcTime),
		Network:   protos.FarcasterNetwork_FARCASTER_NETWORK_MAINNET,
		Body: &protos.MessageData_LinkBody{
			LinkBody: &protos.LinkBody{
				Type: "follow",
				Target: &protos.LinkBody_TargetFid{
					TargetFid: 10222,
				},
			},
		},
	}, []byte{2, 3, 4, 5, 6}, params)
	assert.NoError(t, err)
}

func TestLinkRemove(t *testing.T) {
	params := map[string]interface{}{
		"DbAddress": "postgres://postgres:example@localhost:5432/postgres",
	}

	err := InitBehaviour(params)
	assert.NoError(t, err)

	fcTime, err := FcTime.GetFarcasterTime()
	assert.NoError(t, err)

	err = LinkRemoveHandler(&protos.MessageData{
		Type:      protos.MessageType_MESSAGE_TYPE_LINK_REMOVE,
		Fid:       10126,
		Timestamp: uint32(fcTime),
		Network:   protos.FarcasterNetwork_FARCASTER_NETWORK_MAINNET,
		Body: &protos.MessageData_LinkBody{
			LinkBody: &protos.LinkBody{
				Type: "follow",
				Target: &protos.LinkBody_TargetFid{
					TargetFid: uint64(10222),
				},
			},
		},
	}, []byte{2, 3, 4, 5, 6}, params)

	assert.NoError(t, err)
}

func TestReactionAdd(t *testing.T) {
	params := map[string]interface{}{
		"DbAddress": "postgres://postgres:example@localhost:5432/postgres",
	}

	err := InitBehaviour(params)
	assert.NoError(t, err)

	fcTime, err := FcTime.GetFarcasterTime()
	assert.NoError(t, err)

	reactionTestData := &protos.MessageData{
		Type:      protos.MessageType(protos.ReactionType_REACTION_TYPE_LIKE),
		Fid:       10267,
		Timestamp: uint32(fcTime),
		Network:   protos.FarcasterNetwork_FARCASTER_NETWORK_MAINNET,
		Body: &protos.MessageData_ReactionBody{
			ReactionBody: &protos.ReactionBody{
				Type: protos.ReactionType_REACTION_TYPE_RECAST,
				Target: &protos.ReactionBody_TargetCastId{
					TargetCastId: &protos.CastId{
						Fid:  10423,
						Hash: []byte{1, 2, 3, 4, 5, 6},
					},
				},
			},
		},
	}
	err = ReactionAddHandler(reactionTestData, []byte{5, 6, 7}, params)
	assert.NoError(t, err)
}

func TestReactionRemove(t *testing.T) {
	params := map[string]interface{}{
		"DbAddress": "postgres://postgres:example@localhost:5432/postgres",
	}

	err := InitBehaviour(params)
	assert.NoError(t, err)

	fcTime, err := FcTime.GetFarcasterTime()
	assert.NoError(t, err)

	err = ReactionRemoveHandler(&protos.MessageData{
		Type:      protos.MessageType_MESSAGE_TYPE_REACTION_REMOVE,
		Fid:       10267,
		Timestamp: uint32(fcTime),
		Network:   protos.FarcasterNetwork_FARCASTER_NETWORK_MAINNET,
		Body: &protos.MessageData_ReactionBody{
			ReactionBody: &protos.ReactionBody{
				Type: protos.ReactionType_REACTION_TYPE_RECAST,
				Target: &protos.ReactionBody_TargetCastId{
					TargetCastId: &protos.CastId{
						Fid:  10423,
						Hash: []byte{1, 2, 3, 4, 5, 6},
					},
				},
			},
		},
	}, []byte{0, 1, 0, 1}, params)
	assert.NoError(t, err)
}
