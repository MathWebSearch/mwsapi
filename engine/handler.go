package engine

import "net/http"

// Handler is a handler using an engine
type Handler interface {

	// Name returns the name of this handler and the appropriate path
	Name() string

	// Enabled checks if this handler is a no-op
	Enabled() bool

	// Connect connects this handler
	// which may be a no-op if no settings were provided
	Connect() (err error)

	// ServeHTTP handles a given request
	ServeHTTP(w http.ResponseWriter, r *http.Request) (code int, res interface{}, err error)
}
