package main

import (
	protos "farseer/protos"

	"github.com/charmbracelet/log"
)

type HandlerBehaviour[T any] func (data *protos.MessageData) (T, error) 

type Handler[T any] struct {
	CastAddHandler HandlerBehaviour[T]
	CastRemoveHandler HandlerBehaviour[T]
	FrameActionHandler HandlerBehaviour[T]
	ReactionAddHandler HandlerBehaviour[T]
	ReactionRemoveHandler HandlerBehaviour[T]
}

// T: What will be returned when the message is handled.
func (handler Handler[T]) handleMessages(messages chan *protos.GossipMessage, ll log.Logger) {
	for m := range messages {
		data := m.GetMessage().GetData()
		switch data.Type {
		case protos.MessageType_MESSAGE_TYPE_CAST_ADD:
			{
				if handler.CastAddHandler == nil {
					ll.Warn("CastAdd event was not handled because no function was provided!")
					return
				}
				result, err := handler.CastAddHandler(data)
				if err != nil {
					ll.Info("CastAdd event handled! |", "Result", result)
				}
			}
		case protos.MessageType_MESSAGE_TYPE_CAST_REMOVE:
			{
				if handler.CastRemoveHandler == nil {
					ll.Warn("CastRemove event was not handled because no function was provided!")
					return
				}
				result, err := handler.CastRemoveHandler(data)
				if err != nil {
					ll.Info("CastRemove event handled! |", "Result", result)
				}
			}
		case protos.MessageType_MESSAGE_TYPE_FRAME_ACTION:
			{
				if handler.FrameActionHandler == nil {
					ll.Warn("FrameAction event was not handled because no function was provided!")
					return
				}
				result, err := handler.FrameActionHandler(data)
				if err != nil {
					ll.Info("FrameAction event handled! |", "Result", result)
				}
			}
		case protos.MessageType_MESSAGE_TYPE_REACTION_ADD:
			{
				if handler.ReactionAddHandler == nil {
					ll.Warn("ReactionAdd event was not handled because no function was provided!")
					return
				}
				result, err := handler.ReactionAddHandler(data)
				if err != nil {
					ll.Info("ReactionAdd event handled! |", "Result", result)
				}
			}
		case protos.MessageType_MESSAGE_TYPE_REACTION_REMOVE:
			{
				if handler.ReactionRemoveHandler == nil {
					ll.Warn("ReactionRemove event was not handled because no function was provided!")
					return
				}
				result, err := handler.ReactionRemoveHandler(data)
				if err != nil {
					ll.Info("ReactionRemove event handled! |", "Result", result)
				}
			}
		}
	}
}