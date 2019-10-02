package engine

import (
	"github.com/json-iterator/go"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// Server is a server capable of answering requests to multiple engines
type Server struct {
	router   *mux.Router
	handlers []Handler
}

// NewServer makes a new server
func NewServer() (server *Server) {
	server = &Server{
		router: mux.NewRouter(),
	}
	server.init()
	return
}

// init intializes the server
func (server *Server) init() {
	// add the root route
	server.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !server.ensureMethod(w, r, "GET") {
			return
		}

		server.jsonResponse(w, r, http.StatusOK, server.Status())
	})

	// add a custom 404 route
	server.router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.errorResponse(w, r, http.StatusNotFound, errors.Errorf("Not found"))
	})
}

// StatusResponse is the response for a status request
type StatusResponse struct {
	// Name of this server
	Name string `json:"name"`
	// Server tagline
	TagLine string `json:"tagline"`

	// List of engines that are enabled on this server.
	Engines map[string]bool `json:"engines"`
}

// Status returns the current status
func (server *Server) Status() *StatusResponse {
	engines := make(map[string]bool)

	for _, h := range server.handlers {
		engines[h.Name()] = h.Enabled()
	}

	return &StatusResponse{
		Name:    "mwsapid",
		TagLine: "You know, for math",
		Engines: engines,
	}
}

// AddHandler adds a handler to the server
func (server *Server) AddHandler(handler Handler) {
	// store the handler
	server.handlers = append(server.handlers, handler)

	// and add an actual route
	server.router.HandleFunc("/"+handler.Name()+"/", func(w http.ResponseWriter, r *http.Request) {
		if !server.ensureMethod(w, r, "POST") {
			return
		}

		code, res, err := handler.ServeHTTP(w, r)
		if code == 0 {
			code = 200
		}
		if err != nil {
			server.errorResponse(w, r, code, err)
		} else {
			server.jsonResponse(w, r, code, res)
		}
	})
}

// ListenAndServe starts serving requests on the given port
func (server *Server) ListenAndServe(host string, port int) error {
	return http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), server.logRequest(server.router))
}

// ensureMethod ensures that the method `method` is used in the request
// returns true if it is so
// returns false and sends an error response if it is not so
func (server *Server) ensureMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		server.errorResponse(w, r, http.StatusMethodNotAllowed, errors.Errorf("Invalid request type: Expected %s but got %s", method, r.Method))
		return false
	}

	return true
}

// jsonResponse returns a json or jsonp response
func (server *Server) jsonResponse(w http.ResponseWriter, r *http.Request, code int, response interface{}) {
	// marshal the response
	data, err := jsoniter.Marshal(response)
	if err != nil {
		data = []byte(`{"error":true,"message":"Unknown error"}`)
		code = http.StatusInternalServerError
	}

	// write the code
	w.WriteHeader(code)

	// check if we have a JSONP callback parameter
	callback := r.URL.Query().Get("callback")
	if callback == "" {
		w.Header().Set("Content-Type", "application/json")
	} else {
		w.Header().Set("Content-Type", "application/javascript")
		data = []byte(fmt.Sprintf("%s(%s);", callback, data))
	}

	w.Write(data)
}

// errorResponse returns an error response
func (server *Server) errorResponse(w http.ResponseWriter, r *http.Request, code int, err error) {
	server.jsonResponse(w, r, code, fmt.Sprintf("%v", err))
}

func (server *Server) logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
