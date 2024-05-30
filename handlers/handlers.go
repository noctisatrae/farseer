package handlers

import (
	protos "farseer/protos"

	"github.com/charmbracelet/log"
)

type HandlerBehaviour func(data *protos.MessageData, params map[string]interface{}) error

type Handler struct {
	CastAddHandler        HandlerBehaviour
	CastRemoveHandler     HandlerBehaviour
	FrameActionHandler    HandlerBehaviour
	ReactionAddHandler    HandlerBehaviour
	ReactionRemoveHandler HandlerBehaviour
	LinkAddHandler        HandlerBehaviour
	LinkRemoveHandler     HandlerBehaviour
	Params                map[string]interface{}
}

func (handler Handler) HandleMessages(messages chan *protos.GossipMessage, ll log.Logger) {
	for msgB := range messages { // i hope that the message only gives one message at a time so it's just O(n) and not O(n²)
		for _, m := range msgB.GetMessageBundle().GetMessages() {
			data := m.Data
			switch data.Type {
			case protos.MessageType_MESSAGE_TYPE_CAST_ADD:
				if handler.CastAddHandler == nil {
					ll.Info("New cast published! |", "Body", data.GetCastAddBody())
				} else {
					err := handler.CastAddHandler(data, handler.Params)
					if err != nil {
						ll.Error("CastAdd handler encountered an error! |", "Error", err)
					}
				}
			case protos.MessageType_MESSAGE_TYPE_CAST_REMOVE:
				if handler.CastRemoveHandler == nil {
					ll.Info("Cast was just removed! |", "Body", data.GetCastRemoveBody())
				} else {
					err := handler.CastRemoveHandler(data, handler.Params)
					if err != nil {
						ll.Error("CastRemove handler encountered an error! |", "Error", err)
					}
				}
			case protos.MessageType_MESSAGE_TYPE_FRAME_ACTION:
				if handler.FrameActionHandler == nil {
					ll.Info("New frame interaction! |", "Action", data.GetFrameActionBody())
				} else {
					err := handler.FrameActionHandler(data, handler.Params)
					if err != nil {
						ll.Error("FrameAction handler encountered an error! |", "Error", err)
					}
				}
			case protos.MessageType_MESSAGE_TYPE_REACTION_ADD:
				if handler.ReactionAddHandler == nil {
					ll.Info("New reaction added! |", "Reaction", data.GetReactionBody())
				} else {
					err := handler.ReactionAddHandler(data, handler.Params)
					if err != nil {
						ll.Error("ReactionAdd handler encountered an error! |", "Error", err)
					}
				}
			case protos.MessageType_MESSAGE_TYPE_REACTION_REMOVE:
				if handler.ReactionRemoveHandler == nil {
					ll.Info("A reaction was removed! |", "Reaction", data.GetReactionBody())
				} else {
					err := handler.ReactionRemoveHandler(data, handler.Params)
					if err != nil {
						ll.Error("ReactionRemove handler encountered an error! |", "Error", err)
					}
				}
			case protos.MessageType_MESSAGE_TYPE_LINK_ADD:
				if handler.LinkAddHandler == nil {
					ll.Info("A link was added! |", "Link", data.GetLinkBody())
				} else {
					err := handler.LinkAddHandler(data, handler.Params)
					if err != nil {
						ll.Error("LinkAdd handler encountered an error! |", "Error", err)
					}
				}
			case protos.MessageType_MESSAGE_TYPE_LINK_REMOVE:
				if handler.LinkRemoveHandler == nil {
					ll.Info("A link was removed! |", "Link", data.GetLinkBody())
				} else {
					err := handler.LinkAddHandler(data, handler.Params)
					if err != nil {
						ll.Error("LinkRemove handler encountered an error! |", "Error", err)
					}
				}
			default:
				ll.Warn("Unhandled message type! |", "Type", data.Type)
			}
		}
	}
}
