package main

import (
	"errors"
	"os"

	"github.com/charmbracelet/log"
	"github.com/libp2p/go-libp2p/core/crypto"
)

func getId() (crypto.PrivKey, error) {
	if _, err := os.Stat("./hub_identity"); errors.Is(err, os.ErrNotExist) {
		log.Info("No hub_identity file found! Regenating one!")
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
		log.Info("hub_identity already exist, not regenerating!")
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
	_, err := getId()
	if err != nil {
		log.Fatal(err)
	}
}