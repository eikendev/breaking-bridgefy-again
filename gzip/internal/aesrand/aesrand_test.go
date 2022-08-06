package aesrand

import "testing"

func TestGenerateShortDelay(t *testing.T) {
	r := New(uint64(1))

	for i := 0; i < 1024; i++ {
		delay := r.GenerateShortDelay()

		if delay >= 64 || delay < 4 {
			t.Errorf("short delay %d is out of bounds", delay)
		}
	}
}

func TestGenerateDelay(t *testing.T) {
	r := New(uint64(1))

	for i := 0; i < 4096; i++ {
		delay := r.GenerateDelay()

		if delay >= 512 || delay < 128 {
			t.Errorf("delay %d is out of bounds", delay)
		}
	}
}

func TestGenerateLongDelay(t *testing.T) {
	r := New(uint64(1))

	for i := 0; i < 16384; i++ {
		delay := r.GenerateLongDelay()

		if delay >= 16777216 {
			t.Errorf("long delay %d is out of bounds", delay)
		}
	}
}
