package main

import (
	"context"
	"errors"

	handler "farseer/handlers"
	protos "farseer/protos"
	utils "farseer/utils"

	"github.com/jackc/pgx/v5"
)

const (
	sqlCastAdd = `
	INSERT INTO all_messages (
		fid, 
		timestamp, 
		network, 
		type, 
		cast_add_text, 
		cast_add_parent_cast_id_fid, 
		cast_add_parent_cast_id_hash, 
		cast_add_parent_url, 
		cast_add_embeds_deprecated, 
		cast_add_mentions, 
		cast_add_mentions_positions, 
		cast_add_embeds
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	sqlCastRemove = `
	INSERT INTO cast_remove_messages (message_id, target_hash)
	VALUES ($1, $2)
	`
)

func CheckConfigParams(data *protos.MessageData, params map[string]interface{}, handlerFunc handler.HandlerBehaviour) error {
	msgFilter := params["MessageTypesAllowed"]
	fidFilter := params["FidsAllowed"]

	if msgFilter != nil && fidFilter != nil {
		allowedMsgType := utils.IntersectionOfArrays(msgFilter.([]interface{}), []interface{}{utils.MsgTypeToInt(data.Type)})
		allowedFids := utils.IntersectionOfArrays(fidFilter.([]interface{}), []interface{}{int64(data.Fid)})
		if len(allowedMsgType) > 0 && len(allowedFids) > 0 {
			return handlerFunc(data, params)
		}

		return nil
	} else if msgFilter == nil && fidFilter != nil {
		allowedFids := utils.IntersectionOfArrays(fidFilter.([]interface{}), []interface{}{int64(data.Fid)})
		if len(allowedFids) > 0 {
			return handlerFunc(data, params)
		}
	} else if msgFilter != nil && fidFilter == nil {
		allowedMsgType := utils.IntersectionOfArrays(msgFilter.([]interface{}), []interface{}{utils.MsgTypeToInt(data.Type)})
		if len(allowedMsgType) > 0 {
			return handlerFunc(data, params)
		}
	} else {
		return handlerFunc(data, params)
	}

	return nil
}

// InitBehaviour initializes the plugin by setting up a database connection.
//
// Example config.toml for this plugin:
//
// ```toml
// [handlers.postgresql]
// DbAddress = "your-database-address" # required
// ```
//
// Parameters:
// - params: A map containing configuration parameters.
func InitBehaviour(params map[string]interface{}) error {
	dbAddr := params["DbAddress"]
	if dbAddr == nil {
		return errors.New("no DbAddress was provided, so no connection can be made to the DB")
	}

	hdlCtx := context.Background()
	conn, err := pgx.Connect(hdlCtx, dbAddr.(string))
	if err != nil {
		return err
	}

	params["hdlCtx"] = hdlCtx
	params["dbConn"] = conn

	return nil
}

func CastAddHandler(data *protos.MessageData, params map[string]interface{}) error {
	hdlCtx := params["hdlCtx"].(context.Context)
	conn := params["dbConn"].(*pgx.Conn)

	castAddBody := data.GetCastAddBody()

	_, err := conn.Exec(hdlCtx, sqlCastAdd,
		data.Fid,
		data.Timestamp,
		data.Network.String(),
		data.Type.String(),
		castAddBody.Text,
		castAddBody.GetParentCastId().GetFid(),
		castAddBody.GetParentCastId().GetHash(),
		castAddBody.GetParentUrl(),
		castAddBody.EmbedsDeprecated,
		castAddBody.Mentions,
		castAddBody.MentionsPositions,
		castAddBody.Embeds,
	)
	if err != nil {
		return err
	}

	return nil
}

// Exported variable
var PluginHandler = handler.Handler{
	Name:        "PostgreSQL",
	InitHandler: InitBehaviour,
	CastAddHandler: func(data *protos.MessageData, params map[string]interface{}) error {
		return CheckConfigParams(data, params, CastAddHandler)
	},
}
