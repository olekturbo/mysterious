package api

import (
	"errors"
	"fmt"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	id, err := readCookie(r, idCookie)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			newId := h.service.GenerateID()
			err = h.service.Set(r.Context(), newId, fmt.Sprintf("Hello, %s!", newId))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))
				return
			}

			setCookie(w, idCookie, newId)
			_, _ = w.Write([]byte(newId))
		}
	}

	val, err := h.service.Get(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	_, _ = w.Write([]byte(val))
}
