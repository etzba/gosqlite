package server

import (
	"context"
	"net/http"
	"sync"
)

var mu sync.RWMutex

type ctxKey struct{}

// auth is a middleware that check user session token with the values in the session store
func (s *Server) auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session")
		if err != nil {
			s.Respoder.SendUnauthorized(w)
			return
		}

		mu.RLock()
		uid, ok := s.sessions[c.Value]
		mu.RUnlock()
		if !ok {
			s.Respoder.SendUnauthorized(w)
			return
		}
		next(w, r.WithContext(context.WithValue(r.Context(), ctxKey{}, uid)))
	}
}

// log is a middleware that logs all requests to the server
func (s *Server) log(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Info("Request to " + r.RequestURI + " with method: " + r.Method + " accepted. User agent " + r.UserAgent() + " from address " + r.RemoteAddr)
		next(w, r.WithContext(context.WithValue(r.Context(), r.Method, r.RequestURI)))
	}
}
