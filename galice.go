package galice

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Logger - signature for logging function
type Logger func(error)

// Client - Alice API client
type Client struct {
	autoPings            bool
	autoDanderousContext bool
	logger               Logger
}

// New creates new instance of Alice API client
func New(autoPings bool, autoDanderousContext bool) *Client {
	return &Client{
		autoPings,
		autoDanderousContext,
		func(val error) {
			fmt.Printf("An error accoured while handling Alice request: %v", val)
		},
	}
}

// AliceHandler - signature of Alice request handler
type AliceHandler func(InputData) OutputData

// CreateHandler - creates new handler function for HTTP server based on provided AliceHander
func (c *Client) CreateHandler(fn AliceHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if val := recover(); val != nil {
				c.logger(fmt.Errorf("Unexpected error: %v", val))
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		decoder := json.NewDecoder(r.Body)
		var ai InputData
		err := decoder.Decode(&ai)
		if err != nil {
			c.logger(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if c.autoPings && ai.Request.IsPing() {
			p := Pong(ai)
			b, err := json.Marshal(p)
			if err != nil {
				c.logger(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(b)
			return
		}
		if c.autoDanderousContext && ai.Request.IsDangerousContext() {
			d := Dangerous(ai)
			b, err := json.Marshal(d)
			if err != nil {
				c.logger(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(b)
			return
		}
		ao := fn(ai)
		b, err := json.Marshal(ao)
		if err != nil {
			c.logger(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(b)
	})
}
