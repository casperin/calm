package limiter

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRateHandler(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Errorf("Unable to create new HTTP request. Error: %v", err)
	}
	w := httptest.NewRecorder()
	DefaultRateHandler(w, r)
	if w.Code != 429 {
		t.Errorf("Expected status code 429. Got %v", w.Code)
	}
}
