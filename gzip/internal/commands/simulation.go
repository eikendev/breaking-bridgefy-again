package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/eikendev/breaking-bridgefy-again/gzip/internal/fs"
	"github.com/eikendev/breaking-bridgefy-again/gzip/internal/measure"
	"github.com/eikendev/breaking-bridgefy-again/gzip/internal/payloads"
	"github.com/eikendev/breaking-bridgefy-again/gzip/internal/settings"
)

type SimulateCommand struct {
	Arguments struct {
		PCS string `required:"true" positional-arg-name:"payload content set" description:"The payload content set to use"`
	} `positional-args:"true"`
}

func (c *SimulateCommand) Execute(args []string) error {
	settings.Runner = c
	return nil
}

func (c *SimulateCommand) Run(s *settings.Settings) {
	outfile := fmt.Sprintf("simulation-%s-N%d-H%d.json", c.Arguments.PCS, s.SampleSize, s.SampleHopEnd)
	fs.AssertNotExists(outfile)

	f, err := os.Create(outfile)
	if err != nil {
		log.Fatalf("Cannot create file: %s", err)
	}
	defer f.Close()

	pcs := payloads.Get(c.Arguments.PCS)
	mctx := &measure.MeasurementCtx{
		NewSample: measure.NewMemorySample,
	}

	sample := measure.RunMeasurement(s, s.Seed+pcs.SeedOffset, nil, pcs.PayloadContents, mctx)

	f.WriteString(sample.Format())
}
