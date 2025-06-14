package api

import (
	"errors"
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
			setCookie(w, idCookie, h.service.GenerateID())
		}
	}

	err = h.service.Limit(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusTooManyRequests)
		return
	}

	_, _ = w.Write([]byte("OK"))
}
