package server

import (
	"net/http"

	"github.com/etzba/gosqlite/pkg/logger"
	"github.com/etzba/gosqlite/wire"
	"github.com/gorilla/mux"
)

type Server struct {
	Logger     *logger.Log
	HTTPServer *http.Server
	Mux        *http.ServeMux
	Respoder   wire.Responder
}

func New(logger *logger.Log, address string) *Server {
	responder := wire.Respond{
		Logger: logger,
	}
	server := &Server{
		Logger:   logger,
		Respoder: responder,
	}
	router := server.getRouter()
	server.Mux = http.NewServeMux()
	server.Mux.Handle("/", router)
	server.HTTPServer = &http.Server{
		Addr:    address,
		Handler: router,
	}
	return server
}

func (s *Server) Run() error {
	s.Logger.Info("Start server in port 8080")
	if err := s.HTTPServer.ListenAndServe(); err != nil {
		s.Logger.Error("cannot run http server - listen and serve", err)
	}
	return nil
}

func (s *Server) getRouter() *mux.Router {
	router := mux.NewRouter()
	router.NotFoundHandler = notFound
	router.MethodNotAllowedHandler = methodNotAllowed

	router.HandleFunc("/signup", s.signup()).Methods("GET")
	router.HandleFunc("/login", s.login()).Methods("GET")
	router.HandleFunc("/profile/{id}", s.profile()).Methods("GET")
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
