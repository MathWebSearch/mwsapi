package mwsengine

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MathWebSearch/mwsapi/connection"
	"github.com/MathWebSearch/mwsapi/engine"
	"github.com/MathWebSearch/mwsapi/query"
	"github.com/pkg/errors"
)

// MWSHandler represents an http server capable of answering queries to mws
type MWSHandler struct {
	Host       string
	Port       int
	connection *connection.MWSConnection
}

// Connect connects this handler
func (handler *MWSHandler) Connect() (err error) {
	if handler.Host != "" {
		handler.connection, err = connection.NewMWSConnection(handler.Port, handler.Host)
		if err == nil {
			err = connection.Connect(handler.connection)
		}
		if err != nil {
			err = errors.Wrap(err, "Failed to connect to MWS")
			return
		}
		log.Printf("Connected to MWS at %s", handler.connection.URL())
	}

	return
}

// Name returns the name of this daemon
func (handler *MWSHandler) Name() string {
	return "mws"
}

// Enabled checks if the MWSDaemon is enabled
func (handler *MWSHandler) Enabled() bool {
	return handler.connection != nil
}

// ServeHTTP implements handling of a request
func (handler *MWSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (code int, res interface{}, err error) {
	var request MWSAPIRequest

	// decode the request
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&request)
	if err != nil {
		code = http.StatusBadRequest
		err = errors.Wrap(err, "Failed to read query from request body")
		return
	}

	// default the size to 10
	if request.Size == 0 {
		request.Size = 10
	}

	// validate the request
	err = request.Validate()
	if err != nil {
		code = http.StatusBadRequest
		err = errors.Wrap(err, "Failed to read query from request body")
		return
	}

	// if we have no daemon connection, return
	if !handler.Enabled() {
		code = http.StatusNotFound
		err = errors.Errorf("MathWebSearch is not enabled on this server")
		return
	}

	// run the query, or the count
	if request.Count {
		res, err = Count(handler.connection, request.MWSQuery)
	} else {
		res, err = Run(handler.connection, request.MWSQuery, request.From, request.Size)
	}

	// if there was an error, wrap it and set the status code
	if err != nil {
		code = http.StatusInternalServerError
		err = errors.Wrap(err, "Query failed")
	}

	return
}

// MaxRequestSize is the maximum request size supported by the api
const MaxRequestSize = 100

// MWSAPIRequest is a request sent to mws via the REST API
type MWSAPIRequest struct {
	*query.MWSQuery

	Count bool  `json:"count,omitempty"`
	From  int64 `json:"from,omitempty"`
	Size  int64 `json:"size,omitempty"`
}

// Validate validates an MWSAPI Request
func (req *MWSAPIRequest) Validate() error {
	if req.From < 0 {
		return errors.Errorf("Expected \"from\" to be non-negative, but got %d", req.From)
	}

	if req.Size < 0 || req.Size > MaxRequestSize {
		return errors.Errorf("Expected \"size\" to be between 0 and %d (inclusive), but got %d", MaxRequestSize, req.Size)
	}

	return nil
}

func init() {
	var _ engine.Handler = engine.Handler((*MWSHandler)(nil))
}
