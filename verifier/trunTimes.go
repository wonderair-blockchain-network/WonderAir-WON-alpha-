// Copyright (c) 2018 The wonderair ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package verifier

import (
	"github.com/wonderair ecosystem/go-wonderair ecosystem/params/man"
)

type turnTimes struct {
	beginTimes map[uint32]int64
}

func newTurnTimes() *turnTimes {
	tt := &turnTimes{
		beginTimes: make(map[uint32]int64),
	}

	return tt
}

func (tt *turnTimes) SetBeginTime(consensusTurn uint32, time int64) bool {
	if oldTime, exist := tt.beginTimes[consensusTurn]; exist {
		if time <= oldTime {
			return false
		}
	}
	tt.beginTimes[consensusTurn] = time
	return true
}

func (tt *turnTimes) GetBeginTime(consensusTurn uint32) int64 {
	if beginTime, exist := tt.beginTimes[consensusTurn]; exist {
		return beginTime
	} else {
		return 0
	}
}

func (tt *turnTimes) GetPosEndTime(consensusTurn uint32) int64 {
	posTime := man.LRSPOSOutTime
	if consensusTurn == 0 {
		posTime += man.LRSParentMiningTime
	}

	return tt.GetBeginTime(consensusTurn) + posTime
}

func (tt *turnTimes) CalState(consensusTurn uint32, time int64) (st state, remainTime int64, reelectTurn uint32) {
	posTime := man.LRSPOSOutTime
	if consensusTurn == 0 {
		posTime += man.LRSParentMiningTime
	}

	passTime := time - tt.GetBeginTime(consensusTurn)
	if passTime < posTime {
		return stPos, posTime - passTime, 0
	}

	st = stReelect
	reelectTurn = uint32((passTime-posTime)/man.LRSReelectOutTime) + 1
	remainTime = (passTime - posTime) % man.LRSReelectOutTime
	if remainTime == 0 {
		remainTime = man.LRSReelectOutTime
	}
	return
}

func (tt *turnTimes) CalRemainTime(consensusTurn uint32, reelectTurn uint32, time int64) int64 {
	posTime := man.LRSPOSOutTime
	if consensusTurn == 0 {
		posTime += man.LRSParentMiningTime
	}
	deadLine := tt.GetBeginTime(consensusTurn) + posTime + int64(reelectTurn)*man.LRSReelectOutTime
	return deadLine - time
}
