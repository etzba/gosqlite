package wire

import (
	"encoding/json"
	"net/http"

	"github.com/etzba/gosqlite/pkg/logger"
)

type Responder interface {
	SendOK(w http.ResponseWriter, obj ...interface{})
	SendNothing(w http.ResponseWriter)
	SendError(w http.ResponseWriter, err error)
	SendBadRequest(w http.ResponseWriter)
}

type Respond struct {
	Logger *logger.Log
}

func (r Respond) SendOK(w http.ResponseWriter, obj ...interface{}) {
	if obj != nil {
		body, err := json.Marshal(obj)
		if err != nil {
			r.Logger.Error("could not marshal", err)
			w.WriteHeader(http.StatusInternalServerError)
			res := []byte(`{"message":"could not marshal"}`)
			_, err = w.Write(res)
			if err != nil {
				r.Logger.Error("could not write data", err)
			}
		}
		w.Header().Set("Content-Type", "Application/json")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		_, err = w.Write(body)
		if err != nil {
			r.Logger.Error("could not write body", err)
		}
	} else {
		w.Write([]byte("OK!")) //nolint:errcheck
	}
}

func (r Respond) SendNothing(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("No content")) //nolint:errcheck
}

func (r Respond) SendError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error())) //nolint:errcheck
}

func (r Respond) SendBadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Bad request")) //nolint:errcheck
}
