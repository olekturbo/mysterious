package api

import (
	"errors"
	"net/http"

	"github.com/olekturbo/mysterious/internal/service"
)

type Handler struct {
	idService     *service.ID
	cacheService  *service.Cache
	cookieManager *CookieManager
}

func NewHandler(idService *service.ID, cacheService *service.Cache, cookieManager *CookieManager) *Handler {
	return &Handler{
		idService:     idService,
		cacheService:  cacheService,
		cookieManager: NewCookieManager(),
	}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	id, err := h.cookieManager.Read(r, idCookie)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			h.cookieManager.Write(w, idCookie, h.idService.GenerateID())
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	err = h.cacheService.Limit(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusTooManyRequests)
		return
	}

	_, _ = w.Write([]byte("OK"))
}
