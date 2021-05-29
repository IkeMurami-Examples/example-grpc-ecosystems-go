package example

import (
	"context"
	"fmt"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"

	"github.com/opentracing/opentracing-go"
	//"github.com/slok/go-http-metrics/middleware"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"time"

	"google.golang.org/grpc"
	// "github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	// metrics "github.com/slok/go-http-metrics/metrics/prometheus"

	// t_http_gw_metrics "github.com/z3f1r/grpc-ecosystems-example/pkg/transport/grpc-gw/metric"
	// "github.com/z3f1r/grpc-ecosystems-example/pkg/transport/grpc/clients"
	// "github.com/z3f1r/grpc-ecosystems-example/pkg/transport/grpc/interceptors"
	t_http "github.com/z3f1r/grpc-ecosystems-example/pkg/transport/http"
	t_http_mw "github.com/z3f1r/grpc-ecosystems-example/pkg/transport/http/middleware"
	utils "github.com/z3f1r/grpc-ecosystems-example/pkg/utils"
)

const (
	serverShutdownTimeout = 5
)

var (
	methods = map[string]*desc.MethodDescriptor{}
)

func StartServers(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	interrupt := make(chan os.Signal, 1)

	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	g, ctx := errgroup.WithContext(ctx)

	_ = os.Setenv("JAEGER_SERVICE_NAME", "example")
	jaegerCfg, err := jaegercfg.FromEnv()
	if err != nil {
		return fmt.Errorf("failed to parse Jaeger config: %s", err.Error())
	}
	jaegerLogger := t_http_mw.JaegerLoggerAdapter{Logger: log.Logger}
	jaegerCfg.Reporter.LogSpans = viper.GetBool("DEBUG")
	tracer, closer, err := jaegerCfg.NewTracer(
		jaegercfg.Logger(jaegerLogger),
	)
	if err != nil {
		return fmt.Errorf("could not initialize jaeger tracer: %s", err.Error())
	}
	defer closer.Close()

	sentryDSN := viper.GetString("SENTRY_DSN")
	if sentryDSN != "" {
		utils.SentrySetup(sentryDSN)
	}

	// start gRPC-Gateway server
	grpcGateway, err := startGRPCGateway(ctx, g, tracer)
	if err != nil {
		return fmt.Errorf("failed to start gRPC-Gateway: %v", err)
	}

	toolServer := startToolServer(ctx, g)

	select {
	case <-interrupt:
		break
	case <-ctx.Done():
		break
	}
	log.Warn().Msg("received shutdown signal")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(),
		serverShutdownTimeout*time.Second)
	defer shutdownCancel()

	_ = grpcGateway.Shutdown(shutdownCtx)
	_ = toolServer.Shutdown(shutdownCtx)

	return g.Wait()
}

// startGRPCGateway starts gRPC bridge.
func startGRPCGateway(ctx context.Context, g *errgroup.Group, tracer opentracing.Tracer) (*http.Server, error) {
	var (
		err             error
	)

	endpoint := viper.GetString("HTTP_ENDPOINT")

	/*
	mdlwMetrics := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	 */

	gsrv := grpc.NewServer()

	methods, err = getMethodDescriptors(gsrv)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	/*
	Здесь можно подключить другие сервисы в mux, клиентами которых является этот сервис
	 */

	srv := t_http.NewServer(endpoint, mux, tracer)

	g.Go(func() error {
		log.Info().Str("endpoint", endpoint).Msg("gRPC-Gateway server serving")

		err := srv.ListenAndServe()
		if err != nil {
			log.Warn().Err(err).
				Str("endpoint", endpoint).
				Msg("gRPC-Gateway server stopped")
		}
		return err
	})

	return srv, nil
}

func getMethodDescriptors(srv *grpc.Server) (map[string]*desc.MethodDescriptor, error) {
	var methods = map[string]*desc.MethodDescriptor{}
	sds, err := grpcreflect.LoadServiceDescriptors(srv)
	if err != nil {
		log.Error().Err(err).Msg("grpcreflect.LoadServiceDescriptors")
		return methods, err
	}
	var out = methods
	for _, sd := range sds {
		for _, md := range sd.GetMethods() {
			methodName := fmt.Sprintf("/%s/%s", sd.GetFullyQualifiedName(), md.GetName())
			out[methodName] = md
		}
	}
	return out, nil
}

// startToolServer starts tool Server.
func startToolServer(ctx context.Context, g *errgroup.Group) *http.Server {
	endpoint := viper.GetString("TOOL_SERVER_ENDPOINT")
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	srv := &http.Server{
		Handler: mux,
		Addr:    endpoint,
	}
	g.Go(func() error {
		log.Info().Str("endpoint", endpoint).Msg("tool server serving")

		err := srv.ListenAndServe()
		if err != nil {
			log.Warn().Err(err).
				Str("endpoint", endpoint).
				Msg("Tool server stopped")
		}
		return err
	})
	return srv
}