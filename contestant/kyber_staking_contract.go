package contestant

type KyberStakingContract struct {
}

func NewKyberStakingContract(startBlock uint64, epochDuration uint64) *KyberStakingContract {
	return &KyberStakingContract{}
}

func (sc *KyberStakingContract) Stake(block uint64, amount uint64, staker string) {

}

func (sc *KyberStakingContract) Withdraw(block uint64, amount uint64, staker string) {

}

func (sc *KyberStakingContract) Delegate(block uint64, staker string, representative string) {

}

func (sc *KyberStakingContract) Vote(block uint64, voteid uint64, staker string) {

}

func (sc *KyberStakingContract) GetStake(epoch uint64, staker string) (stake uint64) {
	return 0
}

func (sc *KyberStakingContract) GetDelegatedStake(epoch uint64, staker string) (delegatedStake uint64) {
	return 0
}

func (sc *KyberStakingContract) GetRepresentative(epoch uint64, staker string) (poolmaster string) {
	return "victor"
}

func (sc *KyberStakingContract) GetReward(epoch uint64, staker string) (percentage float64) {
	return float64(0)
}

func (sc *KyberStakingContract) GetPoolReward(epoch uint64, staker string) (percentage float64) {
	return float64(0)
}
