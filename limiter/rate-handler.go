package limiter

import "net/http"

func DefaultRateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(429)
	w.Write([]byte("You have reached maximum request limit."))
}
