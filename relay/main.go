package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	protos "farseer/protos"

	"github.com/charmbracelet/log"

	// libp2p
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-mplex"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multiaddr-dns"

	// "github.com/libp2p/go-libp2p/core/host"
	// "github.com/libp2p/go-libp2p/core/peer"

	// INITIAL secret var
	"github.com/joho/godotenv"
)


func checkConnectionStatus(h host.Host, peerID peer.ID) {
	connected := h.Network().Connectedness(peerID)
	if connected == network.Connected {
		log.Info("Successfully connected to peer! |", "peerID", peerID)
	} else {
		log.Warn("Not connected to peer |", "peerID", peerID)
	}
}

func logMessages(messages chan *protos.GossipMessage, local log.Logger) {
	for msg := range messages {
		local.Info("Received message! |", "msg", msg)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx := context.Background()

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	gossipsubPort, err := strconv.Atoi(os.Getenv("GOSSIPSUB_PORT"))
	if err != nil {
		log.Fatal("Can't parse default gossipsub port, QUITING! |", "Error", err)
	}

	dns_resolver, err := madns.NewResolver()
	if err != nil {
		log.Fatal("Could not start the DNS resolver", "Error", err)
	}
	log.Info("Successfully started the DNS resolver!")

	h, err := libp2p.New(		
		libp2p.Ping(true),
		libp2p.ListenAddrStrings(
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", os.Getenv("GOSSIPSUB_PORT")),
		),
		libp2p.Muxer("/mplex/6.7.0", mplex.DefaultTransport),
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	
	log.Info("Started the libp2p host!", "Addr", fmt.Sprintf("%s/p2p/%s", h.Addrs()[1], h.ID().String()))

	init_peer, err := multiaddr.NewMultiaddr(os.Getenv("INIT_PEER"))
	if err != nil {
		log.Fatal("Couldn't parse multiaddr!", "Error", err)
	}

	resolved_multiaddrs, err := dns_resolver.Resolve(ctx, init_peer)
	if err != nil {
		log.Fatal("Can't resolve from DNS addr", "Error", err)
	}
	peer_addrinfo, err := peer.AddrInfoFromP2pAddr(resolved_multiaddrs[0])
	if err != nil {
		log.Fatal("Can't convert multiaddr to addrinfo", "Error", err)
	}

	log.Info("Connecting to a remote peer! |", "peer", peer_addrinfo)
	err = h.Connect(ctx, *peer_addrinfo)
	if err != nil {
		log.Error("", "Error", err)
	}

	checkConnectionStatus(h, peer_addrinfo.ID)

	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		log.Error(err)
	}

	netw_primary, err := ReceiveMessages(ctx, ps, h.ID(), "primary")
	if err != nil {
		log.Fatal(err.Error())
	}

	netw_contact, err := ReceiveMessages(ctx, ps, h.ID(), "contact_info")
	if err != nil {
		log.Fatal(err.Error())
	}

	netw_discovery, err := ReceiveMessages(ctx, ps, h.ID(), "peer_discovery")
	if err != nil {
		log.Fatal(err.Error())
	}

	go logMessages(netw_primary.NetworkMessage, netw_primary.logger)
	go logMessages(netw_contact.NetworkMessage, netw_contact.logger)
	go logMessages(netw_discovery.NetworkMessage, netw_discovery.logger)

	go func() {
		for range ticker.C {
			netw_contact.PublishContactInfo(&protos.ContactInfoContent{
				HubVersion: "2023",
				Network: 2,
				GossipAddress: &protos.GossipAddressInfo{
					Address: h.Addrs()[0].String(),
					Port: uint32(gossipsubPort),
				},
			})
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	log.Info("Received signal, shutting down...")

	// shut the node down
	if err := h.Close(); err != nil {
		log.Fatal(err.Error())
	}
}