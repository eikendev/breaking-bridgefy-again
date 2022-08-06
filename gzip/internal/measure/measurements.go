package measure

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"runtime"
	"sync"

	"github.com/eikendev/breaking-bridgefy-again/gzip/internal/aesrand"
	"github.com/eikendev/breaking-bridgefy-again/gzip/internal/mesh"
	"github.com/eikendev/breaking-bridgefy-again/gzip/internal/settings"
)

const (
	ReportEvery = 128
)

type NewSampleFunc func() Sample

type MeasurementCtx struct {
	NewSample NewSampleFunc
}

func measureLength(entity *mesh.BleEntity, content string) int {
	entity.SetPayloadContent(content)

	compressed := entity.Compress()

	if false {
		fmt.Printf("%+v\n", entity)
		fmt.Println(base64.StdEncoding.EncodeToString(compressed))
		fmt.Println(hex.EncodeToString(compressed))
	}

	return len(compressed)
}

func RunPartialMeasurement(seed uint64, ctx *RunContext, count int, s *settings.Settings, payloads []string, sample Sample) {
	r := aesrand.New(seed)

	for i := 0; i < count; i++ {
		if len(payloads) > 1 {
			log.Printf("Processing sample %d/%d\n", i+1, count)
		}

		network := mesh.CreateNetwork(s.NetworkSize, s.SampleHopEnd, r, ctx.SenderName, ctx.SenderUUID)
		context := mesh.BuildBroadcastContext(network, "")

		memory := make([]uint64, len(payloads))

		for hop := 0; hop <= s.SampleHopEnd; hop++ {
			entity, err := context.MakeHop()
			if err != nil {
				panic(err)
			}

			if hop < s.SampleHopStart {
				continue
			}

			for j, payload := range payloads {
				length := uint64(measureLength(entity, payload))

				result := &LengthResult{
					Payload:        payload,
					Hop:            uint64(context.Hop),
					LastHop:        int64(context.Hop - 1),
					Length:         length,
					PreviousLength: memory[j],
				}
				sample.Process(result, 1)

				memory[j] = length
			}
		}
	}
}

func RunMeasurement(s *settings.Settings, seed uint64, ctx *RunContext, payloads []string, mctx *MeasurementCtx) Sample {
	sampleSize := int(math.Pow(2, float64(s.SampleSize)))

	bucketSize := bucketSize(sampleSize, runtime.GOMAXPROCS(0)*8)
	log.Printf("Assigning %d rounds to each thread\n", bucketSize)

	if ctx == nil {
		ctx = DefaultRunContext()
	}

	wgSenders := sync.WaitGroup{}
	collector := make(chan Sample)

	for threadID := 0; threadID*bucketSize < sampleSize; threadID++ {
		wgSenders.Add(1)

		count := min(bucketSize, sampleSize-threadID*bucketSize)

		if s.Debug {
			log.Printf("Running thread %d with count %d and seed %d", threadID, count, seed+uint64(threadID))
		}

		go func(localSeed uint64, localCount int, c chan Sample) {
			defer wgSenders.Done()

			sample := mctx.NewSample()

			RunPartialMeasurement(localSeed, ctx, localCount, s, payloads, sample)

			c <- sample
		}(seed+uint64(threadID), count, collector)
	}

	wgReader := sync.WaitGroup{}
	wgReader.Add(1)
	resultSample := mctx.NewSample()

	go func(c <-chan Sample) {
		defer wgReader.Done()

		for sample := range c {
			resultSample.Merge(sample)
		}
	}(collector)

	wgSenders.Wait()
	close(collector)
	wgReader.Wait()

	return resultSample
}
