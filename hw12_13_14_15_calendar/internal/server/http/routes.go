package internalhttp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
)

func (s Server) response(w http.ResponseWriter, r *http.Request, oper string, status int, data string) {
	w.WriteHeader(status)
	if _, err := w.Write([]byte(data)); err != nil {
		_ = s.log.Error(r.Context(), oper, fmt.Sprintf("%v", err))
	}
}

func (s Server) hello(w http.ResponseWriter, r *http.Request) {
	s.response(w, r, "Http.Hello", http.StatusOK, fmt.Sprintf("Hello %s!", getReqSession(r).User))
}

func (s Server) new(w http.ResponseWriter, r *http.Request) {
	s.edit(w, r, "Http.New")
}

func (s Server) reset(w http.ResponseWriter, r *http.Request) {
	s.edit(w, r, "Http.Reset")
}

func (s Server) drop(w http.ResponseWriter, r *http.Request) {
	s.edit(w, r, "Http.Drop")
}

func (s Server) edit(w http.ResponseWriter, r *http.Request, oper string) {
	req, err := io.ReadAll(r.Body)
	if err != nil {
		s.response(w, r, oper, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	evt, err := storage.Unmarshal(req)
	if err != nil {
		s.response(w, r, oper, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	switch {
	case evt.ID == "" && oper == "Http.Drop":
		s.response(w, r, oper, http.StatusBadRequest, "event.ID not found in request")
		return
	case evt.ID != "" && oper == "Http.New":
		s.response(w, r, oper, http.StatusBadRequest, "invalid path for existing event with ID, use /event/reset")
		return
	case evt.Owner != "" && evt.Owner != getReqSession(r).User:
		s.response(w, r, oper, http.StatusUnauthorized,
			fmt.Sprintf("%s, you can't create/update/drop event with Owner=[%s]", getReqSession(r).User, evt.Owner))
		return
	}
	evt.Owner = getReqSession(r).User

	switch {
	case oper == "Http.New" || oper == "Http.Reset":
		if err = s.app.Upsert(r.Context(), evt); err != nil {
			s.response(w, r, oper, http.StatusInternalServerError, fmt.Sprintf("%v", err))
			return
		}
	case oper == "Http.Drop":
		if err = s.app.Drop(r.Context(), evt.ID); err != nil {
			s.response(w, r, oper, http.StatusInternalServerError, fmt.Sprintf("%v", err))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (s Server) get(w http.ResponseWriter, r *http.Request) {
	const oper = "Http.Get"
	var reqArgs appGetEvents

	req, err := io.ReadAll(r.Body)
	if err != nil {
		s.response(w, r, oper, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	if err := json.Unmarshal(req, &reqArgs); err != nil {
		s.response(w, r, oper, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	events, err := s.app.Get(r.Context(), reqArgs.T1, reqArgs.T2)
	if err != nil {
		s.response(w, r, oper, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	data, err := storage.Marshall(events)
	if err != nil {
		s.response(w, r, oper, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	s.response(w, r, oper, http.StatusOK, string(data))
}
