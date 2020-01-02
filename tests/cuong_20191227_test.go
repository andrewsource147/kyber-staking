package tests

import (
	"log"
	"testing"

	"github.com/kyber/staking/contestant"
	"github.com/stretchr/testify/assert"
)

func Test_Cuong_20191227_1(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	kyberSc := contestant.NewKyberStakingContract(10, 10)
	kyberSc.Stake(1, 50, "mike")
	kyberSc.Stake(21, 50, "mike") // epoch 2

	assert.Equal(t, uint64(0), kyberSc.GetStake(0, "mike"))
	assert.Equal(t, uint64(50), kyberSc.GetStake(1, "mike"))
	assert.Equal(t, uint64(50), kyberSc.GetStake(2, "mike"))
	assert.Equal(t, uint64(100), kyberSc.GetStake(3, "mike"))

	// assert.Equal(t, float64(0), kyberSc.GetReward(1, "mike"))
	kyberSc.Vote(21, 1, "mike")
	kyberSc.Vote(21, 1, "mike")
	kyberSc.Vote(21, 2, "mike")

	assert.Equal(t, float64(0), kyberSc.GetReward(1, "mike"))
	assert.Equal(t, float64(1), kyberSc.GetReward(2, "mike"))

}

func Test_Cuong_20191227_2(t *testing.T) {
	assert := assert.New(t)
	kyberSc := contestant.NewKyberStakingContract(10, 10)
	kyberSc.Vote(0, 6, "quang")
	assert.Equal(float64(0), kyberSc.GetReward(0, "quang"))

	kyberSc.Stake(2, 800, "quang")
	assert.Equal(uint64(0), kyberSc.GetStake(0, "quang"))
	assert.Equal(uint64(800), kyberSc.GetStake(1, "quang"))
	assert.Equal(float64(0), kyberSc.GetReward(0, "quang"))
	assert.Equal(float64(0), kyberSc.GetReward(1, "quang"))

	kyberSc.Vote(4, 3, "trung")
	assert.Equal(float64(0), kyberSc.GetReward(0, "trung"))

	kyberSc.Delegate(6, "trung", "mike")
	assert.Equal("trung", kyberSc.GetRepresentative(0, "trung"))
	assert.Equal("mike", kyberSc.GetRepresentative(1, "trung"))
	assert.Equal(float64(0), kyberSc.GetPoolReward(0, "trung"))

	kyberSc.Vote(8, 3, "quang")
	assert.Equal(float64(0), kyberSc.GetReward(0, "quang"))

	kyberSc.Vote(10, 5, "loi")
	assert.Equal(float64(0), kyberSc.GetReward(1, "loi"))

	kyberSc.Delegate(12, "andrew", "loi")
	assert.Equal("andrew", kyberSc.GetRepresentative(1, "andrew"))
	assert.Equal("loi", kyberSc.GetRepresentative(2, "andrew"))
	assert.Equal(float64(0), kyberSc.GetPoolReward(1, "andrew"))

	kyberSc.Vote(14, 6, "mike")
	assert.Equal(float64(0), kyberSc.GetReward(1, "mike"))

	kyberSc.Vote(16, 4, "andrew")
	assert.Equal(float64(0), kyberSc.GetReward(1, "andrew"))

	kyberSc.Delegate(18, "tien", "quang")
	assert.Equal("tien", kyberSc.GetRepresentative(1, "tien"))
	assert.Equal("quang", kyberSc.GetRepresentative(2, "tien"))
	assert.Equal(float64(0), kyberSc.GetPoolReward(1, "tien"))

	kyberSc.Delegate(20, "mike", "trung")
	assert.Equal("mike", kyberSc.GetRepresentative(2, "mike"))
	assert.Equal("trung", kyberSc.GetRepresentative(3, "mike"))
	assert.Equal(float64(0), kyberSc.GetPoolReward(2, "mike"))

	kyberSc.Withdraw(22, 800, "quang")
	assert.Equal(uint64(0), kyberSc.GetStake(2, "quang"))
	assert.Equal(uint64(0), kyberSc.GetStake(3, "quang"))
	assert.Equal(float64(0), kyberSc.GetReward(2, "quang"))
	assert.Equal(float64(0), kyberSc.GetReward(3, "quang"))
}

func Test_Cuong_20191227_3(t *testing.T) {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 2 1
	kyberSc := contestant.NewKyberStakingContract(1, 2)

	// stake 2 216 A1
	kyberSc.Stake(2, 216, "A1")
	// stake 10 12 A1
	kyberSc.Stake(10, 12, "A1")
	// withdraw 10 23 A1
	kyberSc.Withdraw(10, 23, "A1")
	// getStake 5 A1
	assert.Equal(t, uint64(205), kyberSc.GetStake(5, "A1"))
}
