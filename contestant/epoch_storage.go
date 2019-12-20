package contestant

import (
	"sync"
)

type EpochStorage struct {
	previousEpoch *EpochStorage
	stake         map[string]*StakerStorage
	mu            *sync.RWMutex
}

func NewEpochStorage(previousEpoch *EpochStorage) *EpochStorage {
	// clone all state
	if previousEpoch != nil {
		return previousEpoch.Copy()
	}
	return &EpochStorage{
		stake: make(map[string]*StakerStorage),
		mu:    &sync.RWMutex{},
	}
}

func (e *EpochStorage) Copy() *EpochStorage {
	e.mu.RLock()
	defer e.mu.RUnlock()

	stake := make(map[string]*StakerStorage)
	for address, staker := range e.stake {
		stake[address] = staker.Copy()
	}
	return &EpochStorage{
		previousEpoch: e,
		stake:         stake,
		mu:            &sync.RWMutex{},
	}
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
	if a == nil || a.GetLenCamp() == 0 {
		return 0
	}
	totalStake := a.GetHoldAmount()
	delegators := a.GetDelegator()
	for _, dAddr := range delegators {
		delegator := e.GetStaker(dAddr)
		if delegator != nil {
			totalStake = totalStake + delegator.GetHoldAmount()
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
