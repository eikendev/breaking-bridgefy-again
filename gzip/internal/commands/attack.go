package commands

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/eikendev/breaking-bridgefy-again/gzip/internal/fs"
	"github.com/eikendev/breaking-bridgefy-again/gzip/internal/measure"
	"github.com/eikendev/breaking-bridgefy-again/gzip/internal/payloads"
	"github.com/eikendev/breaking-bridgefy-again/gzip/internal/settings"
)

const (
	SeedOffset = 1 << 63
	SeedGap    = 1 << 18
)

type AttackCommand struct {
	Runs      int `long:"runs" description:"How many samples to gather for each payload content"`
	Arguments struct {
		PCS string `required:"true" positional-arg-name:"payload content set" description:"The payload content set to use"`
	} `positional-args:"true"`
}

func (c *AttackCommand) Execute(args []string) error {
	settings.Runner = c
	return nil
}

func (c *AttackCommand) Run(s *settings.Settings) {
	outdir := fmt.Sprintf("attacks-%s-M%d-H%d-n%d", c.Arguments.PCS, s.SampleSize, s.SampleHopEnd, c.Runs)
	fs.AssertNotExists(outdir)

	err := os.MkdirAll(outdir, 0755)
	if err != nil {
		log.Fatalf("Cannot create directory: %s", err)
	}

	pcs := payloads.Get(c.Arguments.PCS)
	localSeed := s.Seed + SeedOffset + pcs.SeedOffset
	mctx := &measure.MeasurementCtx{
		NewSample: measure.NewMemorySample,
	}

	for _, payload := range pcs.PayloadContents {
		outfile := path.Join(outdir, payload+".json")

		for i := 0; i < c.Runs; i++ {
			localSeed += SeedGap

			sample := measure.RunMeasurement(s, localSeed, nil, []string{payload}, mctx)

			fs.AppendString(outfile, sample.Format()+"\n")
		}
	}
}
