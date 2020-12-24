package wsutil

import (
	"encoding/json"

	"github.com/vmihailenco/msgpack/v5"
)

type Encoder interface {
	Encode(v interface{}) ([]byte, error)
	Decode(data []byte, v interface{}) error
}

type JsonEncoder struct{}

func (je JsonEncoder) Encode(v interface{}) ([]byte, error) {
	if v == nil {
		return nil, nil
	}
	return json.Marshal(v)
}

func (je JsonEncoder) Decode(data []byte, v interface{}) error {
	if _, ok := v.(*interface{}); ok {
		return nil
	}
	return json.Unmarshal(data, v)
}

type MsgPackEncoder struct{}

func (mpe MsgPackEncoder) Encode(v interface{}) ([]byte, error) {
	if v == nil {
		return nil, nil
	}
	return msgpack.Marshal(v)
}

func (mpe MsgPackEncoder) Decode(data []byte, v interface{}) error {
	if _, ok := v.(*interface{}); ok {
		return nil
	}
	return msgpack.Unmarshal(data, v)
}
