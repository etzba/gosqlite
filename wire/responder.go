package wire

import (
	"encoding/json"
	"net/http"

	"github.com/etzba/gosqlite/pkg/logger"
)

// Responder creates http handler responses according to status codes
type Responder interface {
	SendOK(w http.ResponseWriter, obj ...interface{})
	SendCreated(w http.ResponseWriter, obj ...interface{})
	SendNothing(w http.ResponseWriter)
	SendError(w http.ResponseWriter, err error)
	SendBadRequest(w http.ResponseWriter)
	SendUnauthorized(w http.ResponseWriter)
	SendNotFound(w http.ResponseWriter)
	SendConflict(w http.ResponseWriter, obj ...interface{})
}

type HttpResponder struct {
	Logger *logger.Log
}

// SendOK
func (r HttpResponder) SendOK(w http.ResponseWriter, obj ...interface{}) {
	w.WriteHeader(http.StatusOK)
	if obj != nil {
		r.writeObjectToBody(w, obj...)
	} else {
		w.Write([]byte("OK!")) //nolint:errcheck
	}
}

// SendCreated
func (r HttpResponder) SendCreated(w http.ResponseWriter, obj ...interface{}) {
	w.WriteHeader(http.StatusCreated)
	if obj != nil {
		r.writeObjectToBody(w, obj...)
	} else {
		w.Write([]byte("Created")) //nolint:errcheck
	}
}

// SendNothing
func (r HttpResponder) SendNothing(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("No content")) //nolint:errcheck
}

// SendError
func (r HttpResponder) SendError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error())) //nolint:errcheck
}

// SendUnauthorized
func (r HttpResponder) SendUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("Unauthorized")) //nolint:errcheck
}

// SendBadRequest
func (r HttpResponder) SendBadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Bad request")) //nolint:errcheck
}

// SendNotFound
func (r HttpResponder) SendNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Not found")) //nolint:errcheck
}

// SendConflict
func (r HttpResponder) SendConflict(w http.ResponseWriter, obj ...interface{}) {
	w.WriteHeader(http.StatusConflict)
	if obj != nil {
		r.writeObjectToBody(w, obj...)
	} else {
		w.Write([]byte("Conflict")) //nolint:errcheck
	}
}

// writeObjectToBody
func (r HttpResponder) writeObjectToBody(w http.ResponseWriter, obj ...interface{}) {
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
}
