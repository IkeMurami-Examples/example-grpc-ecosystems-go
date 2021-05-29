package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/hako/durafmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger logs request/response pair.
func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// log if log level Debug of Trace
		l := zerolog.GlobalLevel()
		if l != zerolog.DebugLevel && l != zerolog.TraceLevel {
			h.ServeHTTP(w, r)
			return
		}

		// We do not want to be spammed by Kubernetes health check.
		// Do not log Kubernetes health check.
		// You can change this behavior as you wish.
		if r.Header.Get("X-Liveness-Probe") == "Healthz" {
			h.ServeHTTP(w, r)
			return
		}

		// Prepare fields to log
		var scheme string
		if r.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
		proto := r.Proto
		method := r.Method
		remoteAddr := r.RemoteAddr
		userAgent := r.UserAgent()
		uri := strings.Join([]string{scheme, "://", r.Host, r.RequestURI}, "")

		// Log HTTP request
		log.Debug().
			Str("http-scheme", scheme).
			Str("http-proto", proto).
			Str("http-method", method).
			Str("remote-addr", remoteAddr).
			Str("user-agent", userAgent).
			Str("uri", uri).
			Msg("request started")

		t := time.Now().In(time.UTC)

		h.ServeHTTP(w, r)

		d := time.Now().In(time.UTC).Sub(t)

		log.Debug().
			Str("http-scheme", scheme).
			Str("http-proto", proto).
			Str("http-method", method).
			Str("remote-addr", remoteAddr).
			Str("user-agent", userAgent).
			Str("uri", uri).
			Str("elapsed", durafmt.Parse(d).String()).
			Msg("request completed")
	})
}
