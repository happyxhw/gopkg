package grpc

import (
	"time"

	grpcRetry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func Client(maxRetry uint) (*grpc.ClientConn, error) {
	opts := []grpcRetry.CallOption{
		grpcRetry.WithBackoff(grpcRetry.BackoffLinear(100 * time.Millisecond)),
		grpcRetry.WithMax(maxRetry),
	}
	conn, err := grpc.Dial(viper.GetString("grpc.greeter"),
		grpc.WithStreamInterceptor(grpcRetry.StreamClientInterceptor(opts...)),
		grpc.WithUnaryInterceptor(grpcRetry.UnaryClientInterceptor(opts...)),
		grpc.WithInsecure(),
	)
	return conn, err
}
