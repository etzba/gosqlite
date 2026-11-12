package server

import (
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	_ "modernc.org/sqlite"
)

// creds are credentials of a user
type creds struct{ Username, Password string }

// user of the app
type user struct {
	id       int64
	Username string
	Created  time.Time
}

// signup will create a user in the database and allow login afterwards
func (s *Server) signup() func(w http.ResponseWriter, r *http.Request) { //nolint:unused
	return func(w http.ResponseWriter, r *http.Request) {
		var c creds
		if json.NewDecoder(r.Body).Decode(&c) != nil || c.Username == "" || c.Password == "" {
			s.Respoder.SendBadRequest(w)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
		if err != nil {
			s.Respoder.SendError(w, err)
			return
		}

		res, err := s.ctx.Db.Exec(`INSERT INTO users(username, password_hash) VALUES(?,?)`, c.Username, hash)
		if err != nil {
			s.Respoder.SendConflict(w, "Username is already taken")
			return
		}

		uid, err := res.LastInsertId()
		if err != nil {
			s.Respoder.SendError(w, err)
			return
		}

		pr, err := s.ctx.Db.Exec(`INSERT INTO profiles(user_id) VALUES(?)`, uid)
		if err != nil {
			s.Respoder.SendError(w, err)
			return
		}

		pid, err := pr.LastInsertId()
		if err != nil {
			s.Respoder.SendError(w, err)
			return
		}

		_, err = s.ctx.Db.Exec(`UPDATE users SET profile_id=? WHERE id=?`, pid, uid)
		if err != nil {
			s.Respoder.SendError(w, err)
			return
		}

		createdUser := user{
			id:       uid,
			Username: c.Username,
			Created:  time.Now(),
		}
		s.Respoder.SendCreated(w, createdUser)
		return
	}
}
