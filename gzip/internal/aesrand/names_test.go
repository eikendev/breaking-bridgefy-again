package aesrand

import "testing"

func TestChoose128(t *testing.T) {
	r := New(uint64(1))

	for i := 0; i < 128*8; i++ {
		r := r.choose128()
		if r >= 128 {
			t.Errorf("random value greater than allowed (%d)", r)
		}
	}
}

func TestGetRandomName(t *testing.T) {
	r := New(uint64(1))

	for i := 0; i < 64; i++ {
		name := r.GetRandomName()
		found := false

		for i := range names {
			if names[i] == name {
				found = true
			}
		}

		if !found {
			t.Errorf("returned name that is not in the list")
		}
	}
}
