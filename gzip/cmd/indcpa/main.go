package main

import (
	"errors"
	"os"

	"github.com/eikendev/breaking-bridgefy-again/gzip/internal/commands"
	"github.com/eikendev/breaking-bridgefy-again/gzip/internal/settings"
	"github.com/jessevdk/go-flags"
)

type command struct {
	settings.Settings
	Q0 commands.Q0Command `command:"q0" description:"Play IND-CPA(0)"`
	Q1 commands.Q1Command `command:"q1" description:"Play IND-CPA(1)"`
}

var (
	cmds   command
	parser = flags.NewParser(&cmds, flags.Default)
)

func main() {
	_, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}

	if cmds.SampleHopEnd >= cmds.NetworkSize {
		panic(errors.New("can not measure message at last hop"))
	}

	if cmds.SampleSize <= 0 {
		panic(errors.New("sample size must be positive"))
	}

	settings.Runner.Run(&cmds.Settings)
}
