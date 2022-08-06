package games

import (
	"crypto/aes"
	"crypto/cipher"
	"log"

	"github.com/eikendev/breaking-bridgefy-again/gzip/internal/aesrand"
	"github.com/eikendev/breaking-bridgefy-again/gzip/internal/mesh"
	"github.com/eikendev/breaking-bridgefy-again/gzip/internal/settings"
)

type IndCpaGame struct {
	quiet   bool
	queries uint
	block   cipher.Block
	b       bool
	s       *settings.Settings
	r       *aesrand.AesRand
	network *mesh.LinearNetwork
}

func NewIndCpaGame(queries uint, quiet bool, s *settings.Settings, r *aesrand.AesRand) *IndCpaGame {
	randBytes := make([]byte, 16)
	r.Read(randBytes)
	left := randBytes[0] & 1

	if !quiet {
		log.Printf("C: Initialized b=%d", left)
	}

	block, err := aes.NewCipher([]byte("abcdefghijklmnop"))
	if err != nil {
		panic("Cannot initialize AES cipher")
	}

	network := mesh.CreateNetwork(s.NetworkSize, s.SampleHopEnd, r, "", "")

	return &IndCpaGame{
		quiet:   quiet,
		queries: queries,
		block:   block,
		b:       left > 0,
		s:       s,
		r:       r,
		network: network,
	}
}

func IndCpaEncrypt(entity *mesh.BleEntity, block *cipher.Block) []byte {
	compressed := entity.Compress()

	length := len(compressed)
	over := length % aes.BlockSize
	if over > 0 {
		length += aes.BlockSize - over
	}

	c := make([]byte, length)
	copied := copy(c, compressed)

	if copied != len(compressed) {
		panic("Failed to encrypt all packet data")
	}

	for i := 0; i < length; i += aes.BlockSize {
		(*block).Encrypt(c[i:i+aes.BlockSize], c[i:i+aes.BlockSize])
	}

	return c
}

func (g *IndCpaGame) Query(left, right string) []byte {
	if !g.quiet {
		log.Printf("A -> C: Payload for b=1 (left)  is %s\n", left)
		log.Printf("A -> C: Payload for b=0 (right) is %s\n", right)
	}

	if len(left) != len(right) {
		log.Fatalf("Must supply equal lengths to the query")
	}

	context := mesh.BuildBroadcastContext(g.network, "")

	entity, err := context.MakeHop()
	if err != nil {
		panic(err)
	}

	var payload string
	if g.b {
		payload = left
	} else {
		payload = right
	}

	entity.SetPayloadContent(payload)
	c := IndCpaEncrypt(entity, &g.block)

	if !g.quiet {
		log.Printf("C -> A: Length is %d\n", len(c))
	}

	return c
}

func b2i(g bool) int {
	if g {
		return 1
	} else {
		return 0
	}
}

func (g *IndCpaGame) Check(prediction bool) bool {
	if !g.quiet && prediction != g.b {
		log.Printf("A: Predicting b=%d when real b=%d\n", b2i(prediction), b2i(g.b))
	}
	return prediction == g.b
}
