package contestant

import (
	"sync"
)

//--------------- This struct for store staking of user -----------------//
type KyberStakingStorage struct {
	epoch map[uint64]*EpochStorage
	mu    *sync.RWMutex
}

func NewKyberStakingStorage() *KyberStakingStorage {
	return &KyberStakingStorage{
		epoch: make(map[uint64]*EpochStorage),
		mu:    &sync.RWMutex{},
	}
}

func (k *KyberStakingStorage) GetEpoch(epochNumber uint64) (epoch *EpochStorage) {
	k.mu.RLock()
	defer k.mu.RUnlock()
	epoch = k.epoch[epochNumber]
	return
}

func (k *KyberStakingStorage) CreateEpoch(epochNumber uint64) *EpochStorage {
	k.mu.Lock()
	defer k.mu.Unlock()
	// find preivous epoch
	var previousEpoch *EpochStorage
	index := epochNumber
	for {
		if index > 0 {
			previousEpoch = k.epoch[index-1]
			if previousEpoch != nil {
				break
			}
		}
		if index == 0 {
			break
		}
		index--
	}
	newEpoch := NewEpochStorage(previousEpoch, epochNumber)
	k.epoch[epochNumber] = newEpoch
	return newEpoch
}

func (k *KyberStakingStorage) GetClosestActiveEpoch(epoch uint64) *EpochStorage {
	k.mu.RLock()
	defer k.mu.RUnlock()

	index := epoch
	for {
		if epoch, ok := k.epoch[index]; ok {
			return epoch
		}
		if index == 0 {
			break
		}
		index--
	}
	return nil
}
