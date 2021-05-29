package middleware

import (
	"fmt"
	"net/http"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog"
)

type JaegerMiddleware struct {
	tracer opentracing.Tracer
}

func NewJaegerMiddleware(tracer opentracing.Tracer) JaegerMiddleware {
	return JaegerMiddleware{tracer: tracer}
}

// Logger logs request/response pair.
func (j *JaegerMiddleware) TraceMiddleware(h http.Handler) http.Handler {
	return nethttp.Middleware(j.tracer, h, nethttp.OperationNameFunc(func(r *http.Request) string {
		return "HTTP " + r.Method + " " + r.RequestURI
	}))
}

type JaegerLoggerAdapter struct {
	Logger zerolog.Logger
}

func (l JaegerLoggerAdapter) Error(msg string) {
	l.Logger.Error().Str("msg", msg).Msg("###")
}

func (l JaegerLoggerAdapter) Infof(msg string, args ...interface{}) {
	l.Logger.Info().Str("msg", fmt.Sprintf(msg, args...)).Msg("###")
}

func (l JaegerLoggerAdapter) Debugf(msg string, args ...interface{}) {
	l.Logger.Debug().Str("msg", fmt.Sprintf(msg, args...)).Msg("###")
}
