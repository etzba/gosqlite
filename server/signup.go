package server

import (
	"net/http"
)

func (s *Server) signup() func(w http.ResponseWriter, r *http.Request) { //nolint:unused
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Info("Server go request" + " method: " + r.Method + " uri: " + r.RequestURI)
		s.Respoder.SendOK(w)
	}
}
