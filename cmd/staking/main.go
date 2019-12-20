package main

import (
	"fmt"
	"github.com/kyber/staking/api"
	"log"
	"os"
	"runtime"
	"sort"

	"github.com/kyber/staking/contestant"
	cli "gopkg.in/urfave/cli.v1"
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

	startBlock := uint64(2000)
	epochDuration := uint64(1000)
	kyberSc := contestant.NewKyberStakingContract(startBlock, epochDuration)

	kyberSc.Stake(uint64(900), uint64(500), "mike")
	kyberSc.GetStake(uint64(0), "mike")
	kyberSc.GetStake(uint64(1), "mike")
	log.Println("staking done")

	server := api.NewServer(kyberSc)
	server.Run()

	return nil
}
