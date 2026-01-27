package server

import (
	"compress/gzip"
	"github.com/kornetvba/metrics-service/internal/config/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	writer io.Writer
}

func (gw *gzipWriter) Write(data []byte) (int, error) {
	return gw.writer.Write(data)
}

type gzReader struct {
	r      io.ReadCloser
	reader *gzip.Reader
}

func (gr *gzReader) Read(data []byte) (int, error) {
	return gr.reader.Read(data)
}

func (gr *gzReader) Close() error {
	return gr.reader.Close()
}

func unCompressData(r io.ReadCloser) (*gzReader, error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &gzReader{reader: gz, r: r}, nil

}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			gzipNewWriter, err := gzip.NewWriterLevel(w, gzip.BestCompression)
			if err != nil {
				logger.Logger.Info("gzip", zap.Error(err))
			}

			gw := &gzipWriter{
				ResponseWriter: w,
				writer:         gzipNewWriter,
			}
			defer gzipNewWriter.Close()
			w = gw

		}
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gz, err := unCompressData(r.Body)
			if err != nil {
				logger.Logger.Info("gzip", zap.Error(err))
			}
			r.Body = gz

		}
		next.ServeHTTP(w, r)

	})

}
