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
	// A query that allows to see if a cast exist in the DB
	CastCheck = `
	SELECT created_at FROM casts WHERE hash = $1 
	`

	// A query that allows to add a cast into the DB
	CastAdd = `
	INSERT INTO casts (
		fid,

		timestamp,

		hash,
		parent_hash,
		parent_fid,
		parent_url,

		text,
		embeds,
		mentions,
		mentions_positions
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	CastAddRemoved = `
	INSERT INTO casts (
		fid, 

		timestamp,
		deleted_at,

		hash
	) VALUES ($1, $2, CURRENT_TIMESTAMP, $3)
	`

	// A query that updates a cast when it's removed
	UpdateCastOnRemove = `
	UPDATE casts
	SET
		deleted_at = CURRENT_TIMESTAMP,
		updated_at = CURRENT_TIMESTAMP
	WHERE
		hash = $1
	`

	// A query to add a link into the database
	LinkAdd = `
	INSERT INTO links (
		timestamp,

		fid,
		target_fid,
		hash,
		type
	) VALUES ($1, $2, $3, $4, $5)
	`

	LinkAddRemoved = `
	INSERT INTO links (
		timestamp,
		deleted_at,

		fid,
		target_fid,
		hash,
		type
	) VALUES ($1, CURRENT_TIMESTAMP, $2, $3, $4, $5)
	`

	LinkRemove = `
	UPDATE links
	SET
		deleted_at = CURRENT_TIMESTAMP,
		updated_at = CURRENT_TIMESTAMP
	WHERE
		target_fid = $1
	`

	// Add a new reaction to the DB
	ReactionAdd = `
	INSERT INTO reactions (
		fid,

		timestamp,
		
		reaction_type,
		hash,
		target_hash,
		target_fid,
		target_url
	) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	ReactionAddRemoved = `
	INSERT INTO reactions (
		fid,

		timestamp,
		deleted_at,
		
		reaction_type,
		hash,
		target_hash,
		target_fid,
		target_url
	) VALUES ($1, $2, CURRENT_TIMESTAMP, $3, $4, $5, $6, $7)
	`

	ReactionRemove = `
	UPDATE reactions
	SET
		deleted_at = CURRENT_TIMESTAMP,
		updated_at = CURRENT_TIMESTAMP
	WHERE
		target_hash = $1
	`

	VerificationAdd = `
	INSERT INTO verifications (
		fid,

		timestamp,

		hash,
		address,
		claim_signature,
		block_hash,
		verification_type,
		chain_id,
		protocol
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
)

// // data.Timestamp => standard timestamp => to SQL compatible timestamp
// func fcTimeToSqlTime(fcTimestamp uint32) (time.Time, error) {
// 	normalTimestamp, err := fctime.FromFarcasterTime(int64(fcTimestamp))
// 	if err != nil {
// 		return time.Time{}, err
// 	}
// 	timestamp := time.Unix(normalTimestamp, 0)
// 	return timestamp, nil
// }

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

	hashStr := utils.BytesToHex(hash)

	castAddBody := data.GetCastAddBody()
	log.Debug("CastAddHandler, handling message", "Hash", hashStr)

	var parentHash []byte = []byte{}
	var parentFid uint64
	var parentUrl string = ""
	if castAddBody.GetParentCastId() != nil {
		parentHash = castAddBody.GetParentCastId().Hash
		parentFid = castAddBody.GetParentCastId().Fid
		parentUrl = castAddBody.GetParentUrl()
	}

	_, err := conn.Exec(hdlCtx, CastAdd,
		data.Fid,
		data.Timestamp,

		hashStr,
		// todo: nil check
		utils.BytesToHex(parentHash),
		parentFid,
		parentUrl,

		castAddBody.Text,
		castAddBody.Embeds,
		castAddBody.Mentions,
		castAddBody.MentionsPositions,
	)

	return err
}

func CastRemoveHandler(data *protos.MessageData, hash []byte, params map[string]interface{}) error {
	hdlCtx := params["hdlCtx"].(context.Context)
	conn := params["dbConn"].(*pgx.Conn)

	castHashToRemove := utils.BytesToHex(data.GetCastRemoveBody().TargetHash)

	cmdTag, err := conn.Exec(hdlCtx, UpdateCastOnRemove, castHashToRemove)
	if err != nil {
		return err
	} else if cmdTag.RowsAffected() == 0 {
		_, err = conn.Exec(
			hdlCtx,
			CastAddRemoved,
			data.Fid,

			data.Timestamp, // timestamp

			castHashToRemove,
		)
		return err
	}

	return err
}

func LinkAddHandler(data *protos.MessageData, hash []byte, params map[string]interface{}) error {
	hdlCtx := params["hdlCtx"].(context.Context)
	conn := params["dbConn"].(*pgx.Conn)

	LinkHash := utils.BytesToHex(hash)

	LinkAddBody := data.GetLinkBody()
	_, err := conn.Exec(hdlCtx, LinkAdd,
		data.Timestamp,

		data.Fid,
		LinkAddBody.GetTargetFid(),
		LinkHash,
		LinkAddBody.Type,
	)

	return err
}

func LinkRemoveHandler(data *protos.MessageData, hash []byte, params map[string]interface{}) error {
	hdlCtx := params["hdlCtx"].(context.Context)
	conn := params["dbConn"].(*pgx.Conn)

	LinkRemoveBody := data.GetLinkBody()
	LinkHash := utils.BytesToHex(hash)

	// TODO: Link check before remove query
	cmdTag, err := conn.Exec(hdlCtx, LinkRemove, LinkRemoveBody.GetTargetFid())
	if err != nil {
		return err
	} else if cmdTag.RowsAffected() == 0 {
		_, err = conn.Exec(hdlCtx, LinkAddRemoved,
			data.Timestamp,
			data.Fid,
			LinkRemoveBody.GetTargetFid(),
			LinkHash,
			LinkRemoveBody.Type,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func ReactionAddHandler(data *protos.MessageData, hash []byte, params map[string]interface{}) error {
	hdlCtx := params["hdlCtx"].(context.Context)
	conn := params["dbConn"].(*pgx.Conn)

	ReactionAddBody := data.GetReactionBody()
	_, err := conn.Exec(hdlCtx, ReactionAdd,
		data.Fid,
		data.Timestamp,
		ReactionAddBody.Type,
		utils.BytesToHex(hash),
		utils.BytesToHex(ReactionAddBody.GetTargetCastId().Hash),
		ReactionAddBody.GetTargetCastId().GetFid(),
		ReactionAddBody.GetTargetUrl(),
	)

	return err
}

func ReactionRemoveHandler(data *protos.MessageData, hash []byte, params map[string]interface{}) error {
	hdlCtx := params["hdlCtx"].(context.Context)
	conn := params["dbConn"].(*pgx.Conn)

	ReactionRemoveBody := data.GetReactionBody()
	cmdTag, err := conn.Exec(hdlCtx, ReactionRemove,
		utils.BytesToHex(ReactionRemoveBody.GetTargetCastId().GetHash()),
	)
	if err != nil {
		return err
	} else if cmdTag.RowsAffected() == 0 {
		_, err := conn.Exec(hdlCtx, ReactionAddRemoved,
			data.Fid,
			data.Timestamp,
			ReactionRemoveBody.Type,
			utils.BytesToHex(hash),
			utils.BytesToHex(ReactionRemoveBody.GetTargetCastId().GetHash()),
			ReactionRemoveBody.GetTargetCastId().Fid,
			ReactionRemoveBody.GetTargetUrl(),
		)
		if err != nil {
			return err
		}
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
	LinkAddHandler: func(data *protos.MessageData, hash []byte, params map[string]interface{}) error {
		return CheckConfigParams(data, params, hash, LinkAddHandler)
	},
	LinkRemoveHandler: func(data *protos.MessageData, hash []byte, params map[string]interface{}) error {
		return CheckConfigParams(data, params, hash, LinkRemoveHandler)
	},
	ReactionAddHandler: func(data *protos.MessageData, hash []byte, params map[string]interface{}) error {
		return CheckConfigParams(data, params, hash, ReactionAddHandler)
	},
	ReactionRemoveHandler: func(data *protos.MessageData, hash []byte, params map[string]interface{}) error {
		return CheckConfigParams(data, params, hash, ReactionRemoveHandler)
	},
}
