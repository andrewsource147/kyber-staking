package contestant

import (
	"log"
)

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
	book2 map[string][]uint64
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

func NewKyberStakingContract(startBlock uint64, epochDuration uint64) *KyberStakingContract {
	return &KyberStakingContract{startBlock, epochDuration, 0,
		make(map[uint64]map[string]*StakerStruct), make(map[string][]uint64), make(map[uint64][]string)}
}

func (sc *KyberStakingContract) GetEpoch(block uint64) uint64{
	if block < sc.startBlock {
		return 0
	}else{
		return (block-sc.startBlock)/sc.epochDuration + 1
	}
}

func (sc *KyberStakingContract) findActiveTarget(epoch uint64, staker string) (uint64, *StakerStruct, uint64){
	if _,ok1 := sc.book2[staker]; ok1 {
		for i := len(sc.book2[staker])-1; i >= 0; i-- {
			item:=sc.book2[staker][i]
			if item==epoch {
				return item, sc.book[item][staker], 0
			} else if item < epoch {return item, sc.book[item][staker], 1}
		}
	}
	return 0, nil, 2
}

func (sc *KyberStakingContract) InitEpoch(epoch uint64, staker string) {// 3M records
	if _, ok := sc.book[epoch]; !ok {
		sc.book[epoch] = make(map[string]*StakerStruct) //1.9s
		sc.voters[epoch] = []string{}					//2.9s
	}
	_, activeTarget, status:= sc.findActiveTarget(epoch, staker)
	if status == 1 {// 3M staker = 7.3
		sc.book[epoch][staker] = &StakerStruct{0, 0, 0, []uint64{}, []string{staker}, []string{staker}, staker, staker}
		sc.book[epoch][staker].representativeForNextEpoch = append([]string{}, activeTarget.representativeForNextEpoch...)
		sc.book[epoch][staker].representativeFor = append([]string{} , activeTarget.representativeForNextEpoch...)
		sc.book[epoch][staker].delegateToNextEpoch = activeTarget.delegateToNextEpoch
		sc.book[epoch][staker].delegateTo = activeTarget.delegateToNextEpoch
		sc.book[epoch][staker].stake = activeTarget.stake + activeTarget.currentStake - activeTarget.withdraw
		sc.book2[staker] = append(sc.book2[staker], epoch)
	} else if status == 2 {// 1 staker = 6.5
		sc.book[epoch][staker] = &StakerStruct{0, 0, 0, []uint64{}, []string{staker}, []string{staker}, staker, staker}
		sc.book2[staker] = []uint64{epoch}
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
	sc.InitEpoch(epoch, staker) //1,2 staker 6.5
	sc.InitEpoch(epoch, representative)//1,2 staker 9s
	sc.InitEpoch(epoch, sc.book[epoch][staker].delegateToNextEpoch)//1,2 staker 10s
	//log.Println("---", staker, sc.book[epoch][staker])
	//log.Println("---", representative, sc.book[epoch][representative])
	if staker == representative {
		delegateToNextEpoch := sc.book[epoch][staker].delegateToNextEpoch
		sc.book[epoch][ delegateToNextEpoch ].representativeForNextEpoch = append([]string{}, remove(sc.book[epoch][ delegateToNextEpoch ].representativeForNextEpoch, staker)...)
		if !containsString(sc.book[epoch][staker].representativeForNextEpoch, staker) {
			sc.book[epoch][staker].representativeForNextEpoch = append(sc.book[epoch][staker].representativeForNextEpoch, staker)
		}
		sc.book[epoch][staker].delegateToNextEpoch = staker
	}else {
		delegateToNextEpoch := sc.book[epoch][staker].delegateToNextEpoch
		if delegateToNextEpoch != representative {
			if !containsString(sc.book[epoch][representative].representativeForNextEpoch, staker) {
				sc.book[epoch][representative].representativeForNextEpoch = append(sc.book[epoch][representative].representativeForNextEpoch, staker)
			}
			sc.book[epoch][ delegateToNextEpoch ].representativeForNextEpoch = append([]string{}, remove(sc.book[epoch][ delegateToNextEpoch ].representativeForNextEpoch, staker)...)

			sc.book[epoch][staker].delegateToNextEpoch = representative
			sc.book[epoch][staker].representativeForNextEpoch = append([]string{},  remove(sc.book[epoch][staker].representativeForNextEpoch, staker)...)
		}
	}

	//log.Println("---", staker, sc.book[epoch][staker])
	//log.Println("---", representative, sc.book[epoch][representative])
	//log.Println("-----------------------------------")
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

	_, activeTarget, status:= sc.findActiveTarget(epoch, staker)
	if status == 1{
		return activeTarget.stake + activeTarget.currentStake - activeTarget.withdraw
	}else if status == 0 {
		return activeTarget.stake - activeTarget.withdraw
	}
	return 0
}

func (sc *KyberStakingContract) GetDelegatedStake(epoch uint64, staker string) (delegatedStake uint64) {
	sum := uint64(0)
	var temp []string
	_, activeTarget, status:= sc.findActiveTarget(epoch, staker)
	if status == 2 { return 0 }
	if status == 1 {
		temp = activeTarget.representativeForNextEpoch
	} else if status == 0 {
		temp = activeTarget.representativeFor
	}

	for _, item := range temp {
		if item != staker{
			_, activeTarget, status:= sc.findActiveTarget(epoch, item)
			if status == 1 {
				sum += activeTarget.stake + activeTarget.currentStake - activeTarget.withdraw
			} else if status == 0 {
				sum += activeTarget.stake - activeTarget.withdraw
			}
		}
	}
	//log.Println(sum)
	return sum
}

func (sc *KyberStakingContract) GetRepresentative(epoch uint64, staker string) (poolmaster string) {
	_, activeTarget, status:= sc.findActiveTarget(epoch, staker)
	if status == 1 { return activeTarget.delegateToNextEpoch }
	if status == 0 { return activeTarget.delegateTo }
	return staker
}

func (sc *KyberStakingContract) GetReward(epoch uint64, staker string) (percentage float64) {

	sum := uint64(0)
	for _, staker := range sc.voters[epoch] {
		//if sc.book[epoch][staker].delegateTo == staker {
		//	sum += uint64(len(sc.book[epoch][staker].vote)) * (sc.book[epoch][staker].stake - sc.book[epoch][staker].withdraw)
		//}

		for _, item := range sc.book[epoch][staker].representativeFor {
			_, activeTarget, status:= sc.findActiveTarget(epoch, item)
			if status == 1 {
				sum += uint64(len(sc.book[epoch][staker].vote)) * (activeTarget.stake + activeTarget.currentStake - activeTarget.withdraw)
			} else if status == 0 {
				sum += uint64(len(sc.book[epoch][staker].vote)) * (activeTarget.stake - activeTarget.withdraw)
			}
		}
	}
	if sum == 0 {
		//log.Println(0)
		return 0
	}else {
		stakerPoint := uint64(0)
		if containsString(sc.voters[epoch], staker){
			for _, item := range sc.book[epoch][staker].representativeFor {
				_, activeTarget, status:= sc.findActiveTarget(epoch, item)
				if status == 1 {
					stakerPoint += uint64(len(sc.book[epoch][staker].vote)) * (activeTarget.stake + activeTarget.currentStake - activeTarget.withdraw)
				} else if status == 0 {
					stakerPoint += uint64(len(sc.book[epoch][staker].vote)) * (activeTarget.stake - activeTarget.withdraw)
				}
			}
		}

		//log.Println(stakerPoint,"/", sum,"=", float64(stakerPoint) / float64(sum))
		return float64(stakerPoint) / float64(sum)
	}
}

func (sc *KyberStakingContract) GetPoolReward(epoch uint64, staker string) (percentage float64) {
	stakerPoint := uint64(0)
	var temp []string
	var ee uint64
	activeEpoch, activeTarget, status:= sc.findActiveTarget(epoch, staker)

	if status == 2 {return 0}
	if status == 1{
		stakerPoint = activeTarget.stake + activeTarget.currentStake - activeTarget.withdraw
		_, a, s:= sc.findActiveTarget(activeEpoch+1, activeTarget.delegateToNextEpoch)
		ee = activeEpoch+1
		if s==1{
			temp = a.representativeForNextEpoch
		}else if s== 0{
			temp = a.representativeFor
		}
	}else {
		stakerPoint = activeTarget.stake - activeTarget.withdraw
		_, a, s:= sc.findActiveTarget(activeEpoch, sc.book[epoch][staker].delegateTo)
		ee = activeEpoch
		if s==1{
			temp = a.representativeForNextEpoch
		}else if s== 0{
			temp = a.representativeFor
		}
	}

	poolPoint := uint64(0)

	for _, item := range temp {
		_, activeTarget, status:= sc.findActiveTarget(ee, item)
		if status == 1 {
			poolPoint += activeTarget.stake + activeTarget.currentStake - activeTarget.withdraw
		} else if status == 0 {
			poolPoint += activeTarget.stake - activeTarget.withdraw
		}
	}

	//log.Println(stakerPoint,"/", poolPoint,"=", float64(stakerPoint)/ float64(poolPoint))
	if poolPoint ==0 {
		if stakerPoint == 0{
			return 0
		}else {
			return 1
		}
	}else {
		return float64(stakerPoint)/ float64(poolPoint)
	}
	return float64(stakerPoint)/ float64(poolPoint)
}

