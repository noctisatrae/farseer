package main

import (
	"context"
	"errors"
	handler "farseer/handlers"
	protos "farseer/protos"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5"
)

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

	conn, err := pgx.Connect(context.Background(), dbAddr.(string))
	if err != nil {
		return err
	}

	params["dbConn"] = conn

	return nil
}

func CastAddHandler(data *protos.MessageData, params map[string]interface{}) error {
	log.Debug("PLUGIN LOADED!", "Params", params, "Data", data)
	return nil
}

func CastRemoveHandler(data *protos.MessageData, params map[string]interface{}) error {
	log.Debug("PLUGIN LOADED!", "Params", params, "Data", data)
	return nil
}

// Exported variable
var PluginHandler = handler.Handler{
	Name:              "PostgreSQL",
	InitHandler:       InitBehaviour,
	CastAddHandler:    CastAddHandler,
	CastRemoveHandler: CastRemoveHandler,
}
