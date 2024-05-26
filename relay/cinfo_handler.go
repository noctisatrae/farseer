package main

import (
	"context"
	"fmt"
	"strconv"

	protos "farseer/protos"

	"github.com/charmbracelet/log"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

// Channel => Message => Content
func HandleContactInfo(contactInfoChan chan *protos.GossipMessage, ll log.Logger, h host.Host, ctx context.Context) {
	for contactInfoMessage := range contactInfoChan {
		remotePeerId, err := peer.IDFromBytes(contactInfoMessage.GetPeerId())
		if err != nil {
			ll.Error("Can't serialize the peer id from message! |", "Error", err)
		} else {
			remotePeerIdStr := remotePeerId.String()
			cinfo := contactInfoMessage.GetContactInfoContent()

			remotePeerAddrFamily := cinfo.GossipAddress.GetFamily()
			remotePeerAddr := cinfo.GossipAddress.GetAddress()
			remotePeerPort := cinfo.GossipAddress.GetPort()

			ll.Info("Received contact info! |", "Addr", remotePeerAddr, "Port", remotePeerPort, "Family", remotePeerAddrFamily)

			remotePeerMultiAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf(
				"/ip%s/%s/tcp/%s/p2p/%s",
				strconv.FormatUint(uint64(remotePeerAddrFamily), 10),
				remotePeerAddr,
				strconv.FormatUint(uint64(remotePeerPort), 10),
				remotePeerIdStr,
			))
			ll.Debug("From contact info! |", "Multiaddr", remotePeerMultiAddr)
			if err != nil {
				ll.Error("Can't parse multiaddrr from contact info! |", "Error", err)
			} else {

				remotePeerAddrInfo, err := peer.AddrInfoFromP2pAddr(remotePeerMultiAddr)
				if err != nil {
					ll.Error("Can't create AddrInfo from contact info! |", "Error", err, "Multiaddr", remotePeerMultiAddr)
				} else {
					err = h.Connect(ctx, *remotePeerAddrInfo)
					if err != nil {
						ll.Error("Couldn't connect to peer from contact info! |", "Error", err)
					} else {
						ll.Info("Connected to peer from contact info! |", "Addr", remotePeerAddr)
					}
				}
			}
		}
	}
}
