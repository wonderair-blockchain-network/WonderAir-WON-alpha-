// Copyright (c) 2018 The wonderair ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package p2p

import (
	"net"

	"github.com/wonderair ecosystem/go-wonderair ecosystem/ca"
	"github.com/wonderair ecosystem/go-wonderair ecosystem/common"
	"github.com/wonderair ecosystem/go-wonderair ecosystem/core/types"
	"github.com/wonderair ecosystem/go-wonderair ecosystem/log"
	"github.com/wonderair ecosystem/go-wonderair ecosystem/mc"
	"github.com/wonderair ecosystem/go-wonderair ecosystem/p2p/discover"
	"github.com/wonderair ecosystem/go-wonderair ecosystem/params"
	"github.com/wonderair ecosystem/go-wonderair ecosystem/rlp"
)

func UdpStart() {
	addr, err := net.ResolveUDPAddr("udp", ":30000")
	if err != nil {
		log.Error("Can't resolve address: ", "p2p udp", err)
		return
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Error("Error listening:", "p2p udp", err)
		return
	}
	defer conn.Close()

	buf := make([]byte, params.MaxUdpBuf)

	for {
		var mxtxs []*types.Transaction_Mx
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Error("UDP read error", "err", err)
			return
		}

		err = rlp.DecodeBytes(buf[:n], &mxtxs)
		if err != nil {
			log.Error("rlp decode error", "err", err)
			continue
		}
		mc.PublishEvent(mc.SendUdpTx, mxtxs)
	}
}

func UdpSend(data interface{}) {
	bytes, err := rlp.EncodeToBytes(data)
	if err != nil {
		log.Error("error", "p2p udp", err)
		return
	}

	ids := make([]discover.NodeID, 0)
	if ca.InDuration() {
		ids = ca.GetRolesByGroupOnlyNextElect(common.RoleValidator | common.RoleBackupValidator)
	} else {
		ids = ca.GetRolesByGroup(common.RoleValidator | common.RoleBackupValidator)
	}
	if len(ids) <= 2 {
		for _, id := range ids {
			send(id, bytes)
		}
		return
	}

	is := Random(len(ids), 2)
	for _, i := range is {
		send(ids[i], bytes)
	}
}

func send(id discover.NodeID, data []byte) {
	node := ServerP2p.ntab.Resolve(id)
	if node == nil {
		log.Error("buckets nodes", "p2p", id)
		return
	}

	addr, err := net.ResolveUDPAddr("udp", node.IP.String()+":30000")
	if err != nil {
		log.Error("Can't resolve address: ", "p2p udp", err)
		return
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Error("Can't dial: ", "p2p udp", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write(data)
	if err != nil {
		log.Error("failed:", "p2p udp", err)
		return
	}
}
