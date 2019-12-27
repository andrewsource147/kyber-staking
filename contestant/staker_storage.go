package contestant

import (
	"sync"
)

type StakerStorage struct {
	representative    string
	tmpRepresentative string
	delegator         []string
	tmpDelegator      []string
	address           string
	tmpStake          uint64
	stakeAmount       uint64
	campArr           []uint64
	mu                *sync.RWMutex
}

func NewStakerStorage(address string) *StakerStorage {
	staker := &StakerStorage{
		representative:    address,
		tmpRepresentative: address,
		delegator:         make([]string, 0),
		tmpDelegator:      make([]string, 0),
		address:           address,
		mu:                &sync.RWMutex{},
	}
	return staker
}

func (s *StakerStorage) CloneForNextEpoch() *StakerStorage {
	s.mu.RLock()
	defer s.mu.RUnlock()

	delegator := make([]string, len(s.tmpDelegator))
	copy(delegator, s.tmpDelegator)

	return &StakerStorage{
		representative:    s.tmpRepresentative,
		tmpRepresentative: s.tmpRepresentative,
		delegator:         delegator,
		tmpDelegator:      delegator,
		address:           s.address,
		stakeAmount:       s.stakeAmount + s.tmpStake,
		mu:                &sync.RWMutex{},
	}
}

func (s *StakerStorage) GetStakeAmount() uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.stakeAmount
}

func (s *StakerStorage) GetDelegator() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.delegator
}

func (s *StakerStorage) StakeAmount(amount uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tmpStake = s.tmpStake + amount
}

func (s *StakerStorage) WithdrawAmount(amount uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if amount < s.tmpStake {
		s.tmpStake = s.tmpStake - amount
		return
	}
	if amount <= s.tmpStake+s.stakeAmount {
		s.stakeAmount = s.tmpStake + s.stakeAmount - amount
		s.tmpStake = 0
		return
	}
}

func (s *StakerStorage) GetRepresentative() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.representative
}

func (s *StakerStorage) SetRepresentative(representative string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tmpRepresentative = representative
}

func (s *StakerStorage) AddDelegator(delegator string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, address := range s.tmpDelegator {
		if address == delegator {
			return
		}
	}
	s.tmpDelegator = append(s.tmpDelegator, delegator)
}

func (s *StakerStorage) GetAddress() string {
	return s.address
}

func (s *StakerStorage) GetLenCamp() uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return uint64(len(s.campArr))
}

func (s *StakerStorage) Vote(voteId uint64) {
	s.mu.RLock()
	campArr := s.campArr
	s.mu.RUnlock()
	// check element is exist
	for _, id := range campArr {
		if id == voteId {
			return
		}
	}

	campArr = append(campArr, voteId)
	s.mu.Lock()
	s.campArr = campArr
	s.mu.Unlock()
}

func (s *StakerStorage) RemoveDelegator(address string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	newDelegator := make([]string, 0)
	for _, delegator := range s.tmpDelegator {
		if delegator != address {
			newDelegator = append(newDelegator, delegator)
		}
	}
	s.tmpDelegator = newDelegator
}

func (s *StakerStorage) GetTmpStakeAmount() uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tmpStake
}

func (s *StakerStorage) GetTmpRepresentative() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tmpRepresentative
}

func (s *StakerStorage) GetTmpDelegator() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tmpDelegator
}
