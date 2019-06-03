// Copyright (c) 2018 The wonderair ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php


package storage

import (
	"hash"
)

const (
	BMTHash  = "BMT"
	SHA3Hash = "SHA3" // http://golang.org/pkg/hash/#Hash
)

type SwarmHash interface {
	hash.Hash
	ResetWithLength([]byte)
}

type HashWithLength struct {
	hash.Hash
}

func (self *HashWithLength) ResetWithLength(length []byte) {
	self.Reset()
	self.Write(length)
}
