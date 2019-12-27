package contestant

import (
	"sync"
)

type EpochStorage struct {
	epocNumber    uint64
	previousEpoch *EpochStorage
	stake         map[string]*StakerStorage
	mu            *sync.RWMutex
}

func NewEpochStorage(previousEpoch *EpochStorage, epocNumber uint64) *EpochStorage {
	// clone all state
	if previousEpoch != nil {
		return previousEpoch.CloneForNextEpoch(epocNumber)
	}
	return &EpochStorage{
		epocNumber: epocNumber,
		stake:      make(map[string]*StakerStorage),
		mu:         &sync.RWMutex{},
	}
}

func (e *EpochStorage) CloneForNextEpoch(epocNumber uint64) *EpochStorage {
	e.mu.RLock()
	defer e.mu.RUnlock()

	stake := make(map[string]*StakerStorage)
	for address, staker := range e.stake {
		stake[address] = staker.CloneForNextEpoch()
	}
	return &EpochStorage{
		previousEpoch: e,
		epocNumber:    epocNumber,
		stake:         stake,
		mu:            &sync.RWMutex{},
	}
}

func (e *EpochStorage) GetEpochNumber() uint64 {
	return e.epocNumber
}

func (e *EpochStorage) GetPreviousEpoch() *EpochStorage {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.previousEpoch
}

func (e *EpochStorage) GetStaker(address string) *StakerStorage {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.stake[address]
}

func (e *EpochStorage) GetOrCreateStaker(address string) *StakerStorage {
	e.mu.RLock()
	staker := e.stake[address]
	e.mu.RUnlock()
	if staker != nil {
		return staker
	}
	newStaker := NewStakerStorage(address)
	e.mu.Lock()
	e.stake[address] = newStaker
	e.mu.Unlock()
	return newStaker
}

func (e *EpochStorage) GetStakerPoint(address string) uint64 {
	e.mu.RLock()
	defer e.mu.RUnlock()

	a := e.GetStaker(address)
	if a == nil {
		return 0
	}

	var totalStake uint64

	if a.GetRepresentative() == a.GetAddress() {
		totalStake = a.GetStakeAmount()
	}

	delegators := a.GetDelegator()
	for _, dAddr := range delegators {
		delegator := e.GetStaker(dAddr)
		if delegator != nil {
			totalStake = totalStake + delegator.GetStakeAmount()
		}
	}

	return totalStake * a.GetLenCamp()
}

func (e *EpochStorage) GetTotalPoint() uint64 {
	e.mu.RLock()
	stake := e.stake
	e.mu.RUnlock()

	var totalPoint uint64
	for _, staker := range stake {
		totalPoint = totalPoint + e.GetStakerPoint(staker.GetAddress())
	}
	return totalPoint
}
