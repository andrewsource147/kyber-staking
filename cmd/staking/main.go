package main

import (
	"fmt"
	"github.com/kyber/staking/api"
	"log"
	"os"
	"runtime"
	"sort"

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
	server := api.NewServer()
	server.Run()

	return nil
}
