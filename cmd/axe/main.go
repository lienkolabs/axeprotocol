package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lienkolabs/axeprotocol/social/state"
	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/breeze/network"
	"github.com/lienkolabs/breeze/network/echo"
)

type ConfigNode struct {
	BlockServiceAddress string
	BlockServiceToken   string
	BlockBoardcastPort  int
	StateDataPath       string
}

func StartNode(credentials crypto.PrivateKey, config ConfigNode) (*echo.SocialNode, error) {
	//blockchain := BlockChain(credentials)
	nodeConfig := echo.SocialNodeConfig{
		Credentials:            credentials,
		SocialCode:             echo.ProtocolCode{0, 0, 0, 1},
		BlockServiceAddress:    config.BlockServiceAddress,
		BlockServiveToken:      crypto.TokenFromString(config.BlockServiceToken),
		BlockBroadcastFirewall: network.AcceptAllConnections,
	}

	genesis := state.NewGenesisState(config.StateDataPath)
	node, err := echo.NewSocialNodeListener(&nodeConfig, genesis)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func main() {
	var config ConfigNode
	if len(os.Args) < 2 {
		log.Fatalln("usage: axe path-to-config-file.json")
	}
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("could not read config file: %v\n", err)
	}
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatalf("could not read config file: %v\n", err)
	}
	if config.BlockServiceAddress == "" {
		log.Fatalf("invalid configuration: block provider address must be provided")
	}
	if token := crypto.TokenFromString(config.BlockServiceToken); token.Equal(crypto.ZeroToken) {
		log.Fatalf("invalid configuration: valid block provider token must be provided")
	}
	if config.BlockBoardcastPort == 0 {
		log.Fatalf("invalid configuration: broadcast port must be provided")
	}
	_, credentials := crypto.RandomAsymetricKey()
	node, err := StartNode(credentials, config)
	if err != nil {
		log.Fatalf("could not start node: %v", err)
	}
	c := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	for {
		<-c
		node.Shutdown()
		return
	}
}
