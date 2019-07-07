package galice

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
//func (c *Client) CreateHandler(fn AliceHandler) http.HandlerFunc
