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
	epochNumber := sc.getEpochNumberbyBlock(block)
	E := sc.storage.GetEpoch(epochNumber)
	if E == nil {
		E = sc.storage.CreateEpoch(epochNumber)
	}

	// cumulative amount for staker
	a := E.GetOrCreateStaker(staker)
	a.StakeAmount(amount)
}

func (sc *KyberStakingContract) Withdraw(block uint64, amount uint64, staker string) {
	epochNumber := sc.getEpochNumberbyBlock(block)
	E := sc.storage.GetEpoch(epochNumber)
	if E == nil {
		E = sc.storage.CreateEpoch(epochNumber)
	}

	a := E.GetStaker(staker)
	if a == nil {
		return
	}

	a.WithdrawAmount(amount)
}

func (sc *KyberStakingContract) Delegate(block uint64, staker string, representative string) {
	epochNumber := sc.getEpochNumberbyBlock(block)
	E := sc.storage.GetEpoch(epochNumber)
	if E == nil {
		E = sc.storage.CreateEpoch(epochNumber)
	}

	a := E.GetOrCreateStaker(staker)

	if a.GetTmpRepresentative() != representative {
		b := E.GetStaker(a.GetTmpRepresentative())
		b.RemoveDelegator(staker)
	}

	a.SetRepresentative(representative)
	b := E.GetOrCreateStaker(representative)
	if representative != staker {
		b.AddDelegator(staker)
	}
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

	if E.GetEpochNumber() == epoch {
		stake = a.GetStakeAmount()
	} else {
		stake = a.GetStakeAmount() + a.GetTmpStakeAmount()
	}
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

	if E.GetEpochNumber() == epoch {
		delegators := a.GetDelegator()
		for _, address := range delegators {
			delegator := E.GetStaker(address)
			if delegator != nil {
				delegatedStake = delegatedStake + delegator.GetStakeAmount()
			}

		}
	} else {
		delegators := a.GetTmpDelegator()
		for _, address := range delegators {
			delegator := E.GetStaker(address)
			if delegator != nil {
				delegatedStake = delegatedStake + delegator.GetStakeAmount() + delegator.GetTmpStakeAmount()
			}
		}
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

	if E.GetEpochNumber() == epoch {
		poolmaster = a.GetRepresentative()
	} else {
		poolmaster = a.GetTmpRepresentative()
	}
	return
}

func (sc *KyberStakingContract) GetReward(epoch uint64, staker string) (percentage float64) {
	E := sc.storage.GetClosestActiveEpoch(epoch)

	if E == nil || E.GetEpochNumber() != epoch {
		return
	}

	a := E.GetStaker(staker)

	if a == nil {
		return
	}

	// get staker amount

	stakerPoint := E.GetStakerPoint(staker)
	if stakerPoint == 0 {
		return
	}

	epochTotalPoint := E.GetTotalPoint()
	if epochTotalPoint == 0 {
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

	if E.GetEpochNumber() == epoch {
		aStakeAmount := a.GetStakeAmount()
		bAddr := a.GetRepresentative()
		b := E.GetStaker(bAddr)
		if b == nil {
			return
		}

		var bTotalAmount uint64

		// Count b stake amount if b not delegate for anyone else
		if b.GetRepresentative() == b.GetAddress() {
			bTotalAmount = b.GetStakeAmount()
		}

		bDelegators := b.GetDelegator()
		for _, address := range bDelegators {
			delegator := E.GetStaker(address)
			if delegator != nil {
				bTotalAmount = bTotalAmount + delegator.GetStakeAmount()
			}
		}

		if bTotalAmount == 0 {
			return
		}
		percentage = float64(aStakeAmount) / float64(bTotalAmount)

	} else {
		aStakeAmount := a.GetStakeAmount() + a.GetTmpStakeAmount()
		bAddr := a.GetTmpRepresentative()
		b := E.GetStaker(bAddr)
		if b == nil {
			return
		}

		var bTotalAmount uint64
		if b.GetTmpRepresentative() == b.GetAddress() {
			bTotalAmount = b.GetStakeAmount() + b.GetTmpStakeAmount()
		}

		bDelegators := b.GetTmpDelegator()
		for _, address := range bDelegators {
			delegator := E.GetStaker(address)
			if delegator != nil {
				bTotalAmount = bTotalAmount + delegator.GetStakeAmount() + delegator.GetTmpStakeAmount()
			}
		}
		if bTotalAmount == 0 {
			return
		}
		percentage = float64(aStakeAmount) / float64(bTotalAmount)
	}
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
