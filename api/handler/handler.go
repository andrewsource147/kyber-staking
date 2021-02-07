package handler

import (
	"github.com/kyber/staking/api/utils"
	"github.com/kyber/staking/contestant"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Handler struct {
	StakingContract *contestant.KyberStakingContract
}

func NewHandler() *Handler {
	return &Handler{
		StakingContract: contestant.NewKyberStakingContract(0, 0),
	}
}

func (self *Handler) GetIndex(c echo.Context) (err error) {
	return c.JSON(http.StatusOK, "OK")
}

// CreateNewContract godoc
// @Summary Create new contract
// @Description Create new contract
// @Accept json
// @Produce plain,json
// @Param request body handler.CreateContractRequest true "Create contract request"
// @Success 201 {object} handler.Response
// @Failure 400 {object} utils.Error
// @Router /newContract [post]
func (self *Handler) CreateNewContract(c echo.Context) (err error) {
	req := &CreateContractRequest{}
	if err = req.bind(c); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	}

	self.StakingContract = contestant.NewKyberStakingContract(req.StartBlock, req.EpochDuration)
	resp := &Response{}
	resp.Data = "Create contract successfully"
	return c.JSON(http.StatusCreated, resp)
}

// Stake godoc
// @Summary Stake
// @Description Stake an amount for a staker at a block
// @Accept json
// @Produce plain,json
// @Param request body handler.StakeRequest true "Stake request"
// @Success 200 {object} handler.Response
// @Failure 400 {object} utils.Error
// @Router /stake [post]
func (self *Handler) Stake(c echo.Context) (err error) {
	req := &StakeRequest{}
	if err = req.bind(c); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	}

	self.StakingContract.Stake(req.Block, req.Amount, req.Staker)
	resp := &Response{}
	resp.Data = "Stake successfully"
	return c.JSON(http.StatusOK, resp)
}

// Withdraw godoc
// @Summary Withdraw
// @Description Staker withdraw an amount at a block
// @Accept json
// @Produce plain,json
// @Param request body handler.StakeRequest true "Withdraw request"
// @Success 200 {object} handler.Response
// @Failure 400 {object} utils.Error
// @Router /withdraw [post]
func (self *Handler) Withdraw(c echo.Context) (err error) {
	req := &StakeRequest{}
	if err = req.bind(c); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	}

	self.StakingContract.Withdraw(req.Block, req.Amount, req.Staker)
	resp := &Response{}
	resp.Data = "Withdraw successfully"
	return c.JSON(http.StatusOK, resp)
}

// Delegate godoc
// @Summary Delegate
// @Description Staker delegate an amount to a representative at a block
// @Accept json
// @Produce plain,json
// @Param request body handler.DelegateRequest true "Delegate request"
// @Success 200 {object} handler.Response
// @Failure 400 {object} utils.Error
// @Router /delegate [post]
func (self *Handler) Delegate(c echo.Context) (err error) {
	req := &DelegateRequest{}
	if err = req.bind(c); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	}

	self.StakingContract.Delegate(req.Block, req.Staker, req.Representative)
	resp := &Response{}
	resp.Data = "Delegate successfully"
	return c.JSON(http.StatusOK, resp)
}

// Vote godoc
// @Summary Vote
// @Description Staker vote for a campain ID at a block
// @Accept json
// @Produce plain,json
// @Param request body handler.VoteRequest true "Vote request"
// @Success 200 {object} handler.Response
// @Failure 400 {object} utils.Error
// @Router /vote [post]
func (self *Handler) Vote(c echo.Context) (err error) {
	req := &VoteRequest{}
	if err = req.bind(c); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	}

	self.StakingContract.Vote(req.Block, req.VoteID, req.Staker)
	resp := &Response{}
	resp.Data = "Vote successfully"
	return c.JSON(http.StatusOK, resp)
}

// GetStake godoc
// @Summary GetStake
// @Description Get amount a staker has staked at an epoch
// @Accept json
// @Produce plain,json
// @Param epoch query uint64 true "Epoch number"
// @Param staker query string true "Staker name"
// @Success 200 {object} handler.Response
// @Failure 400 {object} utils.Error
// @Router /getStake [get]
func (self *Handler) GetStake(c echo.Context) (err error) {
	epochStr := c.QueryParam("epoch")
	staker := c.QueryParam("staker")

	epoch, err := strconv.ParseUint(epochStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	}

	stake := self.StakingContract.GetStake(epoch, staker)
	resp := &Response{}
	resp.Data = stake
	return c.JSON(http.StatusOK, resp)
}

// GetDelegateStake godoc
// @Summary GetDelegateStake
// @Description Get amount a staker has been delegated at an epoch
// @Accept json
// @Produce plain,json
// @Param epoch query uint64 true "Epoch number"
// @Param staker query string true "Staker name"
// @Success 200 {object} handler.Response
// @Failure 400 {object} utils.Error
// @Router /getDelegatedStake [get]
func (self *Handler) GetDelegateStake(c echo.Context) (err error) {
	epochStr := c.QueryParam("epoch")
	staker := c.QueryParam("staker")

	epoch, err := strconv.ParseUint(epochStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	}

	delegatedStake := self.StakingContract.GetDelegatedStake(epoch, staker)
	resp := &Response{}
	resp.Data = delegatedStake
	return c.JSON(http.StatusOK, resp)
}

// GetRepresentative godoc
// @Summary GetRepresentative
// @Description Get representative for a staker at an epoch
// @Accept json
// @Produce plain,json
// @Param epoch query uint64 true "Epoch number"
// @Param staker query string true "Staker name"
// @Success 200 {object} handler.Response
// @Failure 400 {object} utils.Error
// @Router /getRepresentative [get]
func (self *Handler) GetRepresentative(c echo.Context) (err error) {
	epochStr := c.QueryParam("epoch")
	staker := c.QueryParam("staker")

	epoch, err := strconv.ParseUint(epochStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	}

	representative := self.StakingContract.GetRepresentative(epoch, staker)
	resp := &Response{}
	resp.Data = representative
	return c.JSON(http.StatusOK, resp)
}

// GetReward godoc
// @Summary GetReward
// @Description Get reward for a staker at an epoch (in percentage)
// @Accept json
// @Produce plain,json
// @Param epoch query uint64 true "Epoch number"
// @Param staker query string true "Staker name"
// @Success 200 {object} handler.Response
// @Failure 400 {object} utils.Error
// @Router /getReward [get]
func (self *Handler) GetReward(c echo.Context) (err error) {
	epochStr := c.QueryParam("epoch")
	staker := c.QueryParam("staker")

	epoch, err := strconv.ParseUint(epochStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	}

	reward := self.StakingContract.GetReward(epoch, staker)
	resp := &Response{}
	resp.Data = reward
	return c.JSON(http.StatusOK, resp)
}

// GetPoolReward godoc
// @Summary GetPoolReward
// @Description Get reward for a staker from representative at an epoch
// @Accept json
// @Produce plain,json
// @Param epoch query uint64 true "Epoch number"
// @Param staker query string true "Staker name"
// @Success 200 {object} handler.Response
// @Failure 400 {object} utils.Error
// @Router /getPoolReward [get]
func (self *Handler) GetPoolReward(c echo.Context) (err error) {
	epochStr := c.QueryParam("epoch")
	staker := c.QueryParam("staker")

	epoch, err := strconv.ParseUint(epochStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	}

	poolReward := self.StakingContract.GetPoolReward(epoch, staker)
	resp := &Response{}
	resp.Data = poolReward
	return c.JSON(http.StatusOK, resp)
}
