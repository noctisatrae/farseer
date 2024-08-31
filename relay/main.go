package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	protos "github.com/noctisatrae/farseer/protos"
	"github.com/noctisatrae/farseer/config"

	"github.com/charmbracelet/log"

	// libp2p
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-mplex"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/security/noise"

	"github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multiaddr-dns"
)

const HUB_VERSION = "2024.7.24"

type ResolveResult struct {
	ResolvedMultiaddrs []multiaddr.Multiaddr
	Error              error
}

func checkConnectionStatus(h host.Host, peerID peer.ID) {
	connected := h.Network().Connectedness(peerID)
	if connected == network.Connected {
		log.Info("Successfully connected to peer! |", "peerID", peerID)
	} else {
		log.Warn("Not connected to peer |", "peerID", peerID)
	}
}

func logMessages(messages chan *protos.GossipMessage, ll log.Logger) {
	for msg := range messages {
		ll.Info("RECEIVED |", "msg", msg)
	}
}

func getId() (crypto.PrivKey, error) {
	if _, err := os.Stat("./hub_identity"); errors.Is(err, os.ErrNotExist) {
		log.Debug("Privkey file do not exist, creating it!")
		priv, _, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
		if err != nil {
			return nil, err
		}
		privBytes, err := crypto.MarshalPrivateKey(priv)
		if err != nil {
			return nil, err
		}

		err = os.WriteFile("./hub_identity", privBytes, 0644)
		if err != nil {
			return nil, err
		}

		return priv, nil
	} else {
		privBytes, err := os.ReadFile("./hub_identity")
		if err != nil {
			return nil, err
		}

		priv, err := crypto.UnmarshalPrivateKey(privBytes)
		if err != nil {
			return nil, err
		}

		return priv, nil
	}
}

func main() {
	conf, err := config.Load("config.toml")
	if err != nil {
		log.Error("Couldn't parse config file! |", "Error", err)
	}

	if conf.Hub.Debug {
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
		log.Debug("Debugging mode enabled! Have fun :D")
	}

	ctx := context.Background()

	dnsResolver, err := madns.NewResolver()
	if err != nil {
		log.Fatal("Could not start the DNS resolver |", "Error", err)
	}
	log.Info("Successfully started the DNS resolver!")

	privKey, err := getId()
	if err != nil {
		log.Fatal("Couldn't get private key! | ", "Err", err)
	}

	h, err := libp2p.New(
		libp2p.Identity(privKey),
		libp2p.Ping(true),
		libp2p.ListenAddrStrings(
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", strconv.FormatUint(uint64(conf.Hub.GossipPort), 10)),
		),
		libp2p.Security(noise.ID, noise.New),
		libp2p.Muxer("/mplex/6.7.0", mplex.DefaultTransport),
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info("Started the libp2p host! |", "Addr", fmt.Sprintf("%s/p2p/%s", h.Addrs()[1], h.ID().String()))

	resultsChan := make(chan ResolveResult)

	for _, confPeer := range conf.Hub.BootstrapPeers {
		go func(confPeer string) {
			initPeer, err := multiaddr.NewMultiaddr(confPeer)
			if err != nil {
				resultsChan <- ResolveResult{Error: fmt.Errorf("couldn't parse multiaddr: %w", err)}
				return
			}

			resolvedMultiaddrs, err := dnsResolver.Resolve(ctx, initPeer)
			if err != nil {
				resultsChan <- ResolveResult{Error: fmt.Errorf("can't resolve from DNS addr: %w", err)}
				return
			}

			resultsChan <- ResolveResult{ResolvedMultiaddrs: resolvedMultiaddrs}
		}(confPeer)
	}

	for i := 0; i < len(conf.Hub.BootstrapPeers); i++ {
		result := <-resultsChan
		if result.Error != nil {
			log.Error("DNS resolution error", "Error", result.Error)
			continue
		}

		peerAddrinfo, err := peer.AddrInfoFromP2pAddr(result.ResolvedMultiaddrs[0])
		if err != nil {
			log.Fatal("Can't convert multiaddr to addrinfo", "Error", err)
		}

		log.Info("Connecting to a remote peer! |", "peer", peerAddrinfo)
		err = h.Connect(ctx, *peerAddrinfo)
		if err != nil {
			log.Error("", "Error", err)
		}

		checkConnectionStatus(h, peerAddrinfo.ID)
	}

	// psParams := pubsub.DefaultGossipSubParams()
	// psParams.Dlo = 1
	// log.Debug("GossipSub initial params! |", "Params", psParams)

	// params := pubsub.WithGossipSubParams(psParams)

	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		log.Error(err)
	}

	netwPrimary, err := ReceiveMessages(ctx, ps, h.ID(), "primary", conf)
	if err != nil {
		log.Fatal(err.Error())
	}

	netwContact, err := ReceiveMessages(ctx, ps, h.ID(), "contact_info", conf)
	if err != nil {
		log.Fatal(err.Error())
	}

	netwDiscovery, err := ReceiveMessages(ctx, ps, h.ID(), "peer_discovery", conf)
	if err != nil {
		log.Fatal(err.Error())
	}

	// START THE RPC SERVER
	var wg sync.WaitGroup
	stopCh := make(chan struct{})

	wg.Add(1)
	go Start(&wg, stopCh, *netwPrimary)

	// HANDLE THE MESSAGES
	LoadHandlersFromConf(conf, netwPrimary.NetworkMessage, netwPrimary.logger)
	go HandleContactInfo(netwContact.NetworkMessage, netwContact.logger, h, ctx)
	go logMessages(netwDiscovery.NetworkMessage, netwDiscovery.logger)

	// SEND CONTACT_INFO
	go func() {
		ticker := time.NewTicker(time.Duration(conf.Hub.ContactInterval) * time.Second)
		defer ticker.Stop()

		for {
			<-ticker.C
			netwContact.PublishContactInfo(&protos.ContactInfoContent{
				HubVersion: HUB_VERSION,
				Network:    2,
				GossipAddress: &protos.GossipAddressInfo{
					Family:  4, // to know if address ip4/ip6?
					Address: conf.Hub.PublicHubIp,
					Port:    uint32(conf.Hub.GossipPort),
				},
				Body: &protos.ContactInfoContentBody{
					GossipAddress: &protos.GossipAddressInfo{
						Family:  4,
						Address: conf.Hub.PublicHubIp,
						Port:    uint32(conf.Hub.GossipPort),
					},
					HubVersion: HUB_VERSION,
					Network:    2,
					Timestamp:  uint64(time.Now().Unix()),
					AppVersion: "1.0",
				},
				Timestamp: uint64(time.Now().Unix()),
			})
		}
	}()

	h.Network().Notify(&network.NotifyBundle{
		ConnectedF: func(n network.Network, c network.Conn) {
			log.Info("Peer connected!", "Id", c.RemotePeer(), "Multiaddr", c.RemoteMultiaddr())
		},
	})

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	log.Info("Received signal, shutting down...")

	close(stopCh)

	// shut the node down
	if err := h.Close(); err != nil {
		log.Fatal(err.Error())
	}

	wg.Wait()
}
