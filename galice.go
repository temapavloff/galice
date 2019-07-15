package galice

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Logger is a signature for logging function used by Client
type Logger func(error)

// Client represents Alice API client, allows to create HTTP handler function for Alice API incoming webhooks
type Client struct {
	autoPings            bool   // should Alice API healthcheks be handled automatically
	autoDanderousContext bool   // should dangerous context be handled automatically
	logger               Logger // logging function
}

// default logger for Client
var defaultLogger = log.New(os.Stderr, "", 0)

// SetLogger sets logger function to current client.
// Logger used internally for logging any errors occured inside Alice request hanler:
// bad requests, invalid responses, unexpected panics, etc.
// If not called the default logger will be used.
// The default logger writes int stderr.
func (c *Client) SetLogger(logger Logger) {
	c.logger = logger
}

// New creates new Alice API client. The autoPings flag tells client to automatically
// respond to Alice API healthchecks. The autoDanderousContext tells client to
// automatically handle requests marked as dangerous (suicide, hate speech, threats)
// by Alice API.
func New(autoPings bool, autoDanderousContext bool) *Client {
	return &Client{
		autoPings,
		autoDanderousContext,
		func(val error) {
			defaultLogger.Println(val)
		},
	}
}

// AliceHandlerError represents error which may occure while handling Alice request
type AliceHandlerError struct {
	Message      string // Error message
	ResponseCode int    // HTTP status code
}

// Error implements error interface
func (a *AliceHandlerError) Error() string {
	return a.Message
}

// AliceHandler is a signature of Alice request handler. It represents function
// which accepts InputData (go struct, contains Alice API incoming data) and must
// return OutputData (go struct, contains Alice API outcoming data) and optional error
//Notice that error is used only for additional logging, so function mus return correct
// OutputData even if something went wrong
type AliceHandler func(InputData) (OutputData, error)

// CreateHandler creates new http.Handler for Alice API incoming webhooks based on
// provided AliceHandler
func (c *Client) CreateHandler(fn AliceHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if val := recover(); val != nil {
				c.logger(fmt.Errorf("Unexpected error: %v", val))
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		w.Header().Set("Content-Type", "application/json")
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

	var err error
	var i InputData
	if err = json.NewDecoder(r.Body).Decode(&i); err != nil {
		return &AliceHandlerError{fmt.Sprintf("Error while decoding Alice request: %v", err), http.StatusBadRequest}
	}

	var o OutputData
	switch {
	case c.autoPings && i.Request.IsPing():
		o = pong(i)
	case c.autoDanderousContext && i.Request.IsDangerousContext():
		o = dangerous(i)
	default:
		o, err = fn(i)
		if err != nil {
			c.logger(err)
		}
	}

	if err = json.NewEncoder(w).Encode(o); err != nil {
		return &AliceHandlerError{fmt.Sprintf("Error marshaling response: %v", err), http.StatusInternalServerError}
	}

	return nil
}
