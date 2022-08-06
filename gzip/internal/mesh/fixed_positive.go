package mesh

import "github.com/vmihailenco/msgpack/v5"

type FixedNegative struct {
	Value int8
}

func FN(value int) FixedNegative {
	return FixedNegative{
		Value: int8(value),
	}
}

var _ msgpack.Marshaler = (*FixedNegative)(nil)

func (i *FixedNegative) MarshalMsgpack() ([]byte, error) {
	data, err := msgpack.Marshal(i.Value)
	return data[1:], err
}

var _ msgpack.Unmarshaler = (*FixedNegative)(nil)

func (i *FixedNegative) UnmarshalMsgpack(b []byte) error {
	return msgpack.Unmarshal(b, &i.Value)
}
