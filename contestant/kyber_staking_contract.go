package contestant

type KyberStakingContract struct {
	startBlock    uint64
	epochDuration uint64
	storage       *KyberStakingStorage
}

func NewKyberStakingContract(startBlock uint64, epochDuration uint64) *KyberStakingContract {
	return &KyberStakingContract{
		startBlock:    startBlock,
		epochDuration: epochDuration,
		storage:       NewKyberStakingStorage(),
	}
}

// naming explain
// E is epoch instance
// a is staker
// b is current representative

func (sc *KyberStakingContract) Stake(block uint64, amount uint64, staker string) {
	epochNumber := sc.getEpochNumberbyBlock(block) + 1
	E := sc.storage.GetEpoch(epochNumber)
	if E == nil {
		E = sc.storage.CreateEpoch(epochNumber)
	}

	// cumulative amount for staker
	a := E.GetOrCreateStaker(staker)
	a.StakeAmount(amount)

	// handle for reward
	preE := E.GetPreviousEpoch()
	if preE == nil {
		return
	}
	preA := preE.GetStaker(staker)
	if preA == nil {
		return
	}
	preA.AddNewHoldAmount(amount)
}

func (sc *KyberStakingContract) Withdraw(block uint64, amount uint64, staker string) {
	epochNumber := sc.getEpochNumberbyBlock(block) + 1
	E := sc.storage.GetEpoch(epochNumber)
	if E == nil {
		E = sc.storage.CreateEpoch(epochNumber)
	}

	a := E.GetStaker(staker)
	if a == nil {
		return
	}
	a.WithdrawAmount(amount)

	// handle for reward
	preE := E.GetPreviousEpoch()
	if preE == nil {
		return
	}
	preA := preE.GetStaker(staker)
	if preA == nil {
		return
	}
	preA.WithdrawHoldAmount(amount)
}

func (sc *KyberStakingContract) Delegate(block uint64, staker string, representative string) {
	epochNumber := sc.getEpochNumberbyBlock(block) + 1
	E := sc.storage.GetEpoch(epochNumber)
	if E == nil {
		E = sc.storage.CreateEpoch(epochNumber)
	}

	a := E.GetOrCreateStaker(staker)
	a.SetRepresentative(representative)

	b := E.GetOrCreateStaker(representative)
	b.AddDelegator(staker)
}

func (sc *KyberStakingContract) Vote(block uint64, voteid uint64, staker string) {
	epochNumber := sc.getEpochNumberbyBlock(block)
	E := sc.storage.GetEpoch(epochNumber)
	if E == nil {
		E = sc.storage.CreateEpoch(epochNumber)
	}

	a := E.GetOrCreateStaker(staker)
	a.Vote(voteid)
}

func (sc *KyberStakingContract) GetStake(epoch uint64, staker string) (stake uint64) {
	E := sc.storage.GetClosestActiveEpoch(epoch)

	if E == nil {
		return
	}

	a := E.GetStaker(staker)
	if a == nil {
		return
	}

	stake = a.GetStakeAmount()
	return
}

func (sc *KyberStakingContract) GetDelegatedStake(epoch uint64, staker string) (delegatedStake uint64) {
	E := sc.storage.GetClosestActiveEpoch(epoch)

	if E == nil {
		return
	}

	a := E.GetStaker(staker)
	if a == nil {
		return
	}
	delegators := a.GetDelegator()
	for _, address := range delegators {
		delegator := E.GetStaker(address)
		delegatedStake = delegatedStake + delegator.GetStakeAmount()
	}
	return
}

func (sc *KyberStakingContract) GetRepresentative(epoch uint64, staker string) (poolmaster string) {
	E := sc.storage.GetClosestActiveEpoch(epoch)
	if E == nil {
		return staker
	}
	a := E.GetStaker(staker)
	if a == nil {
		return staker
	}
	poolmaster = a.GetRepresentative()
	return
}

func (sc *KyberStakingContract) GetReward(epoch uint64, staker string) (percentage float64) {
	E := sc.storage.GetClosestActiveEpoch(epoch)
	if E == nil {
		return
	}

	stakerPoint := E.GetStakerPoint(staker)
	if stakerPoint == 0 {
		return
	}
	epochTotalPoint := E.GetTotalPoint()
	if stakerPoint == 0 {
		return
	}

	percentage = float64(stakerPoint) / float64(epochTotalPoint)
	return
}

func (sc *KyberStakingContract) GetPoolReward(epoch uint64, staker string) (percentage float64) {
	E := sc.storage.GetClosestActiveEpoch(epoch)
	if E == nil {
		return
	}
	a := E.GetStaker(staker)
	if a == nil {
		return
	}
	bAddr := a.GetRepresentative()

	aHoldAmount := a.GetHoldAmount()
	if aHoldAmount == 0 {
		return
	}

	// get hold amount of
	b := E.GetStaker(bAddr)
	bTotalAmount := b.GetHoldAmount()
	bDelegators := b.GetDelegator()
	for _, address := range bDelegators {
		delegator := E.GetStaker(address)
		bTotalAmount = bTotalAmount + delegator.GetHoldAmount()
	}

	if bTotalAmount == 0 {
		return
	}

	percentage = float64(aHoldAmount) / float64(bTotalAmount)
	return
}

func (sc *KyberStakingContract) getEpochNumberbyBlock(block uint64) (epoch uint64) {
	startBlock := sc.startBlock
	epochDuration := sc.epochDuration

	if block < startBlock {
		return
	}
	epoch = (block-startBlock)/epochDuration + 1
	return
}
