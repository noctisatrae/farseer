package main

import (
	handler "farseer/handlers"
	protos "farseer/protos"

	"github.com/charmbracelet/log"
)

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
	CastAddHandler:    CastAddHandler,
	CastRemoveHandler: CastRemoveHandler,
	// Initialize other handlers as needed...
}
