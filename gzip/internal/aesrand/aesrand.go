package aesrand

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/binary"
)

type AesRand struct {
	counter uint64
	block   cipher.Block
}

func New(seed uint64) *AesRand {
	key := make([]byte, 16)
	binary.LittleEndian.PutUint64(key, seed)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic("cannot initialize AES cipher")
	}

	aesRand := AesRand{
		counter: 0,
		block:   block,
	}

	return &aesRand
}

func (r *AesRand) Read(p []byte) (n int, err error) {
	n = r.block.BlockSize()
	if len(p) != n {
		panic("buffer is too small")
	}

	binary.LittleEndian.PutUint64(p, r.counter)
	r.block.Encrypt(p, p)

	r.counter++
	if r.counter == 0 {
		panic("counter is exhausted")
	}

	return
}

func (r *AesRand) GenerateText16() string {
	buffer := make([]byte, 16)
	r.Read(buffer)
	return base64.StdEncoding.EncodeToString(buffer)[:16]
}

func (r *AesRand) GenerateText16Times(n int) string {
	out := ""
	for i := 0; i < n; i++ {
		out += r.GenerateText16()
	}
	return out
}

func (r *AesRand) GenerateText8() string {
	return r.GenerateText16()[:8]
}

func (r *AesRand) GenerateLongDelay() uint32 {
	buffer := make([]byte, 16)
	r.Read(buffer)
	buffer[0] = 0
	return binary.BigEndian.Uint32(buffer[:4])
}

func (r *AesRand) GenerateShortDelay() uint8 {
	buffer := make([]byte, 16)
	r.Read(buffer)
	buffer[0] |= 0b00000100
	buffer[0] &= 0b00111111
	return buffer[0]
}

func (r *AesRand) GenerateDelay() uint16 {
	buffer := make([]byte, 16)
	r.Read(buffer)
	buffer[0] &= 0b00000001
	buffer[1] |= 0b10000000
	return binary.BigEndian.Uint16(buffer[:2])
}
