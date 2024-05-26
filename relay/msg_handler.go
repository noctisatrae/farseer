package main

import (
	protos "farseer/protos"

	"github.com/charmbracelet/log"
)

type HandlerBehaviour[T any] func(data *protos.MessageData) (T, error)

type Handler[T any] struct {
	CastAddHandler        HandlerBehaviour[T]
	CastRemoveHandler     HandlerBehaviour[T]
	FrameActionHandler    HandlerBehaviour[T]
	ReactionAddHandler    HandlerBehaviour[T]
	ReactionRemoveHandler HandlerBehaviour[T]
}

// T: What will be returned when the message is handled.
func (handler Handler[T]) handleMessages(messagesChan chan *protos.GossipMessage, ll log.Logger) {
	for messageBundle := range messagesChan {
		messages := messageBundle.GetMessageBundle().GetMessages()
		ll.Debug(messages)
		for m := range messages {
			data := messages[m].GetData()
			switch data.Type {
			case protos.MessageType_MESSAGE_TYPE_CAST_ADD:
				{
					if handler.CastAddHandler == nil {
						ll.Warn("CastAdd event was not handled because no function was provided!")
						return
					}
					result, err := handler.CastAddHandler(data)
					if err != nil {
						ll.Error("CastAdd handler encountered an error! |", "Error", err)
					}
					ll.Info("CastAdd event handled! |", "Result", result)
				}
			case protos.MessageType_MESSAGE_TYPE_CAST_REMOVE:
				{
					if handler.CastRemoveHandler == nil {
						ll.Warn("CastRemove event was not handled because no function was provided!")
						return
					}
					result, err := handler.CastRemoveHandler(data)
					if err != nil {
						ll.Error("CastRemove handler encountered an error! |", "Error", err)
					}
					ll.Info("CastRemove event handled! |", "Result", result)
				}
			case protos.MessageType_MESSAGE_TYPE_FRAME_ACTION:
				{
					if handler.FrameActionHandler == nil {
						ll.Warn("FrameAction event was not handled because no function was provided!")
						return
					}
					result, err := handler.FrameActionHandler(data)
					if err != nil {
						ll.Error("FrameAction handler encountered an error! |", "Error", err)
					}
					ll.Info("FrameAction event handled! |", "Result", result)
				}
			case protos.MessageType_MESSAGE_TYPE_REACTION_ADD:
				{
					if handler.ReactionAddHandler == nil {
						ll.Warn("ReactionAdd event was not handled because no function was provided!")
						return
					}
					result, err := handler.ReactionAddHandler(data)
					if err != nil {
						ll.Error("ReactionAdd handler encountered an error! |", "Error", err)
					}
					ll.Info("ReactionAdd event handled! |", "Result", result)
				}
			case protos.MessageType_MESSAGE_TYPE_REACTION_REMOVE:
				{
					if handler.ReactionRemoveHandler == nil {
						ll.Warn("ReactionRemove event was not handled because no function was provided!")
						return
					}
					result, err := handler.ReactionRemoveHandler(data)
					if err != nil {
						ll.Error("ReactionRemove handler encountered an error! |", "Error", err)
					}
					ll.Info("ReactionRemove event handled! |", "Result", result)
				}
			default:
				ll.Warn("Unhandled message type! |", "Type", data.Type)
			}
		}
	}
}
