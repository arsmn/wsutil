package wsutil

import "errors"

var (
	ErrHandlerNotRegistered = errors.New("no handler registered for this method")
)

type Router struct {
	enc      Encoder
	handlers map[Method]Handler
}

func NewRouter(enc Encoder) *Router {
	return &Router{
		enc:      enc,
		handlers: make(map[Method]Handler),
	}
}

func (r *Router) Set(m Method, h Handler) *Router {
	r.handlers[m] = h
	return r
}

func (r *Router) Handle(data []byte) ([]byte, error) {
	req := new(Request)
	if err := r.enc.Decode(data, req); err != nil {
		return nil, err
	}

	h, ok := r.handlers[req.Method]
	if !ok {
		return nil, ErrHandlerNotRegistered
	}

	res, err := h(req)
	if err != nil {
		return nil, err
	}

	return r.enc.Encode(res)
}
