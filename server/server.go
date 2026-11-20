package server

import (
	"net/http"

	"github.com/etzba/gosqlite/dat"
	"github.com/etzba/gosqlite/pkg/env"
	"github.com/etzba/gosqlite/pkg/logger"
	"github.com/etzba/gosqlite/wire"
	"github.com/gorilla/mux"
)

// staticSession is Authorization Bearer token that is set in the server to allow authentication with a token for tests
const staticSession = "c4906878493faa9029a77abbb887578dd79963ef5048afb125152035ec4b233e"

// Server serve http with gorilla mux router and http service to end users
type Server struct {
	// Logger logs events
	Logger     *logger.Log
	HTTPServer *http.Server
	Mux        *http.ServeMux
	// Responder is taking care of responding to end user with status codes and payload
	Respoder wire.Responder
	// ctx is the database context that creates changes in the database
	ctx *dat.Context
	// session is the session store to authenticate users
	sessions map[string]int64
}

// New returns a new instance of Server
func New(logger *logger.Log, address string) (*Server, error) {
	responder := wire.HttpResponder{
		Logger: logger,
	}
	server := &Server{
		Logger:   logger,
		Respoder: responder,
		sessions: make(map[string]int64),
	}

	if env.ReleaseStage == "testing" {
		server.sessions[staticSession] = 9999999999999
	}

	ctx, err := dat.NewContext(logger)
	if err != nil {
		return nil, err
	}
	server.ctx = ctx

	router := server.getRouter()
	server.Mux = http.NewServeMux()
	server.Mux.Handle("/", router)
	server.HTTPServer = &http.Server{
		Addr:    address,
		Handler: router,
	}
	return server, nil
}

// Run the server
func (s *Server) Run() error {
	if err := s.ctx.InitDB(s.Logger); err != nil {
		s.Logger.Error("Failed to initilize database", err)
	}

	s.Logger.Info("Start listening in port 8080")
	if err := s.HTTPServer.ListenAndServe(); err != nil {
		s.Logger.Error("cannot run http server - listen and serve", err)
	}
	return nil
}

// getRouter creates an instance of mux router with endpoints
func (s *Server) getRouter() *mux.Router {
	router := mux.NewRouter()
	router.NotFoundHandler = notFound
	router.MethodNotAllowedHandler = methodNotAllowed

	router.HandleFunc("/signup", s.log(s.signup())).Methods("POST")
	router.HandleFunc("/login", s.log(s.login())).Methods("POST")
	router.HandleFunc("/percentiles", s.log(s.calculatePercentileInteger())).Methods("POST")
	router.HandleFunc("/profile/{id}", s.log(s.auth(s.profile()))).Methods("GET", "PUT")
	router.HandleFunc("/headers", s.log(s.auth(s.getHeaders()))).Methods("GET")
	router.HandleFunc("/headers/alpha", s.log(s.auth(s.sortHeadersAlphaBetically()))).Methods("GET")
	router.HandleFunc("/headers/aahlp", s.log(s.auth(s.sortHeadersCharactersAlphaBetically()))).Methods("GET")
	return router
}

var notFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Not found")) //nolint:errcheck
})

var methodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("Method not allowed")) //nolint:errcheck
})
