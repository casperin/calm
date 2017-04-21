package calm

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Testing adapters
func TestRateHandlerAdapter(t *testing.T) {
	c := Config{}
	if c.rateHandler != nil {
		t.Error("confige.ratehandler should be nil")
	}
	hf := func(w http.ResponseWriter, r *http.Request) {}
	RateHandler(hf)(&c)
	if c.rateHandler == nil {
		t.Error("Should set rate handler")
	}
}

func TestLookupAdapter(t *testing.T) {
	c := Config{}
	if c.lookup != nil {
		t.Error("config.lookup should be nil")
	}
	l := func(r *http.Request) string { return "foo" }
	Lookup(l)(&c)
	if c.lookup == nil {
		t.Error("Should set config.lookup")
	}
}

func TestMethodsAdapter(t *testing.T) {
	c := Config{}
	if len(c.methods) != 0 {
		t.Error("config.methods should be empty array")
	}
	Methods("POST", "PUT")(&c)
	if len(c.methods) != 2 {
		t.Error("Should set config.methods")
	}
}

// Testing actual rate limiter
func TestCalmer1(t *testing.T) {
	calmer := New(
		2,
		time.Second,
		Methods("GET"),
		RateHandler(func(w http.ResponseWriter, r *http.Request) {
			// Let's make sure we got the right request
			if r.FormValue("a") != "foo" {
				t.Errorf("Expected to get \"foo\", got %v", r.FormValue("a"))
			}
			w.Write([]byte("blocked"))
		}),
	)
	// our rate limited HandlerFunc
	hf := calmer(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })

	w1 := httptest.NewRecorder()
	r1, _ := http.NewRequest("GET", "/?a=foo", nil)
	hf(w1, r1)
	if w1.Body.String() != "ok" {
		t.Errorf("Expected body to be ok. Got %v", w1.Body.String())
	}

	w2 := httptest.NewRecorder()
	r2, _ := http.NewRequest("GET", "/?a=foo", nil)
	hf(w2, r2)
	if w2.Body.String() != "ok" {
		t.Errorf("Expected body to be ok. Got %v", w2.Body.String())
	}

	w3 := httptest.NewRecorder()
	r3, _ := http.NewRequest("GET", "/?a=foo", nil)
	hf(w3, r3)
	if w3.Body.String() != "blocked" {
		t.Errorf("Expected body to be blocked. Got %v", w3.Body.String())
	}

	w4 := httptest.NewRecorder()
	r4, _ := http.NewRequest("PUT", "/?a=foo", nil)
	hf(w4, r4)
	if w4.Body.String() != "ok" {
		t.Errorf("Expected body to be ok. Got %v", w4.Body.String())
	}

	time.Sleep(time.Second)

	w5 := httptest.NewRecorder()
	r5, _ := http.NewRequest("GET", "/?a=foo", nil)
	hf(w5, r5)
	if w5.Body.String() != "ok" {
		t.Errorf("Expected body to be ok. Got %v", w5.Body.String())
	}
}

func TestCalmer2(t *testing.T) {
	calmer := New(
		2,
		time.Second,
		Methods("GET"),
		Lookup(func(r *http.Request) string {
			return r.FormValue("a")
		}),
		RateHandler(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("blocked"))
		}),
	)
	// our rate limited HandlerFunc
	hf := calmer(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })

	w1 := httptest.NewRecorder()
	r1, _ := http.NewRequest("GET", "/?a=foo", nil)
	hf(w1, r1)
	if w1.Body.String() != "ok" {
		t.Errorf("Expected body to be ok. Got %v", w1.Body.String())
	}

	w2 := httptest.NewRecorder()
	r2, _ := http.NewRequest("GET", "/?a=foo", nil)
	hf(w2, r2)
	if w2.Body.String() != "ok" {
		t.Errorf("Expected body to be ok. Got %v", w2.Body.String())
	}

	w3 := httptest.NewRecorder()
	r3, _ := http.NewRequest("GET", "/?a=foo", nil)
	hf(w3, r3)
	if w3.Body.String() != "blocked" {
		t.Errorf("Expected body to be blocked. Got %v", w3.Body.String())
	}

	w4 := httptest.NewRecorder()
	r4, _ := http.NewRequest("GET", "/?a=bar", nil)
	hf(w4, r4)
	if w4.Body.String() != "ok" {
		t.Errorf("Expected body to be ok. Got %v", w4.Body.String())
	}
}
