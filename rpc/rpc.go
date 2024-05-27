package main

import (
	handler "farseer/handlers"
	protos "farseer/protos"
	"fmt"
)

func CastAddHandler(data *protos.MessageData) error {
	fmt.Printf("Cast added: %s", data.GetCastAddBody())
	return nil
}

func CastRemoveHandler(data *protos.MessageData) error {
	fmt.Printf("Cast removed: %s", data.GetCastRemoveBody())
	return nil
}

// Exported variable
var PluginHandler = handler.Handler{
	CastAddHandler:        CastAddHandler,
	CastRemoveHandler:     CastRemoveHandler,
	// Initialize other handlers as needed...
}
