package grpcserver

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/logger"
	"google.golang.org/grpc"
)

func loggingInterceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		result, err := handler(ctx, req)

		logger.Info(
			fmt.Sprintf("%s %s",
				method(info.FullMethod),
				latency(start),
			))

		return result, err
	}
}

func method(full string) string {
	return full[strings.LastIndex(full, "/"):]
}

func latency(start time.Time) string {
	return fmt.Sprintf("%dms", time.Since(start).Milliseconds())
}
