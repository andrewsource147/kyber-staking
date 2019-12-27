package tests

import (
	"log"
	"testing"

	"github.com/kyber/staking/contestant"
	"github.com/stretchr/testify/assert"
)

func TestMultiLevelDelegate(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//10 10                   ---> Epoch: 10 blocks, StartBlock: 10
	kyberSc := contestant.NewKyberStakingContract(10, 10)

	// stake 1 200 mike           ---> mike stake 200 KNC ở epoch 0
	kyberSc.Stake(1, 200, "mike")
	// stake 1 300 loi           ---> loi stake 300 KNC ở epoch 0
	kyberSc.Stake(1, 300, "loi")
	// stake 1 500 victor           ---> victor stake 500 KNC ở epoch 0
	kyberSc.Stake(1, 500, "victor")

	// delegate 5 mike loi    ---> mike delegate hết tiền cho loi
	kyberSc.Delegate(5, "mike", "loi")
	// delegate 5 loi victor    ---> loi delegate hết tiền cho victor
	kyberSc.Delegate(5, "loi", "victor")
	// delegate 5 victor mike    ---> victor delegate hết tiền cho mike
	kyberSc.Delegate(5, "victor", "mike")

	// getStake 1 mike             ---> 200 (dù mike đã delegate hết tiền cho loi nhưng vẫn phải trả về số lượng stake hiện tại)
	assert.Equal(t, kyberSc.GetStake(1, "mike"), uint64(200))
	// getDelegatedStake 1 mike  ---> 500 (victor delegate cho mike, ko bao gồm lượng đc delegate từ người khác)
	assert.Equal(t, kyberSc.GetDelegatedStake(1, "mike"), uint64(500))
	// getDelegatedStake 1 loi  ---> 500 (mike delegate cho loi, ko bao gồm lượng đc delegate từ người khác)
	assert.Equal(t, kyberSc.GetDelegatedStake(1, "loi"), uint64(200))
	// getDelegatedStake 1 victor  ---> 500 (loi delegate cho victor, ko bao gồm lượng đc delegate từ người khác)
	assert.Equal(t, kyberSc.GetDelegatedStake(1, "victor"), uint64(300))

	// vote 10 1 mike           ---> mike vote 1 campaign ở epoch 1
	kyberSc.Vote(10, 1, "mike")
	// GetReward 1 mike            ---> 0    (mike delegate hết tiền cho loi nhưng victor lại delegate cho mike)
	assert.Equal(t, kyberSc.GetReward(1, "mike"), float64(1))

	// vote 10 1 loi           ---> loi vote 1 campaign ở epoch 1
	kyberSc.Vote(10, 1, "loi")
	// GetReward 1 mike            ---> 0.71    (mike delegate hết tiền cho loi nhưng lại đc victor delegate)
	assert.Equal(t, kyberSc.GetReward(1, "mike"), float64(500)/float64(700))
	// GetReward 1 loi            ---> 0.28    (loi delegate hết tiền cho victor nhưng lai đc mike delegate)
	assert.Equal(t, kyberSc.GetReward(1, "loi"), float64(200)/float64(700))
	// GetReward 1 victor            ---> 0
	assert.Equal(t, kyberSc.GetReward(1, "victor"), float64(0))

	// vote 10 1 victor           ---> victor vote 1 campaign ở epoch 1
	kyberSc.Vote(10, 1, "victor")
	// GetReward 1 mike            ---> 0.5
	assert.Equal(t, kyberSc.GetReward(1, "mike"), float64(500)/float64(1000))
	// GetReward 1 loi            ---> 0.2
	assert.Equal(t, kyberSc.GetReward(1, "loi"), float64(200)/float64(1000))
	// GetReward 1 victor            ---> 0.3
	assert.Equal(t, kyberSc.GetReward(1, "victor"), float64(300)/float64(1000))
}

func TestExcessWithdrawalAfterDelegateToOther(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	kyberSc := contestant.NewKyberStakingContract(10, 10)

	kyberSc.Stake(1, 10, "mike")
	kyberSc.Stake(1, 20, "mike")
	kyberSc.Stake(1, 20, "mike")
	kyberSc.Stake(1, 120, "loi")
	kyberSc.Stake(1, 150, "victor")

	// withdraw 2 200 mike     ---> vì vượt quá lượng stake hiện tại nên lệnh này bị bỏ qua
	kyberSc.Withdraw(2, 200, "mike")
	// withdraw 2 20 mike
	kyberSc.Withdraw(2, 20, "mike")
	assert.Equal(t, kyberSc.GetStake(1, "mike"), uint64(30))

	// delegate 5 mike victor    ---> mike delegate hết tiền cho victor
	kyberSc.Delegate(5, "mike", "victor")
	// getPoolReward 1 mike       ---> (mike có 30 KNC và victor có 150 KNC)
	assert.Equal(t, kyberSc.GetPoolReward(1, "mike"), float64(30)/float64(180))

	// withdraw 2 20 mike
	kyberSc.Withdraw(10, 20, "mike")
	// getPoolReward 1 mike       ---> (mike h chỉ có 10 KNC và victor có 150 KNC)
	assert.Equal(t, kyberSc.GetPoolReward(1, "mike"), float64(10)/float64(160))

	// stake 11 100 mike     ---> (mike stake thêm 100 KNC)
	kyberSc.Stake(11, 100, "mike")
	// getPoolReward 1 mike       ---> (mike vẫn chỉ có 10 KNC và victor có 150 KNC do 100 KNC stake thêm sẽ tính vào epoch sau)
	assert.Equal(t, kyberSc.GetPoolReward(1, "mike"), float64(10)/float64(160))
	// getPoolReward 2 mike       ---> (mike bây h có 110 KNC và victor có 150 KNC)
	assert.Equal(t, kyberSc.GetPoolReward(2, "mike"), float64(110)/float64(260))
}

func TestNoVoteNoReward(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//2000 1000                   ---> Epoch: 2000 blocks, StartBlock: 1000
	kyberSc := contestant.NewKyberStakingContract(1000, 2000)

	// stake 999 1000 mike           ---> mike stake 1000 KNC ở epoch 0
	kyberSc.Stake(999, 1000, "mike")
	// stake 999 1000 loi           ---> loi stake 1000 KNC ở epoch 0
	kyberSc.Stake(999, 1000, "loi")
	// stake 999 1000 victor           ---> victor stake 1000 KNC ở epoch 0
	kyberSc.Stake(999, 1000, "victor")

	// vote 1000 1 mike           ---> mike vote 1 campaign ở epoch 1
	kyberSc.Vote(1000, 1, "mike")
	// getReward 1 mike            ---> 1.0    (vì victor + loi chưa vote cho campaign nào cả)
	assert.Equal(t, kyberSc.GetReward(1, "mike"), float64(1))

	// vote 1000 1 loi           ---> loi vote 1 campaign ở epoch 1
	kyberSc.Vote(1000, 1, "loi")
	// getReward 1 mike            ---> 0.5    (vì loi đã bắt đầu vote và victor thì chưa)
	assert.Equal(t, kyberSc.GetReward(1, "mike"), float64(1)/float64(2))
	// getReward 1 loi            ---> 0.5    (vì mike đã vote và victor chưa vote)
	assert.Equal(t, kyberSc.GetReward(1, "loi"), float64(1)/float64(2))
	// getReward 1 victor            ---> 0    (victor lười ko vote nên ko có reward ở epoch này)
	assert.Equal(t, kyberSc.GetReward(1, "victor"), float64(0))

	// getReward 2 mike            ---> 0    (ở epoch 2 chưa ai vote nên ko có reward)
	assert.Equal(t, kyberSc.GetReward(2, "mike"), float64(0))
	// getReward 2 loi            ---> 0    (ở epoch 2 chưa ai vote nên ko có reward)
	assert.Equal(t, kyberSc.GetReward(2, "loi"), float64(0))
	// getReward 2 victor            ---> 0    (ở epoch 2 chưa ai vote nên ko có reward)
	assert.Equal(t, kyberSc.GetReward(2, "victor"), float64(0))
}

func TestRevokeDelegate(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	kyberSc := contestant.NewKyberStakingContract(10, 10)

	kyberSc.Stake(1, 50, "mike")
	kyberSc.Stake(1, 150, "loi")
	kyberSc.Stake(1, 200, "victor")

	kyberSc.Delegate(5, "mike", "loi")
	assert.Equal(t, kyberSc.GetPoolReward(1, "mike"), float64(50)/float64(200))

	kyberSc.Delegate(10, "mike", "victor")
	assert.Equal(t, kyberSc.GetPoolReward(1, "mike"), float64(50)/float64(200))
	assert.Equal(t, kyberSc.GetPoolReward(2, "mike"), float64(50)/float64(250))

	kyberSc.Delegate(11, "mike", "mike")
	assert.Equal(t, kyberSc.GetPoolReward(1, "mike"), float64(50)/float64(200))
	assert.Equal(t, kyberSc.GetPoolReward(2, "mike"), float64(1))
}

func TestReceivedRewardAfterVoteMultipleTimes(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	kyberSc := contestant.NewKyberStakingContract(10, 10)

	kyberSc.Stake(1, 100, "mike")
	kyberSc.Stake(1, 200, "loi")
	kyberSc.Stake(1, 700, "victor")

	kyberSc.Vote(10, 1, "mike")
	kyberSc.Vote(10, 1, "loi")
	kyberSc.Vote(10, 1, "victor")

	assert.Equal(t, kyberSc.GetReward(1, "mike"), float64(100)/float64(1000))
	assert.Equal(t, kyberSc.GetReward(1, "loi"), float64(200)/float64(1000))
	assert.Equal(t, kyberSc.GetReward(1, "victor"), float64(700)/float64(1000))

	kyberSc.Vote(11, 2, "mike")
	// vote 12 2 mike            ---> vote cùng 1 block ko đc tính
	kyberSc.Vote(12, 2, "mike")

	assert.Equal(t, kyberSc.GetReward(1, "mike"), float64(200)/float64(1100))
	assert.Equal(t, kyberSc.GetReward(1, "loi"), float64(200)/float64(1100))
	assert.Equal(t, kyberSc.GetReward(1, "victor"), float64(700)/float64(1100))

	kyberSc.Vote(11, 3, "loi")
	kyberSc.Vote(11, 3, "victor")
	kyberSc.Vote(12, 4, "victor")

	assert.Equal(t, kyberSc.GetReward(1, "mike"), float64(200)/float64(2700))
	assert.Equal(t, kyberSc.GetReward(1, "loi"), float64(400)/float64(2700))
	assert.Equal(t, kyberSc.GetReward(1, "victor"), float64(2100)/float64(2700))
}
