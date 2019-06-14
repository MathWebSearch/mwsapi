package temaengine

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MathWebSearch/mwsapi/connection"
	"github.com/MathWebSearch/mwsapi/engine"
	"github.com/MathWebSearch/mwsapi/query"
	"github.com/pkg/errors"
)

// TemaHandler represents an http server capable of answering queries to temasearch
type TemaHandler struct {
	MWSPort     int
	MWSHost     string
	ElasticPort int
	ElasticHost string

	connection *connection.TemaConnection
}

// Connect connects this handler
func (handler *TemaHandler) Connect() (err error) {
	if handler.MWSHost != "" && handler.ElasticHost != "" {
		handler.connection, err = connection.NewTemaConnection(handler.MWSPort, handler.MWSHost, handler.ElasticPort, handler.ElasticHost)
		if err == nil {
			err = connection.Connect(handler.connection)
		}
		if err != nil {
			err = errors.Wrap(err, "Failed to connect to TemaSearch")
			return
		}
		log.Printf("Connected to TemaSearch")
	}

	return
}

// Name returns the name of this daemon
func (handler *TemaHandler) Name() string {
	return "tema"
}

// Enabled checks if the MWSDaemon is enabled
func (handler *TemaHandler) Enabled() bool {
	return handler.connection != nil
}

// ServeHTTP implements handling of a request
func (handler *TemaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (code int, res interface{}, err error) {
	var request TemaAPIRequest

	// decode the request
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&request)
	if err != nil {
		code = http.StatusBadRequest
		err = errors.Wrap(err, "Failed to read query from request body")
		return
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
		err = errors.Errorf("TemaSearch is not enabled on this server")
		return
	}

	// run the query, or the count
	if request.Count {
		res, err = Count(handler.connection, request.Query)
	} else {
		res, err = Run(handler.connection, request.Query, request.From, request.Size)
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

// TemaAPIRequest is a request sent to mws via the REST API
type TemaAPIRequest struct {
	*query.Query

	Count bool  `json:"count,omitempty"`
	From  int64 `json:"from,omitempty"`
	Size  int64 `json:"size,omitempty"`
}

// Validate validates an MWSAPI Request
func (req *TemaAPIRequest) Validate() error {
	if req.From < 0 {
		return errors.Errorf("Expected \"from\" to be non-negative, but got %d", req.From)
	}

	if req.Size < 0 || req.Size > MaxRequestSize {
		return errors.Errorf("Expected \"size\" to be between 0 and %d (inclusive), but got %d", MaxRequestSize, req.Size)
	}

	return nil
}

func init() {
	var _ = engine.Handler((*TemaHandler)(nil))
}
