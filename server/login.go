package server

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (s *Server) login() func(w http.ResponseWriter, r *http.Request) { //nolint:unused
	return func(w http.ResponseWriter, r *http.Request) {
		var c creds
		json.NewDecoder(r.Body).Decode(&c)
		var id int64
		var hash string
		if err := s.ctx.Db.QueryRow(`SELECT id, password_hash FROM users WHERE username=?`, c.Username).Scan(&id, &hash); err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(c.Password)) != nil {
			s.Respoder.SendUnauthorized(w)
			return
		}

		t := newToken()
		mu.Lock()
		s.sessions[t] = id
		mu.Unlock()
		http.SetCookie(w, &http.Cookie{Name: "session", Value: t, Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode})
		s.Respoder.SendOK(w, t)
		return
	}
}

func newToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
