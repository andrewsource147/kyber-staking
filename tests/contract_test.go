package tests

import (
	"log"
	"testing"

	"github.com/kyber/staking/contestant"
	// "github.com/stretchr/testify/assert"
)

func TestSaking(t *testing.T) {
	startBlock := uint64(2000)
	epochDuration := uint64(1000)
	kyberSc := contestant.NewKyberStakingContract(startBlock, epochDuration)

	kyberSc.Stake(uint64(900), uint64(500), "mike")
	kyberSc.GetStake(uint64(0), "mike")
	kyberSc.GetStake(uint64(1), "mike")
	log.Println("staking done")
}
