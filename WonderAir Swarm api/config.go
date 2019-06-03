// Copyright (c) 2018 The wonderair ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php


package api

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wonderair ecosystem/go-wonderair ecosystem/common"
	"github.com/wonderair ecosystem/go-wonderair ecosystem/contracts/ens"
	"github.com/wonderair ecosystem/go-wonderair ecosystem/crypto"
	"github.com/wonderair ecosystem/go-wonderair ecosystem/log"
	"github.com/wonderair ecosystem/go-wonderair ecosystem/node"
	"github.com/wonderair ecosystem/go-wonderair ecosystem/swarm/network"
	"github.com/wonderair ecosystem/go-wonderair ecosystem/swarm/services/swap"
	"github.com/wonderair ecosystem/go-wonderair ecosystem/swarm/storage"
)

const (
	DefaultHTTPListenAddr = "127.0.0.1"
	DefaultHTTPPort       = "8500"
)

// separate bzz directories
// allow several bzz nodes running in parallel
type Config struct {
	// serialised/persisted fields
	*storage.StoreParams
	*storage.ChunkerParams
	*network.HiveParams
	Swap *swap.SwapParams
	*network.SyncParams
	Contract    common.Address
	EnsRoot     common.Address
	EnsAPIs     []string
	Path        string
	ListenAddr  string
	Port        string
	PublicKey   string
	BzzKey      string
	NetworkId   uint64
	SwapEnabled bool
	SyncEnabled bool
	SwapApi     string
	Cors        string
	BzzAccount  string
	BootNodes   string
}

//create a default config with all parameters to set to defaults
func NewDefaultConfig() (self *Config) {

	self = &Config{
		StoreParams:   storage.NewDefaultStoreParams(),
		ChunkerParams: storage.NewChunkerParams(),
		HiveParams:    network.NewDefaultHiveParams(),
		SyncParams:    network.NewDefaultSyncParams(),
		Swap:          swap.NewDefaultSwapParams(),
		ListenAddr:    DefaultHTTPListenAddr,
		Port:          DefaultHTTPPort,
		Path:          node.DefaultDataDir(),
		EnsAPIs:       nil,
		EnsRoot:       ens.TestNetAddress,
		NetworkId:     network.NetworkId,
		SwapEnabled:   false,
		SyncEnabled:   true,
		SwapApi:       "",
		BootNodes:     "",
	}

	return
}

//some config params need to be initialized after the complete
//config building phase is completed (e.g. due to overriding flags)
func (self *Config) Init(prvKey *ecdsa.PrivateKey) {

	address := crypto.PubkeyToAddress(prvKey.PublicKey)
	self.Path = filepath.Join(self.Path, "bzz-"+common.Bytes2Hex(address.Bytes()))
	err := os.MkdirAll(self.Path, os.ModePerm)
	if err != nil {
		log.Error(fmt.Sprintf("Error creating root swarm data directory: %v", err))
		return
	}

	pubkey := crypto.FromECDSAPub(&prvKey.PublicKey)
	pubkeyhex := common.ToHex(pubkey)
	keyhex := crypto.Keccak256Hash(pubkey).Hex()

	self.PublicKey = pubkeyhex
	self.BzzKey = keyhex

	self.Swap.Init(self.Contract, prvKey)
	self.SyncParams.Init(self.Path)
	self.HiveParams.Init(self.Path)
	self.StoreParams.Init(self.Path)
}
