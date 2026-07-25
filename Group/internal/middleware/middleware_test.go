package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoggingMiddleware(t *testing.T) {
	// 1. Create a dummy handler that acts as our "next" application route.
	// We tell it to explicitly return a 404 status to test our interceptor.
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound) // 404 Error
		w.Write([]byte("not found"))
	})

	// 2. Wrap the dummy handler inside our Logging Middleware
	handlerToTest := LoggingMiddleware(testHandler)

	// 3. Create a fake HTTP request hitting a dummy path
	req := httptest.NewRequest("GET", "/test-route", nil)

	// 4. Create a ResponseRecorder to act as a fake browser client
	recorder := httptest.NewRecorder()

	// 5. Execute the middleware handler using the fake request and recorder
	handlerToTest.ServeHTTP(recorder, req)

	// 6. Assertions: Verify that the response code matches what we expect
	if recorder.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, but got %d", http.StatusNotFound, recorder.Code)
	}

	// Verify that the body text was passed through successfully
	expectedBody := "not found"
	if recorder.Body.String() != expectedBody {
		t.Errorf("Expected body %q, but got %q", expectedBody, recorder.Body.String())
	}
}
