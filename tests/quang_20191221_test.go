package tests

import (
	"github.com/kyber/staking/contestant"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

const (
	Mike = "mike"
	Victor = "victor"
	Loi = "loi"
	Andrew = "andrew"
	Quang = "quang"
	Tien = "tien"
	Trung = "trung"
)

const MaxUint = ^uint64(0)

func TestStaking_Quang_20191221_TC1(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	assert := assert.New(t)
	//2000 1000                   ---> Epoch: 2000 blocks, StartBlock: 1000
	kyberSc := contestant.NewKyberStakingContract(1000, 2000)
	// stake 500 1000 mike
	kyberSc.Stake(500, 1000, Mike)
	// stake 900 500 victor
	kyberSc.Stake(900, 500, Victor)
	// stake 1000 500 mike
	kyberSc.Stake(1000, 500, Mike)
	// getStake 0 mike			---> 0
	assert.Equal(uint64(0), kyberSc.GetStake(0, Mike))
	// getStake 1 mike			---> 1000
	assert.Equal(uint64(1000), kyberSc.GetStake(1, Mike))
	// getStake 1 victor 		---> 500
	assert.Equal(uint64(500), kyberSc.GetStake(1, Victor))
	// getStake 1 loi 			---> 0
	assert.Equal(uint64(0), kyberSc.GetStake(1, Loi))

	// vote 1003 1 mike
	kyberSc.Vote(1001, 1, Mike)
	// vote 1004 1 victor
	kyberSc.Vote(1002, 1, Victor)
	// vote 1004 1 victor
	kyberSc.Vote(1003, 1, Loi)
	// getReward 0 mike			---> mike 1000 + victor 500
	assert.Equal(float64(0), kyberSc.GetReward(0, Mike))
	// getReward 1 mike			---> mike 1000 + victor 500
	assert.Equal(float64(2)/float64(3), kyberSc.GetReward(1, Mike))
	// getReward 1 victor		---> mike 1000 + victor 500
	assert.Equal(float64(1)/float64(3), kyberSc.GetReward(1, Victor))
	// getReward 1 loi			---> 0
	assert.Equal(float64(0), kyberSc.GetReward(1, Loi))

	// delegate 1001 mike loi
	kyberSc.Delegate(1003, Mike, Loi)
	// stake 1002 1000 mike
	kyberSc.Stake(1004, 1000, Mike)
	// getStake 1 mike 			---> 1000
	assert.Equal(uint64(1000), kyberSc.GetStake(1, Mike))
	// getStake 1 loi			---> 0
	assert.Equal(uint64(0), kyberSc.GetStake(1, Loi))
	// getStake 2 mike			---> mike stake +1500 at epoch 1, +500 at epoch 0
	assert.Equal(uint64(2500), kyberSc.GetStake(2, Mike))
	// getStake 2 loi			---> 0
	assert.Equal(uint64(0), kyberSc.GetStake(2, Loi))
	// getDelegatedStake 2 loi	---> 2500 from mike
	assert.Equal(uint64(2500), kyberSc.GetDelegatedStake(2, Loi))
	// getRepresentative 2 mike ---> loi
	assert.Equal(Loi, kyberSc.GetRepresentative(2, Mike))

	// vote 3000 2 victor
	kyberSc.Vote(3000, 2, Victor)
	// vote 3001 2 loi
	kyberSc.Vote(3001, 2, Loi)
	// getReward 2 loi			---> 5/6 	(loi: 2500, victor: 500, mike: 0)
	assert.Equal(float64(5)/float64(6), kyberSc.GetReward(2, Loi))
	// getReward 2 victor		---> 1/6 	(loi: 2500, victor: 500, mike: 0)
	assert.Equal(float64(1)/float64(6), kyberSc.GetReward(2, Victor))
	// getReward 2 mike			---> 0/6 	(loi: 2500, victor: 500, mike: 0)
	assert.Equal(float64(0), kyberSc.GetReward(2, Mike))
	// getPoolReward 2 mike 	---> 1		(loi: 0, mike: 2500)
	assert.Equal(float64(1), kyberSc.GetPoolReward(2, Mike))
	// withdraw 3002 500 mike
	kyberSc.Withdraw(3002, 500, Mike)
	// getReward 2 loi
	assert.Equal(float64(4)/float64(5), kyberSc.GetReward(2, Loi))
	// getReward 2 victor
	assert.Equal(float64(1)/float64(5), kyberSc.GetReward(2, Victor))
}

func TestStaking_Quang_20191221_TC2(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

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
	assert.Equal(float64(0), kyberSc.GetReward(1, "quang"))

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

	kyberSc.Withdraw(24, 1000, "tien")
	assert.Equal(uint64(0), kyberSc.GetStake(2, "tien"))
	assert.Equal(float64(0), kyberSc.GetReward(2, "tien"))

	kyberSc.Vote(26, 2, "loi")
	assert.Equal(float64(0), kyberSc.GetReward(2, "loi"))

	kyberSc.Stake(28, 200, "tien")
	assert.Equal(uint64(0), kyberSc.GetStake(2, "tien"))
	assert.Equal(uint64(200), kyberSc.GetStake(3, "tien"))
	assert.Equal(float64(0), kyberSc.GetReward(2, "tien"))
	assert.Equal(float64(0), kyberSc.GetReward(3, "tien"))

	kyberSc.Delegate(30, "tien", "tien")
	assert.Equal("quang", kyberSc.GetRepresentative(3, "tien"))
	assert.Equal("tien", kyberSc.GetRepresentative(4, "tien"))
	assert.Equal(float64(1), kyberSc.GetPoolReward(3, "tien"))
	assert.Equal(float64(1), kyberSc.GetPoolReward(4, "tien"))

	kyberSc.Delegate(32, "loi", "loi")
	assert.Equal("loi", kyberSc.GetRepresentative(3, "loi"))
	assert.Equal(float64(0), kyberSc.GetPoolReward(3, "loi"))

	kyberSc.Vote(34, 2, "andrew")
	assert.Equal(float64(0), kyberSc.GetReward(3, "andrew"))

	kyberSc.Withdraw(36, 200, "quang")
	assert.Equal(uint64(0), kyberSc.GetStake(3, "quang"))
	assert.Equal(float64(0), kyberSc.GetReward(3, "quang"))

	kyberSc.Vote(38, 6, "mike")
	assert.Equal(float64(0), kyberSc.GetReward(3, "mike"))

	kyberSc.Delegate(40, "trung", "quang")
	assert.Equal("mike", kyberSc.GetRepresentative(4, "trung"))
	assert.Equal("quang", kyberSc.GetRepresentative(5, "trung"))
	assert.Equal(float64(0), kyberSc.GetPoolReward(4, "trung"))

	kyberSc.Stake(42, 400, "andrew")
	assert.Equal(uint64(0), kyberSc.GetStake(4, "andrew"))
	assert.Equal(uint64(400), kyberSc.GetStake(5, "andrew"))
	assert.Equal(float64(0), kyberSc.GetReward(4, "andrew"))
	assert.Equal(float64(0), kyberSc.GetReward(5, "andrew"))

	kyberSc.Withdraw(44, 600, "trung")
	assert.Equal(uint64(0), kyberSc.GetStake(4, "trung"))
	assert.Equal(float64(0), kyberSc.GetReward(4, "trung"))

	kyberSc.Vote(46, 2, "loi")
	assert.Equal(float64(0), kyberSc.GetReward(4, "loi"))

	kyberSc.Delegate(48, "andrew", "victor")
	assert.Equal("loi", kyberSc.GetRepresentative(4, "andrew"))
	assert.Equal("victor", kyberSc.GetRepresentative(5, "andrew"))
	assert.Equal(float64(0), kyberSc.GetPoolReward(4, "andrew"))
	assert.Equal(float64(1), kyberSc.GetPoolReward(5, "andrew"))

	kyberSc.Withdraw(50, 600, "mike")
	assert.Equal(uint64(0), kyberSc.GetStake(5, "mike"))
	assert.Equal(float64(0), kyberSc.GetReward(5, "mike"))

	kyberSc.Withdraw(52, 600, "quang")
	assert.Equal(uint64(0), kyberSc.GetStake(5, "quang"))
	assert.Equal(float64(0), kyberSc.GetReward(5, "quang"))

	kyberSc.Delegate(54, "quang", "loi")
	assert.Equal("quang", kyberSc.GetRepresentative(5, "quang"))
	assert.Equal("loi", kyberSc.GetRepresentative(6, "quang"))
	assert.Equal(float64(0), kyberSc.GetPoolReward(5, "quang"))

	kyberSc.Delegate(56, "quang", "quang")
	assert.Equal("quang", kyberSc.GetRepresentative(5, "quang"))
	assert.Equal("quang", kyberSc.GetRepresentative(6, "quang"))
	assert.Equal(float64(0), kyberSc.GetPoolReward(5, "quang"))

	kyberSc.Withdraw(58, 400, "quang")
	assert.Equal(uint64(0), kyberSc.GetStake(5, "quang"))
	assert.Equal(uint64(0), kyberSc.GetStake(6, "quang"))
	assert.Equal(float64(0), kyberSc.GetReward(5, "quang"))

	kyberSc.Delegate(60, "andrew", "tien")
	assert.Equal("victor", kyberSc.GetRepresentative(6, "andrew"))
	assert.Equal("tien", kyberSc.GetRepresentative(7, "andrew"))
	assert.Equal(uint64(200), kyberSc.GetStake(6, "tien"))
	assert.Equal(uint64(400), kyberSc.GetStake(6, "andrew"))
	assert.Equal(float64(1), kyberSc.GetPoolReward(6, "andrew"))
	assert.Equal(float64(2)/float64(3), kyberSc.GetPoolReward(7, "andrew"))

	kyberSc.Stake(62, 800, "quang")
	assert.Equal(uint64(0), kyberSc.GetStake(6, "quang"))
	assert.Equal(uint64(800), kyberSc.GetStake(7, "quang"))
	assert.Equal(float64(0), kyberSc.GetReward(6, "quang"))
	assert.Equal(float64(0), kyberSc.GetReward(7, "quang"))

	kyberSc.Stake(64, 1000, "trung")
	assert.Equal(uint64(0), kyberSc.GetStake(6, "trung"))
	assert.Equal(float64(0), kyberSc.GetReward(6, "trung"))
	assert.Equal(float64(0), kyberSc.GetReward(7, "trung"))
}

func TestStaking_Quang_20191221_TC3(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	assert := assert.New(t)
	kyberSc := contestant.NewKyberStakingContract(10, 10)
	kyberSc.Stake(0, 100, Mike)
	kyberSc.Stake(1, 100, Victor)
	kyberSc.Stake(2, 100, Loi)
	kyberSc.Stake(3, 100, Andrew)
	kyberSc.Stake(4, 100, Tien)
	kyberSc.Stake(5, 100, Trung)

	kyberSc.Delegate(10, Mike, Victor)
	kyberSc.Delegate(11, Victor, Loi)
	kyberSc.Delegate(12, Loi, Andrew)
	kyberSc.Delegate(13, Andrew, Tien)
	kyberSc.Delegate(14, Tien, Trung)
	kyberSc.Delegate(15, Trung, Mike)

	assert.Equal(Victor, kyberSc.GetRepresentative(2, Mike))
	assert.Equal(Loi, kyberSc.GetRepresentative(2, Victor))
	assert.Equal(Andrew, kyberSc.GetRepresentative(2, Loi))
	assert.Equal(Tien, kyberSc.GetRepresentative(2, Andrew))
	assert.Equal(Trung, kyberSc.GetRepresentative(2, Tien))
	assert.Equal(Mike, kyberSc.GetRepresentative(2, Trung))

	assert.Equal(uint64(100), kyberSc.GetDelegatedStake(2, Mike))
	assert.Equal(uint64(100), kyberSc.GetDelegatedStake(2, Victor))
	assert.Equal(uint64(100), kyberSc.GetDelegatedStake(2, Loi))
	assert.Equal(uint64(100), kyberSc.GetDelegatedStake(2, Andrew))
	assert.Equal(uint64(100), kyberSc.GetDelegatedStake(2, Tien))
	assert.Equal(uint64(100), kyberSc.GetDelegatedStake(2, Trung))

	kyberSc.Vote(20, 1, Victor)
	kyberSc.Vote(21, 1, Loi)
	kyberSc.Vote(22, 1, Andrew)
	kyberSc.Vote(23, 1, Tien)
	kyberSc.Vote(24, 1, Trung)
	kyberSc.Vote(25, 1, Mike)

	assert.Equal(float64(1)/float64(6), kyberSc.GetReward(2, Mike))
	assert.Equal(float64(1)/float64(6), kyberSc.GetReward(2, Victor))
	assert.Equal(float64(1)/float64(6), kyberSc.GetReward(2, Loi))
	assert.Equal(float64(1)/float64(6), kyberSc.GetReward(2, Andrew))
	assert.Equal(float64(1)/float64(6), kyberSc.GetReward(2, Tien))
	assert.Equal(float64(1)/float64(6), kyberSc.GetReward(2, Trung))

	assert.Equal(float64(1), kyberSc.GetPoolReward(2, Mike))
	assert.Equal(float64(1), kyberSc.GetPoolReward(2, Victor))
	assert.Equal(float64(1), kyberSc.GetPoolReward(2, Loi))
	assert.Equal(float64(1), kyberSc.GetPoolReward(2, Andrew))
	assert.Equal(float64(1), kyberSc.GetPoolReward(2, Tien))
	assert.Equal(float64(1), kyberSc.GetPoolReward(2, Trung))

	kyberSc.Delegate(30, Mike, Trung)
	kyberSc.Delegate(31, Victor, Mike)
	kyberSc.Delegate(32, Loi, Victor)
	kyberSc.Delegate(33, Andrew, Loi)
	kyberSc.Delegate(34, Tien, Andrew)
	kyberSc.Delegate(35, Trung, Tien)

	assert.Equal(float64(1), kyberSc.GetPoolReward(4, Mike))
	assert.Equal(float64(1), kyberSc.GetPoolReward(4, Victor))
	assert.Equal(float64(1), kyberSc.GetPoolReward(4, Loi))
	assert.Equal(float64(1), kyberSc.GetPoolReward(4, Andrew))
	assert.Equal(float64(1), kyberSc.GetPoolReward(4, Tien))
	assert.Equal(float64(1), kyberSc.GetPoolReward(4, Trung))
}

func TestStaking_Quang_20191221_TC4(t *testing.T) {
	assert := assert.New(t)
	kyberSc := contestant.NewKyberStakingContract(10, 10)
	kyberSc.Stake(0, 100, Mike)
	kyberSc.Stake(1, 100, Victor)
	kyberSc.Stake(2, 100, Loi)
	kyberSc.Stake(3, 100, Andrew)

	kyberSc.Delegate(10, Mike, Victor)
	kyberSc.Delegate(11, Loi, Victor)
	kyberSc.Delegate(12, Andrew, Victor)

	assert.Equal(Victor, kyberSc.GetRepresentative(2, Mike))
	assert.Equal(Victor, kyberSc.GetRepresentative(2, Victor))
	assert.Equal(Victor, kyberSc.GetRepresentative(2, Loi))
	assert.Equal(Victor, kyberSc.GetRepresentative(2, Andrew))

	kyberSc.Vote(20, 1, Victor)
	assert.Equal(float64(1), kyberSc.GetReward(2, Victor))
	assert.Equal(float64(1)/float64(4), kyberSc.GetPoolReward(2, Mike))
	assert.Equal(float64(1)/float64(4), kyberSc.GetPoolReward(2, Victor))
	assert.Equal(float64(1)/float64(4), kyberSc.GetPoolReward(2, Loi))
	assert.Equal(float64(1)/float64(4), kyberSc.GetPoolReward(2, Andrew))

	assert.Equal(uint64(0), kyberSc.GetDelegatedStake(2, Mike))
	assert.Equal(uint64(300), kyberSc.GetDelegatedStake(2, Victor))
	assert.Equal(uint64(0), kyberSc.GetDelegatedStake(2, Loi))
	assert.Equal(uint64(0), kyberSc.GetDelegatedStake(2, Andrew))

	kyberSc.Delegate(21, Victor, Trung)
	assert.Equal(uint64(300), kyberSc.GetDelegatedStake(3, Victor))
	assert.Equal(uint64(100), kyberSc.GetDelegatedStake(3, Trung))
	assert.Equal(float64(1), kyberSc.GetPoolReward(3, Victor))
	assert.Equal(float64(0), kyberSc.GetPoolReward(3, Trung))
}
