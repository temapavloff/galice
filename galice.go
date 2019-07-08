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

// SetLogger sets logger function to current client
func (c *Client) SetLogger(logger Logger) {
	c.logger = logger
}

// New creates new instance of Alice API client
func New(autoPings bool, autoDanderousContext bool) *Client {
	return &Client{
		autoPings,
		autoDanderousContext,
		func(val error) {
			fmt.Printf("An error accured while handling Alice request: %v", val)
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

		if r.Body == nil {
			c.logger(fmt.Errorf("Empty request body"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		decoder := json.NewDecoder(r.Body)
		var ai InputData
		err := decoder.Decode(&ai)
		if err != nil {
			c.logger(fmt.Errorf("Error while decoding Alice request: %v", err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if c.autoPings && ai.Request.IsPing() {
			p := Pong(ai)
			b, err := json.Marshal(p)
			if err != nil {
				c.logger(fmt.Errorf("Error while marshaling ping response: %v", err))
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
