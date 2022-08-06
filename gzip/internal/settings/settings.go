package settings

type Settings struct {
	Test           string `short:"t" default:""`
	Seed           uint64 `long:"seed" default:"42" description:"The seed for the PRNG"`
	SampleSize     int    `short:"n" long:"sample-size" default:"10" description:"Generate 2^n packets per payload per hop"`
	NetworkSize    int    `long:"network-size" default:"50" description:"The size a network should be"`
	SampleHopStart int    `long:"hop-start" default:"0" description:"Which hop to start capture the broadcast at"`
	SampleHopEnd   int    `long:"hop-end" default:"10" description:"Which hop to stop capture the broadcast at"`
	Debug          bool   `long:"debug" description:"Print debug output"`
}

// Runnable defines an interface for subcommands that take the global settings and a password.
type Runnable interface {
	Run(*Settings)
}

// Runner is the subcommand to run after all arguments were parsed.
var Runner Runnable
