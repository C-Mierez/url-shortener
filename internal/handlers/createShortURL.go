package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/c-mierez/url-shortener/internal/store"
)

type CreateShortURLHandler struct {
	shortURLStore store.ShortURLStore
	generateSlug  func() string
}

type CreateShortURLHandlerParams struct {
	ShortURLStore store.ShortURLStore
	GenerateSlug  func() string
}

func NewCreateShortURLHandler(params CreateShortURLHandlerParams) *CreateShortURLHandler {
	return &CreateShortURLHandler{
		shortURLStore: params.ShortURLStore,
		generateSlug:  params.GenerateSlug,
	}
}

func (h *CreateShortURLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	requestData := struct {
		Destination string `json:"destination"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	slug := h.generateSlug()

	newShortURL, err := h.shortURLStore.CreateShortURL(store.CreateShortURLParams{
		Destination: requestData.Destination,
		Slug:        slug,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(newShortURL); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
