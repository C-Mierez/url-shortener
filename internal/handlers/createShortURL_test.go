package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/c-mierez/url-shortener/internal/store"
	mockstore "github.com/c-mierez/url-shortener/internal/store/mock"
	"github.com/stretchr/testify/assert"
)

func TestCreateShortURL(t *testing.T) {

	testCases := []struct {
		name                 string
		createMockResult     store.ShortURL
		payload              string
		expectedStatusCode   int
		expectedBody         []byte
		expectedCreateParams store.CreateShortURLParams
	}{
		{
			name: "Create short URL successfully",
			createMockResult: store.ShortURL{
				Slug:        "abc123",
				Destination: "https://www.google.com",
			},
			payload:            `{"destination": "https://www.google.com"}`,
			expectedStatusCode: http.StatusCreated,
			expectedBody:       []byte(`{"id": 0, "destination": "https://www.google.com", "slug": "abc123"}`),
			expectedCreateParams: store.CreateShortURLParams{
				Destination: "https://www.google.com",
				Slug:        "abc123",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Println("Test case: ", tc.name)

			assert := assert.New(t)

			// Create mock store
			mockStore := &mockstore.MockShortURLStore{}

			mockStore.On("CreateShortURL", tc.expectedCreateParams).Return(tc.createMockResult, nil)

			// Create handler
			handler := NewCreateShortURLHandler(CreateShortURLHandlerParams{
				ShortURLStore: mockStore,
				GenerateSlug: func() string {
					return "abc123"
				},
			})

			request := httptest.NewRequest("POST", "/shorten", strings.NewReader(tc.payload))
			responseRecorder := httptest.NewRecorder()

			/// Call handler
			handler.ServeHTTP(responseRecorder, request)

			response := responseRecorder.Result()
			defer response.Body.Close()

			body, err := io.ReadAll(response.Body)
			assert.NoError(err)

			assert.Equal(tc.expectedStatusCode, response.StatusCode)
			assert.JSONEq(string(tc.expectedBody), string(body))

			mockStore.AssertExpectations(t)

		})
	}
}
