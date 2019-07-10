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

// AliceHandlerError - error representation which may occure while handling Alice request
type AliceHandlerError struct {
	Message      string
	ResponseCode int
}

func (a *AliceHandlerError) Error() string {
	return a.Message
}

// AliceHandler - signature of Alice request handler
type AliceHandler func(InputData) (OutputData, error)

// CreateHandler - creates new handler function for HTTP server based on provided AliceHander
func (c *Client) CreateHandler(fn AliceHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if val := recover(); val != nil {
				c.logger(fmt.Errorf("Unexpected error: %v", val))
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		err := c.handleRequest(w, r, fn)
		if err != nil {
			c.logger(err)
			w.WriteHeader(err.ResponseCode)
		}
	})
}

func (c *Client) handleRequest(w http.ResponseWriter, r *http.Request, fn AliceHandler) *AliceHandlerError {
	if r.Body == nil {
		return &AliceHandlerError{"Empty request body", http.StatusBadRequest}
	}
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	var i InputData
	err := decoder.Decode(&i)
	if err != nil {
		return &AliceHandlerError{fmt.Sprintf("Error while decoding Alice request: %v", err), http.StatusBadRequest}
	}

	var o OutputData
	switch {
	case c.autoPings && i.Request.IsPing():
		o = Pong(i)
	case c.autoDanderousContext && i.Request.IsDangerousContext():
		o = Dangerous(i)
	default:
		o, err = fn(i)
		if err != nil {
			c.logger(err)
		}
	}

	b, err := json.Marshal(o)
	if err != nil {
		return &AliceHandlerError{fmt.Sprintf("Error marshaling response: %v", err), http.StatusInternalServerError}
	}

	w.Write(b)

	return nil
}
