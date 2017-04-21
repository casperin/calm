package calm

import (
	"net/http"
	"time"

	"github.com/casperin/calm/buffer"
	"github.com/casperin/calm/limiter"
	"github.com/casperin/calm/lookup"
)

type Config struct {
	numRequests int
	duration    time.Duration
	ipLookups   []string
	methods     []string
	rateHandler http.HandlerFunc
	lookup      func(r *http.Request) string
}

type Adapter func(*Config)

func New(c int, d time.Duration, adapters ...Adapter) func(http.HandlerFunc) http.HandlerFunc {
	config := Config{
		numRequests: c,
		duration:    d,
		ipLookups:   []string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"},
		methods:     []string{},
		rateHandler: limiter.DefaultRateHandler,
	}
	for _, adapter := range adapters {
		adapter(&config)
	}
	return func(next http.HandlerFunc) http.HandlerFunc {
		requests := map[string]*buffer.Buffer{}
		return func(w http.ResponseWriter, r *http.Request) {
			var key string
			if config.lookup != nil {
				key = config.lookup(r)
			} else {
				key = lookup.RemoteIP(config.ipLookups, r)
			}
			if len(config.methods) > 0 && !contains(config.methods, r.Method) {
				next(w, r)
				return
			}
			if _, ok := requests[key]; !ok {
				requests[key] = buffer.New(config.numRequests)
			}
			if b := requests[key]; b.IsOkay(time.Now(), config.duration) {
				b.AddNow()
				next(w, r)
				return
			}
			config.rateHandler(w, r)
		}
	}
}

// By default this library will give a 429 and a message to a user being rate
// limited. This adapter allows you to handle the rate limited user with your
// own handler.
//
//	calmer := calm.New(1, time.Second, calm.RateHandler(func(w http.ResponseWriter, r *http.Request) {
//		w.Write([]byte("Calm down my friend."))
//	})
func RateHandler(handler http.HandlerFunc) Adapter {
	return func(config *Config) {
		config.rateHandler = handler
	}
}

// Sets a function that takes a request and returns a string unique to that
// user.
//
// For example, if for some reason the user is defined by their query param:
//
//	calmer := calm.New(1, time.Second, calm.Lookup(func(r *http.Request) string {
//		return r.FormValue("userid")
//	})
func Lookup(lookup func(r *http.Request) string) Adapter {
	return func(config *Config) {
		config.lookup = lookup
	}
}

// Sets the methods that should be rate limited.
//
//	calmer := calm.New(1, time.Second, calm.Methods("GET", "PUT"))
//
// Will only rate limit GET and PUT requests.
func Methods(methods ...string) Adapter {
	return func(config *Config) {
		config.methods = methods
	}
}

func contains(strs []string, s string) bool {
	for _, str := range strs {
		if s == str {
			return true
		}
	}
	return false
}
