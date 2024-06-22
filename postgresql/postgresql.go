package main

import (
	"context"
	"errors"

	handler "farseer/handlers"
	protos "farseer/protos"
	utils "farseer/utils"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5"
)

const (
	sqlCastAdd = `
	INSERT INTO all_messages (
		fid, 
		timestamp, 
		network, 
		type,
		hash,
		cast_add_text,
		cast_add_parent_cast_id_fid, 
		cast_add_parent_cast_id_hash, 
		cast_add_parent_url, 
		cast_add_embeds_deprecated, 
		cast_add_mentions, 
		cast_add_mentions_positions, 
		cast_add_embeds
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	sqlCastAddRemoved = `
	INSERT INTO all_messages (
		fid, 
		timestamp,
		network,
		type,
		hash,
		removed_at
	) VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
	`

	sqlCastUpdateRemoved = `
UPDATE all_messages
SET 
    removed_at = CURRENT_TIMESTAMP
WHERE 
    hash = $1;
	`
)

func CheckConfigParams(data *protos.MessageData, params map[string]interface{}, hash []byte, handlerFunc handler.HandlerBehaviour) error {
	msgFilter := params["MessageTypesAllowed"]
	fidFilter := params["FidsAllowed"]

	if msgFilter != nil && fidFilter != nil {
		allowedMsgType := utils.IntersectionOfArrays(msgFilter.([]interface{}), []interface{}{utils.MsgTypeToInt(data.Type)})
		allowedFids := utils.IntersectionOfArrays(fidFilter.([]interface{}), []interface{}{int64(data.Fid)})
		if len(allowedMsgType) > 0 && len(allowedFids) > 0 {
			return handlerFunc(data, hash, params)
		}

		return nil
	} else if msgFilter == nil && fidFilter != nil {
		allowedFids := utils.IntersectionOfArrays(fidFilter.([]interface{}), []interface{}{int64(data.Fid)})
		if len(allowedFids) > 0 {
			return handlerFunc(data, hash, params)
		}
	} else if msgFilter != nil && fidFilter == nil {
		allowedMsgType := utils.IntersectionOfArrays(msgFilter.([]interface{}), []interface{}{utils.MsgTypeToInt(data.Type)})
		if len(allowedMsgType) > 0 {
			return handlerFunc(data, hash, params)
		}
	} else {
		return handlerFunc(data, hash, params)
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

func CastAddHandler(data *protos.MessageData, hash []byte, params map[string]interface{}) error {
	hdlCtx := params["hdlCtx"].(context.Context)
	conn := params["dbConn"].(*pgx.Conn)

	castAddBody := data.GetCastAddBody()
	log.Debug("CastAddHandler, handling message", "Hash", utils.BytesToHex(hash))

	_, err := conn.Exec(hdlCtx, sqlCastAdd,
		data.Fid,
		data.Timestamp,
		data.Network.Number(),
		data.Type.Number(),
		utils.BytesToHex(hash),	
		castAddBody.Text,
		castAddBody.GetParentCastId().GetFid(),
		utils.BytesToHex(castAddBody.GetParentCastId().GetHash()),
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

func CastRemoveHandler(data *protos.MessageData, hash []byte, params map[string]interface{}) error {
	hdlCtx := params["hdlCtx"].(context.Context)
	conn := params["dbConn"].(*pgx.Conn)

	var id int
	castIdToRemove := utils.BytesToHex(data.GetCastRemoveBody().TargetHash) 

	err := conn.QueryRow(hdlCtx, sqlCastUpdateRemoved, castIdToRemove).Scan(&id)
	if err == pgx.ErrNoRows {
		_, err := conn.Exec(
			hdlCtx, 
			sqlCastAddRemoved, 
			// args
			data.Fid,
			data.Timestamp,
			data.Network.Number(),
			data.Type.Number(),
			castIdToRemove,
			data.Timestamp,
		)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

// Exported variable
var PluginHandler = handler.Handler{
	Name:        "PostgreSQL",
	InitHandler: InitBehaviour,
	CastAddHandler: func(data *protos.MessageData, hash []byte, params map[string]interface{}) error {
		return CheckConfigParams(data, params, hash, CastAddHandler)
	},
	CastRemoveHandler: func(data *protos.MessageData, hash []byte, params map[string]interface{}) error {
		return CheckConfigParams(data, params, hash, CastRemoveHandler)
	},
}
