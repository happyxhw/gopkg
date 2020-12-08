package grpc

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcZap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
)

func Server(logger *zap.Logger, metricsAddr string) *grpc.Server {
	s := grpc.NewServer(
		grpc.StreamInterceptor(grpcMiddleware.ChainStreamServer(
			grpcZap.StreamServerInterceptor(logger),
			grpcRecovery.StreamServerInterceptor(),
			grpcPrometheus.StreamServerInterceptor,
		)),
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			grpcZap.UnaryServerInterceptor(logger),
			grpcRecovery.UnaryServerInterceptor(),
			grpcPrometheus.UnaryServerInterceptor,
		)),
	)
	if metricsAddr != "" {
		grpcPrometheus.Register(s)
		http.Handle("/metrics", promhttp.Handler())
		go func() {
			err := http.ListenAndServe(metricsAddr, nil)
			if err != nil {
				logger.Error("failed to start metrics")
			}
		}()
	}
	return s
}
