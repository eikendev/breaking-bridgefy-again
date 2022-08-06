package payloads

import "log"

type PayloadContentSet struct {
	PayloadContents []string
	SeedOffset      uint64
}

var pcSets = map[string]*PayloadContentSet{
	"bytes":        getPCSBytes(),
	"rockyou8n256": getPCSRockyou8n256(),
	"rockyou6n256": getPCSRockyou6n256(),
	"rockyouXn128": getPCSRockyouXn128(),
	"rockyouXn256": getPCSRockyouXn256(),
}

func Get(name string) *PayloadContentSet {
	if pcs, ok := pcSets[name]; ok {
		return pcs
	}

	log.Fatal("Please supply a valid payload content set")
	return nil
}

func getPCSBytes() *PayloadContentSet {
	return &PayloadContentSet{
		PayloadContents: pcsBytes[:],
		SeedOffset:      0 << 60,
	}
}

func getPCSRockyou8n256() *PayloadContentSet {
	return &PayloadContentSet{
		PayloadContents: pcsRockyou8n256[:],
		SeedOffset:      3 << 60,
	}
}

func getPCSRockyou6n256() *PayloadContentSet {
	return &PayloadContentSet{
		PayloadContents: pcsRockyou6n256[:],
		SeedOffset:      6 << 60,
	}
}

func getPCSRockyouXn128() *PayloadContentSet {
	return &PayloadContentSet{
		PayloadContents: pcsRockyouXn128[:],
		SeedOffset:      1 << 60,
	}
}

func getPCSRockyouXn256() *PayloadContentSet {
	return &PayloadContentSet{
		PayloadContents: pcsRockyouXn256[:],
		SeedOffset:      7 << 60,
	}
}
