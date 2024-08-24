package internalhttp

import (
	"fmt"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/logger"
	"net/http"
	"strings"
	"time"
)

func loggingMiddleware(next http.HandlerFunc, logger logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{w, http.StatusOK}
		next(rw, r)

		logger.Info(
			fmt.Sprintf("%s %s %s %s %d %s %s",
				requestAddr(r),
				r.Method,
				r.RequestURI,
				r.Proto,
				rw.code,
				latency(start),
				userAgent(r),
			))
	}
}

func requestAddr(r *http.Request) string {
	return strings.Split(r.RemoteAddr, ":")[0]
}

func userAgent(r *http.Request) string {
	userAgents := r.Header["User-Agent"]
	if len(userAgents) > 0 {
		return "\"" + userAgents[0] + "\""
	}
	return ""
}

func latency(start time.Time) string {
	return fmt.Sprintf("%dms", time.Since(start).Milliseconds())
}
