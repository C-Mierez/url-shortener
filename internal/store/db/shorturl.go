package db

import (
	"log/slog"

	"github.com/c-mierez/url-shortener/internal/store"
)

type ShortURLStore struct {
	shortURLs []store.ShortURL
	logger    *slog.Logger
}

// NewShortURLStore
type NewShortURLStoreParams struct {
	Logger *slog.Logger
}

func NewShortURLStore(params NewShortURLStoreParams) *ShortURLStore {

	shortURLs := []store.ShortURL{}

	return &ShortURLStore{
		shortURLs,
		params.Logger,
	}
}

// CreateShortURL
func (s *ShortURLStore) CreateShortURL(params store.CreateShortURLParams) (store.ShortURL, error) {

	shortURL := store.ShortURL{
		ID:          len(s.shortURLs),
		Destination: params.Destination,
		Slug:        params.Slug,
	}

	s.shortURLs = append(s.shortURLs, shortURL)

	s.logger.Info("Created short URL", slog.Any("values", shortURL))

	return shortURL, nil
}

// GetShortURLBySlug
func (s *ShortURLStore) GetShortURLBySlug(slug string) (*store.ShortURL, error) {

	// Iterate over the slice of short URLs
	for _, shortURL := range s.shortURLs {
		if shortURL.Slug == slug {
			return &shortURL, nil
		}
	}

	return nil, store.ErrShortURLNotFound
}
