package measure

import (
	"encoding/json"
)

type LengthResult struct {
	Payload        string
	Hop            uint64
	LastHop        int64
	Length         uint64
	PreviousLength uint64
}

type Sample interface {
	Process(*LengthResult, uint64)
	Format() string
	Merge(Sample)
}

type LengthMemory map[string]map[uint64]map[uint64]map[uint64]uint64

func (s LengthMemory) Process(payload string, hop, length, previousLength, count uint64) {
	if _, exists := s[payload]; !exists {
		s[payload] = make(map[uint64]map[uint64]map[uint64]uint64)
	}
	if _, exists := s[payload][hop]; !exists {
		s[payload][hop] = make(map[uint64]map[uint64]uint64)
	}
	if _, exists := s[payload][hop][length]; !exists {
		s[payload][hop][length] = make(map[uint64]uint64)
	}
	if _, exists := s[payload][hop][length][previousLength]; !exists {
		s[payload][hop][length][previousLength] = 0
	}
	s[payload][hop][length][previousLength] += count
}

type MemorySample struct {
	Memory LengthMemory
}

func NewMemorySample() Sample {
	memory := make(LengthMemory)

	return &MemorySample{
		Memory: memory,
	}
}

func (s *MemorySample) Process(r *LengthResult, count uint64) {
	s.Memory.Process(r.Payload, r.Hop, r.Length, r.PreviousLength, count)
}

func (s *MemorySample) Format() string {
	encoded, _ := json.Marshal(s.Memory)
	return string(encoded)
}

func (s *MemorySample) Merge(other Sample) {
	o := other.(*MemorySample)

	for payload, hops := range o.Memory {
		for hop, lengths := range hops {
			for length, previousLengths := range lengths {
				for previousLength, count := range previousLengths {
					s.Memory.Process(payload, hop, length, previousLength, count)
				}
			}
		}
	}
}
