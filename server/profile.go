package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) profile() func(w http.ResponseWriter, r *http.Request) { //nolint:unused
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			uid := r.Context().Value(ctxKey{}).(int64)
			var out struct {
				ProfileID int64  `json:"profile_id"`
				UserID    int64  `json:"user_id"`
				Posts     string `json:"posts"`
			}
			if err := s.ctx.Db.QueryRow(`SELECT id, user_id, posts FROM profiles WHERE user_id=?`, uid).
				Scan(&out.ProfileID, &out.UserID, &out.Posts); err != nil {
				s.Respoder.SendNotFound(w)
				return
			}

			if err := json.NewEncoder(w).Encode(out); err != nil {
				s.Respoder.SendError(w, err)
				return
			}
			s.Respoder.SendOK(w, out)
			return
		}

		if r.Method == "UPDATE" {
			uid := r.Context().Value(ctxKey{}).(int64)
			var in struct{ Posts string }
			if json.NewDecoder(r.Body).Decode(&in) != nil || len(in.Posts) > 50 {
				s.Respoder.SendBadRequest(w)
				return
			}

			_, err := s.ctx.Db.Exec(`UPDATE profiles SET posts=? WHERE user_id=?`, in.Posts, uid)
			if err != nil {
				s.Respoder.SendError(w, err)
				return
			}

			s.Respoder.SendNothing(w)
			return
		}
	}
}
