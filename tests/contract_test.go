package tests

import (
	"log"
	"testing"

	"github.com/kyber/staking/contestant"
	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/assert"
)

func TestStaking(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//2000 1000                   ---> Epoch: 2000 blocks, StartBlock: 1000
	kyberSc := contestant.NewKyberStakingContract(1000, 2000)
	// stake 900 500 mike          --->
	kyberSc.Stake(900, 500, "mike")
	// getStake 0 mike             ---> 0      (mike không có KNC nào đc tính ở epoch 0)
	assert.Equal(t, kyberSc.GetStake(0, "mike"), uint64(0))
	// getStake 1 mike             ---> 500    (mike có 500 KNC đc tính ở epoch 1)
	assert.Equal(t, kyberSc.GetStake(1, "mike"), uint64(500))
	// stake 1000 500 victor       --->
	kyberSc.Stake(1000, 500, "victor")
	// stake 1001 1200 loi         --->
	kyberSc.Stake(1001, 1200, "loi")
	// stake 1001 500 mike         --->
	kyberSc.Stake(1001, 500, "mike")
	// withdraw 1002 200 mike      --->
	kyberSc.Withdraw(1002, 200, "mike")
	// getStake 1 mike             ---> 500    (mike vẫn chỉ có 500 KNC đc tính ở epoch 1)
	assert.Equal(t, kyberSc.GetStake(1, "mike"), uint64(500))
	// getStake 2 mike             ---> 800    (mike có 800 KNC đc tính ở epoch 2)
	assert.Equal(t, kyberSc.GetStake(2, "mike"), uint64(800))
	// vote 1002 1 mike            --->
	kyberSc.Vote(1002, 1, "mike")
	// vote 1003 1 victor          --->
	kyberSc.Vote(1004, 1, "victor")
	// vote 1004 1 loi             --->
	kyberSc.Vote(1004, 1, "loi")
	// getReward 1 mike            ---> 1.0    (vì victor + loi ko có KNC đc tính ở epoch 1)
	assert.Equal(t, kyberSc.GetReward(1, "mike"), float64(1))
	// getReward 1 victor          ---> 0.0    (vì victor ko có KNC đc tính ở epoch 1)
	assert.Equal(t, kyberSc.GetReward(1, "victor"), float64(0))
	// getReward 1 loi             ---> 0.0    (vì loi ko có KNC đc tính ở epoch 1)
	assert.Equal(t, kyberSc.GetReward(1, "loi"), float64(0))

	// delegate 1005 loi victor    --->
	kyberSc.Delegate(1005, "loi", "victor")
	// getDelegatedStake 1 victor  ---> 0      (delegate chỉ có hiệu lực ở epoch sau)
	assert.Equal(t, kyberSc.GetDelegatedStake(1, "victor"), uint64(0))
	// getDelegatedStake 2 victor  ---> 1200   (stake loi đã delegate cho victor)
	assert.Equal(t, kyberSc.GetDelegatedStake(2, "victor"), uint64(1200))
	// getRepresentative 1 loi     ---> loi    (delegate chỉ có hiệu lực ở epoch sau)
	assert.Equal(t, kyberSc.GetRepresentative(1, "loi"), "loi")
	// getRepresentative 2 loi     ---> victor (loi delegate cho victor ở epoch trước đó)
	assert.Equal(t, kyberSc.GetRepresentative(2, "loi"), "victor")

	// vote 3000 2 mike           --->
	kyberSc.Vote(3000, 2, "mike")
	// vote 3001 2 victor         --->
	kyberSc.Vote(3001, 2, "victor")
	// getReward 2 mike           ---> 0.32     (mike: 800, loi+victor: 1700)
	assert.Equal(t, kyberSc.GetReward(2, "mike"), float64(0.32))
	// getReward 2 victor         ---> 0.68     (loi+victor: 1700, mike: 800)
	assert.Equal(t, kyberSc.GetReward(2, "victor"), float64(0.68))
	// getPoolReward 2 loi        ---> 12/17    (loi: 1200, victor: 500)
	assert.Equal(t, kyberSc.GetPoolReward(2, "loi"), float64(12)/float64(17))
	// getPoolReward 2 victor     ---> 5/17     (victor: 500, loi: 1200)
	assert.Equal(t, kyberSc.GetPoolReward(2, "victor"), float64(5)/float64(17))
	// getPoolReward 2 mike       ---> 1        (mike ko delegate cho ai, và ko ai delegate cho mike)
	assert.Equal(t, kyberSc.GetPoolReward(2, "mike"), float64(1))

	// vote 3002 3 mike           --->
	kyberSc.Vote(3002, 3, "mike")
	// getReward 2 mike           ---> 16/33    (mike: 800*2, loi+victor: 1700)
	assert.Equal(t, kyberSc.GetReward(2, "mike"), float64(16)/float64(33))
	// getReward 2 victor         ---> 17/33    (loi+victor: 1700, mike: 800*2)
	assert.Equal(t, kyberSc.GetReward(2, "victor"), float64(17)/float64(33))
	// withdraw 3003 200 mike     --->
	kyberSc.Withdraw(3003, 200, "mike")
	// getReward 2 mike           ---> 12/29    (mike: 600*2, loi+victor: 1700)
	assert.Equal(t, kyberSc.GetReward(2, "mike"), float64(12)/float64(29))
	// getReward 2 victor         ---> 17/29    (loi+victor: 1700, mike: 600*2)
	assert.Equal(t, kyberSc.GetReward(2, "victor"), float64(17)/float64(29))

}
