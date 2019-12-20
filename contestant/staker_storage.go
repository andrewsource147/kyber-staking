package contestant

import "sync"

type StakerStorage struct {
	representative string
	delegator      []string
	address        string
	stakeAmount    uint64
	holdAmount     uint64
	newStakeAmount uint64
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

func (s *StakerStorage) Copy() *StakerStorage {
	s.mu.RLock()
	defer s.mu.RUnlock()

	delegator := make([]string, 0)
	copy(delegator, s.delegator)

	return &StakerStorage{
		representative: s.representative,
		delegator:      delegator,
		address:        s.address,
		stakeAmount:    s.stakeAmount,
		holdAmount:     s.stakeAmount,
		mu:             &sync.RWMutex{},
	}
}

func (s *StakerStorage) GetStakeAmount() uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.stakeAmount
}
func (s *StakerStorage) GetHoldAmount() uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.holdAmount
}

func (s *StakerStorage) GetDelegator() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.delegator
}

func (s *StakerStorage) StakeAmount(amount uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stakeAmount = s.stakeAmount + amount
	s.holdAmount = s.holdAmount + amount
}

func (s *StakerStorage) WithdrawAmount(amount uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if amount > s.stakeAmount {
		s.stakeAmount = 0
		s.holdAmount = 0
	}
	s.stakeAmount = s.stakeAmount - amount
	s.holdAmount = s.holdAmount - amount
}

func (s *StakerStorage) AddNewHoldAmount(amount uint64) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.newStakeAmount = s.newStakeAmount + amount
}

func (s *StakerStorage) WithdrawHoldAmount(amount uint64) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if amount <= s.newStakeAmount {
		s.newStakeAmount = s.newStakeAmount - amount
		return
	}
	if amount <= s.newStakeAmount+s.holdAmount {
		s.holdAmount = s.holdAmount + s.newStakeAmount - amount
		s.newStakeAmount = 0
		return
	}

	s.holdAmount = 0
	s.newStakeAmount = 0
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
