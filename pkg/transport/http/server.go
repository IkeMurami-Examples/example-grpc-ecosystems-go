// Package http implements HTTP-specific transport features like initialization, logging, etc.
package http

import (
	"net/http"

	"github.com/opentracing/opentracing-go"

	"github.com/z3f1r/grpc-ecosystems-example/pkg/transport/http/middleware"
)

func NewServer(endpoint string, mux http.Handler, tracer opentracing.Tracer) *http.Server {
	jaeger := middleware.NewJaegerMiddleware(tracer)
	return &http.Server{
		Addr: endpoint,
		// add handler with middleware
		Handler: jaeger.TraceMiddleware(
			middleware.Logger(mux)),
	}
}
