package aesrand

import (
	"strings"
	"testing"
)

func TestNewUUID(t *testing.T) {
	r := New(uint64(1))

	for i := 0; i < 64; i++ {
		str := r.NewUUID()

		count := strings.Count(str, "-")
		if count != 4 {
			t.Errorf("UUID must have 4 hyphens but has %d", count)
		}
	}
}
