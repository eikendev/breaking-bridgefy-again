package commands

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"math"
	"strings"

	"github.com/eikendev/breaking-bridgefy-again/gzip/internal/aesrand"
	"github.com/eikendev/breaking-bridgefy-again/gzip/internal/games"
	"github.com/eikendev/breaking-bridgefy-again/gzip/internal/mesh"
	"github.com/eikendev/breaking-bridgefy-again/gzip/internal/settings"
)

type Q0Command struct {
	Quiet       bool `short:"q" long:"quiet" description:"Whether to print game output or not"`
	PayloadSize int  `long:"payload-size" default:"4" description:"The number of bytes for one payload half"`
	Rounds      int  `long:"rounds" default:"8192" description:"How often to play the game"`
}

func (c *Q0Command) Execute(args []string) error {
	settings.Runner = c
	return nil
}

func (c *Q0Command) playGame(s *settings.Settings, r *aesrand.AesRand, block *cipher.Block) bool {
	oracle := games.NewIndCpaGame(0, c.Quiet, s, r)

	size := c.PayloadSize
	blocks := 1 + ((size - 1) / aes.BlockSize)
	payloadL := strings.Repeat("0", size)
	payloadR := r.GenerateText16Times(blocks)[:size]

	length := len(oracle.Query(payloadL, payloadR))

	network := mesh.CreateNetwork(s.NetworkSize, 0, r, "", "")
	context := mesh.BuildBroadcastContext(network, "")

	entity, err := context.MakeHop()
	if err != nil {
		panic(err)
	}

	entity.SetPayloadContent(payloadL)
	lengthL := len(games.IndCpaEncrypt(entity, block))

	entity.SetPayloadContent(payloadR)
	lengthR := len(games.IndCpaEncrypt(entity, block))

	mean := (lengthL + lengthR) / 2
	prediction := length < mean

	return oracle.Check(prediction)
}

func (c *Q0Command) Run(s *settings.Settings) {
	seed := uint32(s.Seed)
	r := aesrand.New(uint64(seed))

	block, err := aes.NewCipher([]byte("ABCDEFGHIJKLMNOP"))
	if err != nil {
		panic("Cannot initialize AES cipher")
	}

	successCount := 0

	for i := 0; i < c.Rounds; i++ {
		if c.playGame(s, r, &block) {
			successCount++
		}
	}

	advantage := 2 * math.Abs(float64(successCount)/float64(c.Rounds)-0.5)
	fmt.Printf("%d,%d,%f\n", c.PayloadSize, c.Rounds, advantage)
}
