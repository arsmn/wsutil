package wsutil

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/fasthttp/websocket"
)

type Client struct {
	mutex   sync.Mutex
	conn    *websocket.Conn
	pending map[uint64]*Call
	counter uint64
	enc     Encoder
}

type Call struct {
	Req  Request
	Res  Response
	Done chan bool
	Err  error
}

func NewCall(req Request) *Call {
	done := make(chan bool)
	return &Call{
		Req:  req,
		Done: done,
	}
}

func NewClient(enc Encoder) *Client {
	return &Client{
		pending: make(map[uint64]*Call, 1),
		counter: 1,
		enc:     enc,
	}
}

func (c *Client) read() {
	var (
		data []byte
		err  error
	)

	for err == nil {
		if _, data, err = c.conn.ReadMessage(); err != nil {
			err = fmt.Errorf("error reading message: %q", err)
			continue
		}

		var res Response
		if err = c.enc.Decode(data, &res); err != nil {
			err = fmt.Errorf("error decoding message: %q", err)
			continue
		}

		c.mutex.Lock()
		call := c.pending[res.ID]
		delete(c.pending, res.ID)
		c.mutex.Unlock()
		if call == nil {
			err = errors.New("no pending request found")
			continue
		}
		call.Res = res
		call.Done <- true
	}
	c.mutex.Lock()
	for _, call := range c.pending {
		call.Err = err
		call.Done <- true
	}
	c.mutex.Unlock()
}

func (c *Client) Connect(url string, h http.Header) error {
	conn, _, err := websocket.DefaultDialer.Dial(url, h)
	if err != nil {
		return err
	}
	c.conn = conn
	go c.read()
	return nil
}

func (c *Client) Request(req Request) (Response, error) {
	c.mutex.Lock()
	id := c.counter
	c.counter++
	req.ID = id
	call := NewCall(req)
	data, err := c.enc.Encode(req)
	if err != nil {
		c.mutex.Unlock()
		return Response{}, err
	}

	c.pending[id] = call
	if err := c.conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		delete(c.pending, id)
		c.mutex.Unlock()
		return Response{}, err
	}
	c.mutex.Unlock()
	select {
	case <-call.Done:
	case <-time.After(2 * time.Second):
		call.Err = errors.New("request timeout")
	}
	if call.Err != nil {
		return Response{}, call.Err
	}
	return call.Res, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
