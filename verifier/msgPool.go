// Copyright (c) 2018 The wonderair ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package verifier

import (
	"github.com/wonderair ecosystem/go-wonderair ecosystem/common"
	"github.com/wonderair ecosystem/go-wonderair ecosystem/core/types"
	"github.com/wonderair ecosystem/go-wonderair ecosystem/mc"
	"github.com/wonderair ecosystem/go-wonderair ecosystem/params/man"
	"github.com/pkg/errors"
	"time"
)

type msgPool struct {
	parentHeader     *types.Header
	posNotifyCache   []*mc.BlockPOSFinishedNotify
	inquiryReqCache  map[common.Address]*mc.HD_ReelectInquiryReqMsg
	rlConsensusCache map[uint32]*mc.HD_ReelectLeaderConsensus
}

func newMsgPool() *msgPool {
	return &msgPool{
		parentHeader:     nil,
		posNotifyCache:   make([]*mc.BlockPOSFinishedNotify, 0),
		inquiryReqCache:  make(map[common.Address]*mc.HD_ReelectInquiryReqMsg),
		rlConsensusCache: make(map[uint32]*mc.HD_ReelectLeaderConsensus),
	}
}

func (mp *msgPool) SavePOSNotifyMsg(msg *mc.BlockPOSFinishedNotify) error {
	if nil == msg || (msg.Header.Leader == common.Address{}) {
		return ErrMsgIsNil
	}

	for _, oldMsg := range mp.posNotifyCache {
		if oldMsg.ConsensusTurn == msg.ConsensusTurn && oldMsg.Header.Leader == msg.Header.Leader {
			return ErrMsgExistInCache
		}
	}

	mp.posNotifyCache = append(mp.posNotifyCache, msg)
	return nil
}

func (mp *msgPool) GetPOSNotifyMsg(leader common.Address, consensusTurn uint32) (*mc.BlockPOSFinishedNotify, error) {
	for _, msg := range mp.posNotifyCache {
		if msg.ConsensusTurn == msg.ConsensusTurn && msg.Header.Leader == msg.Header.Leader {
			return msg, nil
		}
	}
	return nil, ErrNoMsgInCache
}

func (mp *msgPool) SaveInquiryReqMsg(msg *mc.HD_ReelectInquiryReqMsg) {
	if nil == msg || (msg.From == common.Address{}) {
		return
	}

	old, exist := mp.inquiryReqCache[msg.From]
	if exist && old.TimeStamp > msg.TimeStamp {
		return
	}
	mp.inquiryReqCache[msg.From] = msg
}

func (mp *msgPool) GetInquiryReqMsg(leader common.Address) (*mc.HD_ReelectInquiryReqMsg, error) {
	msg, OK := mp.inquiryReqCache[leader]
	if !OK {
		return nil, ErrNoMsgInCache
	}

	passTime := time.Now().Unix() - msg.TimeStamp
	if passTime > man.LRSReelectInterval {
		delete(mp.inquiryReqCache, leader)
		return nil, errors.Errorf("消息已过期, timestamp=%d, passTime=%d", msg.TimeStamp, passTime)
	}
	return msg, nil
}

func (mp *msgPool) SaveRLConsensusMsg(msg *mc.HD_ReelectLeaderConsensus) {
	if nil == msg || (msg.Req.InquiryReq.Master == common.Address{}) {
		return
	}
	mp.rlConsensusCache[msg.Req.InquiryReq.ConsensusTurn] = msg
}

func (mp *msgPool) GetRLConsensusMsg(consensusTurn uint32) (*mc.HD_ReelectLeaderConsensus, error) {
	msg, OK := mp.rlConsensusCache[consensusTurn]
	if !OK {
		return nil, ErrNoMsgInCache
	}
	return msg, nil
}

func (mp *msgPool) SaveParentHeader(header *types.Header) {
	if nil == header {
		return
	}
	mp.parentHeader = header
}

func (mp *msgPool) GetParentHeader() *types.Header {
	return mp.parentHeader
}
