package mesh

import (
	"bytes"
	"compress/gzip"
	"log"

	"github.com/vmihailenco/msgpack/v5"
)

// Message is the object used to bundle relevant information for a Bridgefy message.
type Message struct {
	Ct string  `msgpack:"ct"`
	Mt float64 `msgpack:"mt"`
	Ku float64 `msgpack:"ku"`
	Mi string  `msgpack:"mi"`
	Nm string  `msgpack:"nm"`
	Ds float64 `msgpack:"ds"`
	Et float64 `msgpack:"et"`
}

// ForwardPacket is a container for a message to be transmitted over the mesh network.
type ForwardPacket struct {
	Added         uint64        `msgpack:"added"`
	ID            string        `msgpack:"id"`
	Payload       Message       `msgpack:"payload"`
	Enc_payload   FixedNegative `msgpack:"enc_payload"`
	Sender        string        `msgpack:"sender"`
	Creation      uint64        `msgpack:"creation"`
	Expiration    uint16        `msgpack:"expiration"`
	Receiver_type FixedPositive `msgpack:"receiver_type"`
	Hops          FixedPositive `msgpack:"hops"`
	Profile       FixedPositive `msgpack:"profile"`
	Track         []uint32      `msgpack:"track"`
}

// ForwardTransaction is a collection of ForwardPacket instances.
type ForwardTransaction struct {
	Dump   bool          `msgpack:"dump"`
	Sender string        `msgpack:"sender"`
	Mesh   ForwardPacket `msgpack:"mesh"`
}

// BleEntity is an object that is transmitted over the mesh.
type BleEntity struct {
	ID string             `msgpack:"id"`
	Et FixedPositive      `msgpack:"et"`
	Ct ForwardTransaction `msgpack:"ct"`
}

// Set the original sender of the message.
func (b *BleEntity) SetMessageSender(sender string) {
	b.Ct.Mesh.Sender = sender
}

// Set the payload content of the message.
func (b *BleEntity) SetPayloadContent(content string) {
	b.Ct.Mesh.Payload.Ct = content
}

// Serialize and compress the entity.
func (b *BleEntity) Compress() []byte {
	serialized, err := msgpack.Marshal(b)
	if err != nil {
		panic(err)
	}

	var compressed bytes.Buffer
	g := gzip.NewWriter(&compressed)

	if _, err = g.Write(serialized); err != nil {
		log.Fatal(err)
		return nil
	}
	if err = g.Close(); err != nil {
		log.Fatal(err)
		return nil
	}

	return compressed.Bytes()
}
