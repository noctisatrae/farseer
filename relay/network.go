package main

import (
	"context"
	"fmt"
	"os"
	"time"

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

	m := protos.GossipMessage{
		Topics:    []string{netw.topic.String()},
		PeerId:    peerIdEncoded,
		Version:   protos.GossipVersion_GOSSIP_VERSION_V1_1,
		Timestamp: uint32(time.Now().Unix()),
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

func ReceiveMessages(ctx context.Context, ps *pubsub.PubSub, selfId peer.ID, topicReq string) (*Network, error) {
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

	conf, err := Load("./relay/config.toml")
	if err != nil {
		log.Error("Couldn't parse config file! |", "Error", err)
	}

	ll := log.NewWithOptions(os.Stderr, log.Options{
		Prefix: topicReq,
	})

	if conf.Debug {
		ll.SetLevel(log.DebugLevel)
	}

	netw := &Network{
		ctx:            ctx,
		ps:             ps,
		topic:          topic,
		sub:            sub,
		NetworkMessage: make(chan *protos.GossipMessage, conf.BufferSize),
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
