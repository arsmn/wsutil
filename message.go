package wsutil

type (
	Method  uint
	Payload map[string]interface{}
	Handler func(*Request) (*Response, error)
)

type Request struct {
	ID     uint64      `msgpack:"id,omitempty" json:"id"`
	Method Method      `msgpack:"method,omitempty" json:"method,omitempty"`
	Status int         `msgpack:"status,omitempty" json:"status,omitempty"`
	Data   interface{} `msgpack:"data,omitempty" json:"data,omitempty"`
}

type Response struct {
	ID     uint64      `msgpack:"id,omitempty" json:"id"`
	Method Method      `msgpack:"method,omitempty" json:"method,omitempty"`
	Status int         `msgpack:"status,omitempty" json:"status,omitempty"`
	Data   interface{} `msgpack:"data,omitempty" json:"data,omitempty"`
}
