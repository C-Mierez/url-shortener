package store

import "errors"

var ErrShortURLNotFound = errors.New("short URL not found")

type ShortURL struct {
	ID          int    `json:"id"`
	Destination string `json:"destination"`
	Slug        string `json:"slug"`
}

type CreateShortURLParams struct {
	Destination string
	Slug        string
}

type ShortURLStore interface {
	CreateShortURL(params CreateShortURLParams) (ShortURL, error)
	GetShortURLBySlug(slug string) (*ShortURL, error)
}
