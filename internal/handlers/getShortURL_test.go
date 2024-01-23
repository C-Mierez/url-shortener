package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/c-mierez/url-shortener/internal/store"
	"github.com/c-mierez/url-shortener/internal/store/mock"
	"github.com/stretchr/testify/assert"
)

func TestGetShortURL(t *testing.T) {
	testCases := []struct {
		name                      string
		url                       string
		expectedStatusCode        int
		expectedGetShortURLParams string
		getShortURLResult         *store.ShortURL
	}{
		{
			name:                      "Successfully get short URL from slug",
			url:                       "/abc123",
			expectedStatusCode:        http.StatusMovedPermanently,
			expectedGetShortURLParams: "abc123",
			getShortURLResult: &store.ShortURL{
				ID:          0,
				Destination: "https://www.google.com",
				Slug:        "abc123",
			},
		},
		{
			name:                      "Failed get short URL from slug with 404",
			url:                       "/abc123",
			expectedStatusCode:        http.StatusNotFound,
			expectedGetShortURLParams: "abc123",
			getShortURLResult:         nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			assert := assert.New(t)

			mockShortURLStore := mock.MockShortURLStore{}

			mockShortURLStore.On(
				"GetShortURLBySlug",
				tc.expectedGetShortURLParams).Return(tc.getShortURLResult, nil)

			handler := NewGetShortURLHandler(NewGetShortURLHandlerParams{
				ShortURLStore: &mockShortURLStore,
			})

			request := httptest.NewRequest("GET", tc.url, nil)
			responseRecorder := httptest.NewRecorder()

			handler.ServeHTTP(responseRecorder, request)

			response := responseRecorder.Result()
			defer response.Body.Close()

			assert.Equal(tc.expectedStatusCode, response.StatusCode)

			mockShortURLStore.AssertExpectations(t)

		})
	}
}
