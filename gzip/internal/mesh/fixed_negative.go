package mesh

import "github.com/vmihailenco/msgpack/v5"

type FixedPositive struct {
	Value int8
}

func FP(value int) FixedPositive {
	return FixedPositive{
		Value: int8(value),
	}
}

var _ msgpack.Marshaler = (*FixedPositive)(nil)

func (i *FixedPositive) MarshalMsgpack() ([]byte, error) {
	data, err := msgpack.Marshal(i.Value)
	return data[1:], err
}

var _ msgpack.Unmarshaler = (*FixedPositive)(nil)

func (i *FixedPositive) UnmarshalMsgpack(b []byte) error {
	return msgpack.Unmarshal(b, &i.Value)
}
