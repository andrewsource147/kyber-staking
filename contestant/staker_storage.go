package contestant

import "sync"

type StakerStorage struct {
	representative string
	delegator      []string
	address        string
	tmpStake       uint64
	stakeAmount    uint64
	campArr        []uint64
	mu             *sync.RWMutex
}

func NewStakerStorage(address string) *StakerStorage {
	staker := &StakerStorage{
		representative: address,
		address:        address,
		mu:             &sync.RWMutex{},
	}
	return staker
}

func (s *StakerStorage) CloneForNextEpoch() *StakerStorage {
	s.mu.RLock()
	defer s.mu.RUnlock()

	delegator := make([]string, 0)
	copy(delegator, s.delegator)

	return &StakerStorage{
		representative: s.representative,
		delegator:      delegator,
		address:        s.address,
		stakeAmount:    s.stakeAmount + s.tmpStake,
		mu:             &sync.RWMutex{},
	}
}

func (s *StakerStorage) GetStakeAmount() uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.stakeAmount
}

func (s *StakerStorage) GetTmpStakeAmount() uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tmpStake
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
		s.tmpStake = 0
		s.stakeAmount = s.tmpStake + s.stakeAmount - amount
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
	s.representative = representative
}

func (s *StakerStorage) AddDelegator(delegator string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.delegator = append(s.delegator, delegator)
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
	for _, delegator := range s.delegator {
		if delegator != address {
			newDelegator = append(newDelegator, delegator)
		}
	}
	s.delegator = newDelegator
}
