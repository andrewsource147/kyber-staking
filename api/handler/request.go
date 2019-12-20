package handler

import (
	"github.com/labstack/echo/v4"
	"log"
)

type CreateContractRequest struct {
	StartBlock    uint64 `json:"start_block" validate:"required,numeric" example:"2000"`
	EpochDuration uint64 `json:"epoch_duration" validate:"required,numeric" example:"1000"`
}

func (self *CreateContractRequest) bind(c echo.Context) error {
	if err := c.Bind(self); err != nil {
		log.Println(err)
		return err
	}

	if err := c.Validate(self); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

type StakeRequest struct {
	Block uint64 `json:"block" validate:"required,numeric" example:"3000"`
	Amount uint64 `json:"amount" validate:"required,numeric" example:"500"`
	Staker string `json:"staker" validate:"required" example:"Alice"`
}

func (self *StakeRequest) bind(c echo.Context) error {
	if err := c.Bind(self); err != nil {
		log.Println(err)
		return err
	}

	if err := c.Validate(self); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

type DelegateRequest struct {
	Block uint64 `json:"block" validate:"required,numeric" example:"3000"`
	Staker string `json:"staker" validate:"required" example:"Alice"`
	Representative string `json:"representative" validate:"required" example:"Bob"`
}

func (self *DelegateRequest) bind(c echo.Context) error {
	if err := c.Bind(self); err != nil {
		log.Println(err)
		return err
	}

	if err := c.Validate(self); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

type VoteRequest struct {
	Block uint64 `json:"block" validate:"required,numeric" example:"3000"`
	VoteID uint64  `json:"vote_id" validate:"required,numeric" example:"1"`
	Staker string `json:"staker" validate:"required" example:"Alice"`
}

func (self *VoteRequest) bind(c echo.Context) error {
	if err := c.Bind(self); err != nil {
		log.Println(err)
		return err
	}

	if err := c.Validate(self); err != nil {
		log.Println(err)
		return err
	}

	return nil
}