//go:build profile
// +build profile

package commands

import (
	"testing"

	"github.com/eikendev/breaking-bridgefy-again/gzip/internal/settings"
)

func TestRun(t *testing.T) {
	cmd := SimulateCommand{}
	s := &settings.Settings{
		Seed:           42,
		SampleSize:     100,
		NetworkSize:    5,
		SampleHopStart: 0,
		SampleHopEnd:   4,
	}
	cmd.Run(s)
}
