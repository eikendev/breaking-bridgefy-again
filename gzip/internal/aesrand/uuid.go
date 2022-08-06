package aesrand

// Inspired from github.com/google/uuid.

import (
	"encoding/hex"
	"io"
)

type UUID [16]byte

var (
	Nil UUID // empty UUID, all zeros
)

// NewRandomFromReader returns a UUID based on bytes read from a given io.Reader.
func (r *AesRand) NewUUID() string {
	var uuid UUID
	_, err := io.ReadFull(r, uuid[:])
	if err != nil {
		panic(err)
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10

	var buf [36]byte
	encodeHex(buf[:], uuid)
	return string(buf[:])
}

func encodeHex(dst []byte, uuid UUID) {
	hex.Encode(dst, uuid[:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], uuid[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], uuid[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], uuid[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:], uuid[10:])
}
