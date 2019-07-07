package galice

import (
	"encoding/json"
	"net/http"
)

// Client - Alice API client
type Client struct {
	autoPings            bool
	autoDanderousContext bool
}

// New creates new instance of Alice API client
func New(autoPings bool, autoDanderousContext bool) *Client {
	return &Client{
		autoPings,
		autoDanderousContext,
	}
}

// AliceHandler - signature of Alice request handler
type AliceHandler func(InputData) OutputData

// CreateHandler - creates new handler function for HTTP server based on provided AliceHander
func (c *Client) CreateHandler(fn AliceHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				panic(err) // TODO Handle it properly
			}
		}()
		decoder := json.NewDecoder(r.Body)
		var ai InputData
		err := decoder.Decode(&ai)
		if err != nil {
			panic(err) // TODO Handle it properly
		}
		if c.autoPings && ai.Request.IsPing() {
			p := Pong(ai)
			b, err := json.Marshal(p)
			if err != nil {
				panic(err) // TODO Handle it properly
			}
			w.Write(b)
			return
		}
		if c.autoDanderousContext && ai.Request.IsDangerousContext() {
			d := Dangerous(ai)
			b, err := json.Marshal(d)
			if err != nil {
				panic(err) // TODO Handle it properly
			}
			w.Write(b)
			return
		}
		ao := fn(ai)
		b, err := json.Marshal(ao)
		if err != nil {
			panic(err) // TODO Handle it properly
		}
		w.Write(b)
	})
}
