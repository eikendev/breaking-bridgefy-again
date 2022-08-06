package mesh

import (
	"fmt"
	"hash/crc32"
)

func makeTimestamp() uint64 {
	return 1619088192134
}

type BroadcastContext struct {
	Network *LinearNetwork
	Hop     int
	Entity  BleEntity
}

// Assembles the context for a broadcast message.
func BuildBroadcastContext(network *LinearNetwork, content string) *BroadcastContext {
	startTime := makeTimestamp() + uint64(network.Rand.GenerateLongDelay())
	creationOffset := uint64(network.Rand.GenerateShortDelay())

	message := Message{
		Ct: content,
		Ku: 1,
		Mi: network.Rand.NewUUID(),
		Et: 2,
		Mt: 0,
		Nm: network.SenderName,
		Ds: float64(startTime),
	}

	forwardPacket := ForwardPacket{
		Added:         startTime + creationOffset,
		ID:            network.Rand.NewUUID(),
		Payload:       message,
		Enc_payload:   FN(-1),
		Sender:        network.getNodeUUID(0),
		Creation:      startTime + creationOffset,
		Expiration:    3600,
		Receiver_type: FP(1),
		Hops:          FP(network.Size + 1),
		Profile:       FP(1),
		Track:         []uint32{crc32.ChecksumIEEE([]byte(network.getNodeUUID(0)))},
	}

	forwardTransaction := ForwardTransaction{
		Dump:   false,
		Sender: network.getNodeUUID(0),
		Mesh:   forwardPacket,
	}

	bleEntity := BleEntity{
		ID: network.Rand.NewUUID(),
		Et: FP(3),
		Ct: forwardTransaction,
	}

	return &BroadcastContext{
		Network: network,
		Hop:     -1,
		Entity:  bleEntity,
	}
}

func (b *BroadcastContext) MakeHop() (*BleEntity, error) {
	if b.Hop > b.Network.MaxHop {
		return nil, fmt.Errorf("node %d does not send any message in this mesh network", b.Hop)
	}

	b.Hop++

	b.Entity.ID = b.Network.Rand.NewUUID()

	b.Entity.Ct.Sender = b.Network.getNodeUUID(b.Hop)

	b.Entity.Ct.Mesh.Added = b.Entity.Ct.Mesh.Added + uint64(b.Network.Rand.GenerateDelay())
	b.Entity.Ct.Mesh.Hops.Value--

	newChecksum := crc32.ChecksumIEEE([]byte(b.Network.getNodeUUID(b.Hop + 1)))
	b.Entity.Ct.Mesh.Track = append(b.Entity.Ct.Mesh.Track, newChecksum)

	return &b.Entity, nil
}

func (b *BroadcastContext) SetMessageSender(sender string) {
	b.Entity.SetMessageSender(sender)
}

func (b *BroadcastContext) SetPayloadContent(content string) {
	b.Entity.SetPayloadContent(content)
}
