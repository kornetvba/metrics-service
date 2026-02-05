package logger

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func LogMiddlewarePost(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeNow := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(timeNow)

		Logger.Info(
			"info post",
			zap.String("url", fmt.Sprintf("%v", r.URL)),
			zap.String("method", r.Method),
			zap.Duration("time duration", duration),
		)
	})

}

type LogResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	sizeResponse int
}

func (rw *LogResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *LogResponseWriter) Write(data []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(data)
	if err != nil {
		return 0, err
	}
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}

	rw.sizeResponse += size
	return size, nil
}

func LogMiddlewareGet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := &LogResponseWriter{ResponseWriter: w}
		next.ServeHTTP(writer, r)
		Logger.Info(
			"info response",
			zap.Int("status code", writer.statusCode),
			zap.Int("size response", writer.sizeResponse),
		)
	})

}
