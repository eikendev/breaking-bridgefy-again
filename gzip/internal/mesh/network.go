package mesh

import (
	"github.com/eikendev/breaking-bridgefy-again/gzip/internal/aesrand"
)

type LinearNetwork struct {
	Size       int
	MaxHop     int
	NodeUUIDs  []string
	SenderName string
	Rand       *aesrand.AesRand
}

// Create a mesh network with `size` nodes.
func CreateNetwork(size, maxHop int, r *aesrand.AesRand, senderName, senderUUID string) *LinearNetwork {
	if senderName == "" {
		senderName = r.GetRandomName()
	}

	if senderUUID == "" {
		senderUUID = r.NewUUID()
	}

	uuids := make([]string, maxHop+2)
	uuids[0] = senderUUID

	for i := 1; i <= maxHop+1; i++ {
		uuids[i] = r.NewUUID()
	}

	return &LinearNetwork{
		Size:       size,
		MaxHop:     maxHop,
		NodeUUIDs:  uuids,
		SenderName: senderName,
		Rand:       r,
	}
}

func (n *LinearNetwork) getNodeUUID(index int) string {
	return n.NodeUUIDs[index]
}
