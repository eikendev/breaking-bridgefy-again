package commands

import (
	"crypto/aes"
	"encoding/hex"
	"fmt"
	"math"

	"github.com/eikendev/breaking-bridgefy-again/gzip/internal/aesrand"
	"github.com/eikendev/breaking-bridgefy-again/gzip/internal/games"
	"github.com/eikendev/breaking-bridgefy-again/gzip/internal/settings"
)

type Q1Command struct {
	Quiet       bool `short:"q" long:"quiet" description:"Whether to print game output or not"`
	PayloadSize int  `long:"payload-size" default:"4" description:"The number of bytes for one payload half"`
	Rounds      int  `long:"rounds" default:"8192" description:"How often to play the game"`
}

func (c *Q1Command) Execute(args []string) error {
	settings.Runner = c
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func (c *Q1Command) playGame(s *settings.Settings, r *aesrand.AesRand) bool {
	oracle := games.NewIndCpaGame(1, c.Quiet, s, r)

	size := c.PayloadSize
	blocks := 1 + ((size - 1) / aes.BlockSize)
	payloadL := r.GenerateText16Times(blocks * 2)[:2*size]
	payloadR := r.GenerateText16Times(blocks * 2)[:2*size]

	c1 := oracle.Query(payloadL, payloadR)
	c2 := oracle.Query(payloadL, payloadL)

	minlen := min(len(c1), len(c2))
	match := false

	m := make(map[string]bool)

	for i := 0; i < minlen; i += aes.BlockSize {
		m[hex.EncodeToString(c1[i:i+aes.BlockSize])] = true
	}

	for i := 0; i < minlen; i += aes.BlockSize {
		if _, ok := m[hex.EncodeToString(c2[i:i+aes.BlockSize])]; ok {
			match = true
			break
		}
	}

	return oracle.Check(match)
}

func (c *Q1Command) Run(s *settings.Settings) {
	seed := uint32(s.Seed)
	r := aesrand.New(uint64(seed))

	successCount := 0

	for i := 0; i < c.Rounds; i++ {
		if c.playGame(s, r) {
			successCount++
		}
	}

	advantage := 2 * math.Abs(float64(successCount)/float64(c.Rounds)-0.5)
	fmt.Printf("%d,%d,%f\n", c.PayloadSize, c.Rounds, advantage)
}
