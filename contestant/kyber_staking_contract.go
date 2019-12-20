package contestant

import "log"

func contains(slice []uint64, item uint64) bool {
	set := make(map[uint64]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}
func containsString(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

type KyberStakingContract struct {
	startBlock, epochDuration, latestEpoch uint64
	book map[uint64]map[string]*StakerStruct
	book2 map[string]uint64
	voters map[uint64][]string
}

type StakerStruct struct {
	currentStake, stake, withdraw uint64
	vote []uint64
	representativeForNextEpoch []string
	representativeFor []string
	delegateToNextEpoch string
	delegateTo string
}

func NewKyberStakingContract() *KyberStakingContract {
	return &KyberStakingContract{0, 0, 0,
		make(map[uint64]map[string]*StakerStruct), make(map[string]uint64), make(map[uint64][]string)}
}
func (sc *KyberStakingContract) Init(startBlock uint64, epochDuration uint64){
	sc.startBlock = startBlock
	sc.epochDuration = epochDuration
}
func (sc *KyberStakingContract) GetEpoch(block uint64) uint64{
	//block < @startBlock ? 0 : ((block - @startBlock)/@epoch + 1)
	if block < sc.startBlock {
		return 0
	}else{
		return (block-sc.startBlock)/sc.epochDuration + 1
	}
}

func (sc *KyberStakingContract) GetBook() (map[uint64]map[string]*StakerStruct){
	return sc.book
}

func (sc *KyberStakingContract) InitEpoch(epoch uint64, staker string) {
	if _, ok := sc.book[epoch]; !ok {
		sc.book[epoch] = make(map[string]*StakerStruct)
		sc.voters[epoch] = []string{}
	}

	if _, ok := sc.book[epoch][staker]; !ok {
		sc.book[epoch][staker] = &StakerStruct{0, 0, 0, []uint64{}, []string{staker}, []string{staker}, staker, staker}
		if _,ok1 := sc.book2[staker]; ok1 {
			prevTemp := sc.book[sc.book2[staker]][staker]
			sc.book[epoch][staker].representativeForNextEpoch = prevTemp.representativeForNextEpoch
			sc.book[epoch][staker].representativeFor = prevTemp.representativeForNextEpoch
			sc.book[epoch][staker].delegateToNextEpoch = prevTemp.delegateToNextEpoch
			sc.book[epoch][staker].delegateTo = prevTemp.delegateToNextEpoch
			sc.book[epoch][staker].stake = prevTemp.stake + prevTemp.currentStake - prevTemp.withdraw
		}
		if epoch >= sc.latestEpoch {
			sc.book2[staker] = epoch
		}
	}
}

func (sc *KyberStakingContract) Stake(block uint64, amount uint64, staker string) {
	epoch := sc.GetEpoch(block)
	if epoch > sc.latestEpoch { sc.latestEpoch = epoch }
	sc.InitEpoch(epoch, staker)
	sc.book[epoch][staker].currentStake += amount
}

func (sc *KyberStakingContract) Withdraw(block uint64, amount uint64, staker string) {
	epoch := sc.GetEpoch(block)
	if epoch > sc.latestEpoch { sc.latestEpoch = epoch }
	sc.InitEpoch(epoch, staker)

	if sc.book[epoch][staker].currentStake >= amount {
		sc.book[epoch][staker].currentStake -= amount
	}else {
		if sc.book[epoch][staker].withdraw + amount <= sc.book[epoch][staker].currentStake + sc.book[epoch][staker].stake{
			sc.book[epoch][staker].withdraw += amount - sc.book[epoch][staker].currentStake
			sc.book[epoch][staker].currentStake = 0
		}else {
			log.Println("skip cause withdraw too much")
		}
	}
}

func (sc *KyberStakingContract) Delegate(block uint64, staker string, representative string) {
	epoch := sc.GetEpoch(block)
	if epoch > sc.latestEpoch { sc.latestEpoch = epoch }
	sc.InitEpoch(epoch, staker)
	sc.InitEpoch(epoch, representative)
	sc.InitEpoch(epoch, sc.book[epoch][staker].delegateToNextEpoch)
	if staker == representative {
		remove(sc.book[epoch][ sc.book[epoch][staker].delegateToNextEpoch ].representativeForNextEpoch, staker)
		if !containsString(sc.book[epoch][staker].representativeForNextEpoch, staker) {
			sc.book[epoch][staker].representativeForNextEpoch = append(sc.book[epoch][staker].representativeForNextEpoch, staker)
			sc.book[epoch][staker].delegateToNextEpoch = staker
		}
	}else {
		if !containsString(sc.book[epoch][representative].representativeForNextEpoch, staker) {
			sc.book[epoch][representative].representativeForNextEpoch = append(sc.book[epoch][representative].representativeForNextEpoch, staker)
		}
		remove(sc.book[epoch][ sc.book[epoch][staker].delegateToNextEpoch ].representativeForNextEpoch, staker)

		sc.book[epoch][staker].delegateToNextEpoch = representative
		remove(sc.book[epoch][staker].representativeForNextEpoch, staker)
	}
}

func (sc *KyberStakingContract) Vote(block uint64, voteid uint64, staker string) {
	epoch := sc.GetEpoch(block)

	if epoch > sc.latestEpoch { sc.latestEpoch = epoch }
	sc.InitEpoch(epoch, staker)

	if !contains(sc.book[epoch][staker].vote, voteid) {
		sc.book[epoch][staker].vote = append(sc.book[epoch][staker].vote, voteid)
	}

	if !containsString(sc.voters[epoch], staker) {
		sc.voters[epoch] = append(sc.voters[epoch], staker)
	}

}

func (sc *KyberStakingContract) GetStake(epoch uint64, staker string) (stake uint64) {
	if _, ok := sc.book[epoch][staker]; !ok{
		prevTemp := sc.book[sc.book2[staker]][staker]
		log.Println(prevTemp.stake + prevTemp.currentStake - prevTemp.withdraw)
		return prevTemp.stake + prevTemp.currentStake - prevTemp.withdraw
	}else {
		log.Println(sc.book[epoch][staker].stake - sc.book[epoch][staker].withdraw)
		return sc.book[epoch][staker].stake - sc.book[epoch][staker].withdraw
	}
}

func (sc *KyberStakingContract) GetDelegatedStake(epoch uint64, staker string) (delegatedStake uint64) {
	sum := uint64(0)
	var temp []string
	if _, ok := sc.book[epoch][staker]; !ok{
		temp = sc.book[sc.book2[staker]][staker].representativeForNextEpoch
	}else {
		temp = sc.book[epoch][staker].representativeFor
	}

	for _, item := range temp {
		if item != staker{
			if _, ok := sc.book[epoch][item]; !ok{
				prevTemp := sc.book[sc.book2[item]][item]
				sum += prevTemp.stake + prevTemp.currentStake - prevTemp.withdraw

			}else{
				sum += sc.book[epoch][item].stake - sc.book[epoch][item].withdraw

			}
		}
	}
	log.Println(sum)
	return sum
}

func (sc *KyberStakingContract) GetRepresentative(epoch uint64, staker string) (poolmaster string) {
	if _, ok := sc.book[epoch][staker]; !ok{
		prevTemp := sc.book[sc.book2[staker]][staker]
		log.Println(prevTemp.delegateToNextEpoch)
		return prevTemp.delegateToNextEpoch
	}else{
		log.Println(sc.book[epoch][staker].delegateTo)
		return sc.book[epoch][staker].delegateTo
	}
}

func (sc *KyberStakingContract) GetReward(epoch uint64, staker string) (percentage float64) {

	sum := uint64(0)
	for _, staker := range sc.voters[epoch] {
		//sum += len(sc.book[epoch][v].vote)
		if sc.book[epoch][staker].delegateTo == staker {
			sum += uint64(len(sc.book[epoch][staker].vote)) * (sc.book[epoch][staker].stake - sc.book[epoch][staker].withdraw)
		}

		for _, item := range sc.book[epoch][staker].representativeFor {
			if item != staker{
				if _, ok := sc.book[epoch][item]; !ok{
					prevTemp := sc.book[sc.book2[item]][item]
					sum += uint64(len(sc.book[epoch][staker].vote)) * (prevTemp.stake + prevTemp.currentStake - prevTemp.withdraw)
				}else{
					sum += uint64(len(sc.book[epoch][staker].vote)) * (sc.book[epoch][item].stake - sc.book[epoch][item].withdraw)
				}
			}
		}
	}
	if sum == 0 {
		log.Println(0)
		return 0
	}else {
		stakerPoint := uint64(0)
		if sc.book[epoch][staker].delegateTo == staker{
			if _, ok := sc.book[epoch][staker]; !ok{
				prevTemp := sc.book[sc.book2[staker]][staker]
				stakerPoint += uint64(len(sc.book[epoch][staker].vote)) * (prevTemp.stake + prevTemp.currentStake - prevTemp.withdraw)
			}else{
				stakerPoint += uint64(len(sc.book[epoch][staker].vote)) * (sc.book[epoch][staker].stake - sc.book[epoch][staker].withdraw)
			}
		}

		for _, item := range sc.book[epoch][staker].representativeFor {
			if item != staker{
				if _, ok := sc.book[epoch][item]; !ok{
					prevTemp := sc.book[sc.book2[item]][item]
					stakerPoint += uint64(len(sc.book[epoch][staker].vote)) * (prevTemp.stake + prevTemp.currentStake - prevTemp.withdraw)
				}else{
					stakerPoint += uint64(len(sc.book[epoch][staker].vote)) * (sc.book[epoch][item].stake - sc.book[epoch][item].withdraw)
				}
			}
		}

		log.Println(stakerPoint,"/", sum,"=", float64(stakerPoint) / float64(sum))
		return float64(stakerPoint) / float64(sum)
	}
}

func (sc *KyberStakingContract) GetPoolReward(epoch uint64, staker string) (percentage float64) {
	stakerPoint := uint64(0)
	var temp []string
	if _, ok := sc.book[epoch][staker]; !ok{
		prevTemp := sc.book[sc.book2[staker]][staker]
		stakerPoint = prevTemp.stake + prevTemp.currentStake - prevTemp.withdraw
		log.Println("====", prevTemp.delegateTo)
		vutien := prevTemp.delegateToNextEpoch
		if _, ok := sc.book[epoch][ vutien ]; !ok {
			temp = sc.book[sc.book2[vutien]][vutien].representativeForNextEpoch
		}else {
			temp = sc.book[epoch][vutien].representativeForNextEpoch
		}
	}else{
		stakerPoint = sc.book[epoch][staker].stake - sc.book[epoch][staker].withdraw
		vutien := sc.book[epoch][staker].delegateTo
		if _, ok := sc.book[epoch][ vutien ]; !ok {
			temp = sc.book[sc.book2[vutien]][vutien].representativeFor
		}else {
			temp = sc.book[epoch][vutien].representativeFor
		}
	}

	poolPoint := uint64(0)

	for _, item := range temp {
		if _, ok := sc.book[epoch][item]; !ok{
			prevTemp := sc.book[sc.book2[item]][item]
			poolPoint += prevTemp.stake + prevTemp.currentStake - prevTemp.withdraw
		}else{
			poolPoint += sc.book[epoch][item].stake - sc.book[epoch][item].withdraw
		}
	}

	log.Println(stakerPoint,"/", poolPoint,"=", float64(stakerPoint)/ float64(poolPoint))
	return float64(stakerPoint)/ float64(poolPoint)
}

