package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCheckMethodAndContentType(t *testing.T) {
	// Создаем тестовый HTTP-сервер с применением middleware.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	middlewareHandler := CheckMethodAndContentType(handler)

	// Создаем тестовый запрос с разными вариантами входных данных.
	testCases := []struct {
		method       string
		contentType  string
		requestURI   string
		expectedCode int
		expectedBody string
	}{
		{http.MethodPost, "text/plain", "/update/counter/someMetric/527", http.StatusOK, "OK"},
		{http.MethodGet, "text/plain", "/update/counter/someMetric/527", http.StatusMethodNotAllowed, "Only POST requests are allowed!"},
		{http.MethodPost, "application/json", "/update/counter/someMetric/527", http.StatusMethodNotAllowed, "Only Content-Type is text/plain!"},
		{http.MethodPost, "text/plain", "/update/counter/someMetric/527/metric1", http.StatusNotFound, "Invalid URL format!"},
	}

	for _, tc := range testCases {
		t.Run("Test "+tc.method+" "+tc.contentType+" "+tc.requestURI, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.requestURI, nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", tc.contentType)

			recorder := httptest.NewRecorder()
			middlewareHandler.ServeHTTP(recorder, req)

			if recorder.Code != tc.expectedCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedCode, recorder.Code)
			}

			if body := strings.TrimSpace(recorder.Body.String()); body != tc.expectedBody {
				t.Errorf("Expected response body %q, got %q", tc.expectedBody, body)
			}
		})
	}
}
