package main

import (
	"context"
	"fmt"
	"os"

	"farseer/config"
	"farseer/time"
	protos "farseer/protos"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
	"google.golang.org/protobuf/proto"

	"github.com/charmbracelet/log"
)

type Network struct {
	NetworkMessage chan *protos.GossipMessage

	ctx   context.Context
	ps    *pubsub.PubSub
	topic *pubsub.Topic
	sub   *pubsub.Subscription

	logger log.Logger
	self   peer.ID
}

func (netw *Network) PublishContactInfo(contact *protos.ContactInfoContent) {
	peerIdEncoded, err := netw.self.Marshal()
	if err != nil {
		netw.logger.Error("An empty PeerId will be sent because we can't marshall the provided one. |", "Error", err)
		peerIdEncoded = []byte{}
	}

	contactInfoTime, err := time.GetFarcasterTime()
	if err != nil {
		contactInfoTime = 0
		netw.logger.Error("Couldn't get Farcaster time for the message!")
	}

	m := protos.GossipMessage{
		Topics:    []string{netw.topic.String()},
		PeerId:    peerIdEncoded,
		Version:   protos.GossipVersion_GOSSIP_VERSION_V1_1,
		Timestamp: uint32(contactInfoTime),
		Content: &protos.GossipMessage_ContactInfoContent{
			ContactInfoContent: contact,
		},
	}

	netw.logger.Info("Sending!")

	if err := netw.Publish(&m); err != nil {
		netw.logger.Error("Error publishing message! |", "Error", err)
	}
}

func (netw *Network) Publish(m *protos.GossipMessage) error {
	mEncoded, err := proto.Marshal(m)
	if err != nil {
		netw.logger.Error("Couldn't encode the gossip message! |", "Error", err)
	}

	err = netw.topic.Publish(netw.ctx, mEncoded)
	return err
}

func ReceiveMessages(ctx context.Context, ps *pubsub.PubSub, selfId peer.ID, topicReq string, conf config.Config) (*Network, error) {
	req := fmt.Sprint("f_network_1_", topicReq)
	log.Info("Suscribing to a new topic! |", "Topic", req)

	topic, err := ps.Join(req)
	if err != nil {
		log.Fatal(err.Error())
	}

	sub, err := topic.Subscribe()
	if err != nil {
		log.Fatal(err.Error())
	}

	ll := log.NewWithOptions(os.Stderr, log.Options{
		Prefix: topicReq,
	})

	if conf.Hub.Debug {
		ll.SetLevel(log.DebugLevel)
	}

	netw := &Network{
		ctx:            ctx,
		ps:             ps,
		topic:          topic,
		sub:            sub,
		NetworkMessage: make(chan *protos.GossipMessage, conf.Hub.BufferSize),
		self:           selfId,
		logger:         *ll,
	}

	go netw.readLoop()
	return netw, nil
}

func (netw *Network) readLoop() {
	for {
		msg, err := netw.sub.Next(netw.ctx)
		if err != nil {
			log.Error(err.Error())
			close(netw.NetworkMessage)
			return
		}

		// if message received is from me => don't care
		if msg.ReceivedFrom == netw.self {
			continue
		}

		netwMsg := new(protos.GossipMessage)
		err = proto.Unmarshal(msg.Data, netwMsg)
		if err != nil {
			log.Error("Could not parse the incoming message! |", "error", err)
			continue
		}
		netw.NetworkMessage <- netwMsg
	}
}
