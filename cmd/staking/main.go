package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"sort"
	"stake-go/contestant"
	"strconv"
	"strings"
	"time"

	"gopkg.in/urfave/cli.v1"
)

func main() {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	//set log for server
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	app := cli.NewApp()
	app.Name = "Kyber Staking Contract"
	app.Usage = "For testing kyber staking contract"
	app.Version = "1.0.0"

	app.Flags = []cli.Flag{}

	app.Commands = []cli.Command{}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Action = cmdMain

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func cmdMain(ctx *cli.Context) error {
	time1 := time.Now()
	var startBlock uint64
	var epochDuration uint64
	//var startBlock, epochDuration uint64
	kyberSc := contestant.NewKyberStakingContract()
	file, err := os.Open("test1.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cmd := strings.Split(scanner.Text(), " ")
		if len(cmd) > 0 {
			if len(cmd) == 2{
				if v, err := strconv.Atoi(cmd[0]); err == nil { epochDuration = uint64(v) }
				if v, err := strconv.Atoi(cmd[1]); err == nil { startBlock = uint64(v) }
				kyberSc.Init(startBlock, epochDuration)
			}else {
				if cmd[0] == "stake" {
					block, _ := strconv.Atoi(cmd[1])
					stake, _ := strconv.Atoi(cmd[2])
					kyberSc.Stake(uint64(block), uint64(stake), cmd[3])
				}else if cmd[0] == "withdraw" {
					block, _ := strconv.Atoi(cmd[1])
					stake, _ := strconv.Atoi(cmd[2])
					kyberSc.Withdraw(uint64(block), uint64(stake), cmd[3])
				}else if cmd[0] == "getStake" {
					epoch, _ :=  strconv.Atoi(cmd[1])
					kyberSc.GetStake(uint64(epoch), cmd[2])
				}else if cmd[0] == "vote" {
					block, _ := strconv.Atoi(cmd[1])
					voteid, _ :=  strconv.Atoi(cmd[2])
					kyberSc.Vote(uint64(block),uint64(voteid), cmd[3])
				}else if cmd[0] == "getReward" {
					epoch, _ :=  strconv.Atoi(cmd[1])
					kyberSc.GetReward(uint64(epoch), cmd[2]) //delegate getDelegatedStake  getRepresentative
				}else if cmd[0] == "delegate" {
					block, _ := strconv.Atoi(cmd[1])
					kyberSc.Delegate(uint64(block), cmd[2], cmd[3])
				}else if cmd[0] == "getDelegatedStake" {
					epoch, _ :=  strconv.Atoi(cmd[1])
					kyberSc.GetDelegatedStake(uint64(epoch), cmd[2])
				}else if cmd[0] == "getRepresentative" {
					epoch, _ :=  strconv.Atoi(cmd[1])
					kyberSc.GetRepresentative(uint64(epoch), cmd[2])
				}else if cmd[0] == "getPoolReward" {
					epoch, _ :=  strconv.Atoi(cmd[1])
					kyberSc.GetPoolReward(uint64(epoch), cmd[2])
				}
			}
		}
	}
	log.Println(time.Now().Sub(time1))
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return nil
}
