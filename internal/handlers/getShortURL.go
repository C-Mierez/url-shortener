package handlers

import (
	"net/http"

	"github.com/c-mierez/url-shortener/internal/store"
)

type GetShortURLHandler struct {
	shortURLStore store.ShortURLStore
}

type NewGetShortURLHandlerParams struct {
	ShortURLStore store.ShortURLStore
}

func NewGetShortURLHandler(params NewGetShortURLHandlerParams) *GetShortURLHandler {
	return &GetShortURLHandler{
		shortURLStore: params.ShortURLStore,
	}

}

// Handler
func (h *GetShortURLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Path[1:]

	shortURL, err := h.shortURLStore.GetShortURLBySlug(slug)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if shortURL == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.Redirect(w, r, shortURL.Destination, http.StatusMovedPermanently)

}
