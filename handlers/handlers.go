package handlers

import (
	protos "github.com/noctisatrae/farseer/protos"
	"github.com/charmbracelet/log"
)

type InitBehaviour func(params map[string]interface{}) error
// This is the definition of the type of function that will handle incoming messages.
type HandlerBehaviour func(data *protos.MessageData, hash []byte, params map[string]interface{}) error

type Handler struct {
	// The name of the plugin/custom handler - used for contextualization in logs.
	Name string
	// The function that is run at the start of the handling of messages. If you need to create an instance of a DB,
	// you need to mutate the params map passed to the functions without overriding something you might need.
	// Here's a little example:
	// func InitBehaviour(params map[string]interface{}) error {
	// 	params["dbConn"] = conn.New()
	// 	return nil
	// }
	// Obviously, it's over-simplified without all the error handling & log reporting... I would also advise you to make a
	// rough scheme of the look of your params map. Because it's an interface and not a strongly typed struct, it can lead to
	// panicking if things aren't well queried.
	// You also might want some checks to see if the info you need from the config are there!
	InitHandler               InitBehaviour
	CastAddHandler            HandlerBehaviour
	CastRemoveHandler         HandlerBehaviour
	FrameActionHandler        HandlerBehaviour
	ReactionAddHandler        HandlerBehaviour
	ReactionRemoveHandler     HandlerBehaviour
	LinkAddHandler            HandlerBehaviour
	LinkRemoveHandler         HandlerBehaviour
	VerificationAddHandler    HandlerBehaviour
	VerificationRemoveHandler HandlerBehaviour
}

func (handler Handler) HandleMessages(messages chan *protos.GossipMessage, ll log.Logger, params map[string]interface{}) {
	if handler.InitHandler == nil {
	} else {
		err := handler.InitHandler(params)
		if err != nil {
			ll.Error("A handler encountered a problem! |", "Name", handler.Name, "Error", err)
		}
	}
	for msgB := range messages { // i hope that the chan only gives one message at a time so it's just O(n) and not O(n²)
		for _, m := range msgB.GetMessageBundle().GetMessages() {
			data := m.Data
			hash := m.Hash
			switch data.Type {
			case protos.MessageType_MESSAGE_TYPE_CAST_ADD:
				if handler.CastAddHandler == nil {
					ll.Info("New cast published! |", "Body", data)
				} else {
					err := handler.CastAddHandler(data, hash, params)
					if err != nil {
						ll.Error("CastAdd handler encountered an error! |", "Error", err)
					}
				}
			case protos.MessageType_MESSAGE_TYPE_CAST_REMOVE:
				if handler.CastRemoveHandler == nil {
					ll.Info("Cast was just removed! |", "Body", data)
				} else {
					err := handler.CastRemoveHandler(data, hash, params)
					if err != nil {
						ll.Error("CastRemove handler encountered an error! |", "Error", err)
					}
				}
			case protos.MessageType_MESSAGE_TYPE_FRAME_ACTION:
				if handler.FrameActionHandler == nil {
					ll.Info("New frame interaction! |", "Action", data)
				} else {
					err := handler.FrameActionHandler(data, hash, params)
					if err != nil {
						ll.Error("FrameAction handler encountered an error! |", "Error", err)
					}
				}
			case protos.MessageType_MESSAGE_TYPE_REACTION_ADD:
				if handler.ReactionAddHandler == nil {
					ll.Info("New reaction added! |", "Reaction", data)
				} else {
					err := handler.ReactionAddHandler(data, hash, params)
					if err != nil {
						ll.Error("ReactionAdd handler encountered an error! |", "Error", err)
					}
				}
			case protos.MessageType_MESSAGE_TYPE_REACTION_REMOVE:
				if handler.ReactionRemoveHandler == nil {
					ll.Info("A reaction was removed! |", "Reaction", data)
				} else {
					err := handler.ReactionRemoveHandler(data, hash, params)
					if err != nil {
						ll.Error("ReactionRemove handler encountered an error! |", "Error", err)
					}
				}
			case protos.MessageType_MESSAGE_TYPE_LINK_ADD:
				if handler.LinkAddHandler == nil {
					ll.Info("A link was added! |", "Link", data)
				} else {
					err := handler.LinkAddHandler(data, hash, params)
					if err != nil {
						ll.Error("LinkAdd handler encountered an error! |", "Error", err)
					}
				}
			case protos.MessageType_MESSAGE_TYPE_LINK_REMOVE:
				if handler.LinkRemoveHandler == nil {
					ll.Info("A link was removed! |", "Link", data)
				} else {
					err := handler.LinkAddHandler(data, hash, params)
					if err != nil {
						ll.Error("LinkRemove handler encountered an error! |", "Error", err)
					}
				}
			case protos.MessageType_MESSAGE_TYPE_VERIFICATION_ADD_ETH_ADDRESS:
				if handler.VerificationAddHandler == nil {
					ll.Info("A ETH address was just verified! |", "VerificationBody", data)
				} else {
					err := handler.VerificationAddHandler(data, hash, params)
					if err != nil {
						ll.Error("VerificationAdd handler encountered an error! |", "Error", err)
					}
				}
			case protos.MessageType_MESSAGE_TYPE_VERIFICATION_REMOVE:
				if handler.VerificationRemoveHandler == nil {
					ll.Info("A ETH address was just removed! |", "VerificationBody", data)
				} else {
					err := handler.VerificationAddHandler(data, hash, params)
					if err != nil {
						ll.Error("VerificationRemove handler encountered an error! |", "Error", err)
					}
				}
			default:
				ll.Warn("Unhandled message type! |", "Type", data.Type)
			}
		}
	}
}
