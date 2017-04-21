package lookup

import (
	"net/http"
	"strings"
)

func RemoteIP(ipLookups []string, r *http.Request) string {
	forwardedFor := r.Header.Get("X-Forwarded-For")
	realIP := r.Header.Get("X-Real-IP")

	for _, lookup := range ipLookups {
		if lookup == "RemoteAddr" {
			idx := strings.LastIndex(r.RemoteAddr, ":")
			if idx == -1 {
				return r.RemoteAddr
			}
			return r.RemoteAddr[:idx]
		}

		if lookup == "X-Forwarded-For" && forwardedFor != "" {
			parts := strings.Split(forwardedFor, ",")
			return strings.TrimSpace(parts[0])
		}

		if lookup == "X-Real-IP" && realIP != "" {
			return realIP
		}
	}

	return ""
}
