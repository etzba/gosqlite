package server

import (
	"net/http"

	"github.com/etzba/gosqlite/dat"
	"github.com/etzba/gosqlite/pkg/logger"
	"github.com/etzba/gosqlite/wire"
	"github.com/gorilla/mux"
)

// Server
type Server struct {
	Logger     *logger.Log
	HTTPServer *http.Server
	Mux        *http.ServeMux
	Respoder   wire.Responder
	ctx        *dat.Context
	sessions   map[string]int64
}

// New
func New(logger *logger.Log, address string) (*Server, error) {
	responder := wire.Respond{
		Logger: logger,
	}
	server := &Server{
		Logger:   logger,
		Respoder: responder,
		sessions: make(map[string]int64),
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

// Run
func (s *Server) Run() error {
	if err := dat.InitDB(s.Logger); err != nil {
		s.Logger.Error("Failed to initilize database", err)
	}

	s.Logger.Info("Start listening in port 8080")
	if err := s.HTTPServer.ListenAndServe(); err != nil {
		s.Logger.Error("cannot run http server - listen and serve", err)
	}
	return nil
}

// getRouter
func (s *Server) getRouter() *mux.Router {
	router := mux.NewRouter()
	router.NotFoundHandler = notFound
	router.MethodNotAllowedHandler = methodNotAllowed

	router.HandleFunc("/signup", s.log(s.signup())).Methods("POST")
	router.HandleFunc("/login", s.log(s.login())).Methods("POST")
	router.HandleFunc("/profile/{id}", s.log(s.auth(s.profile()))).Methods("GET", "POST", "PUT")
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
