package api

import (
	"bytes"
	"compress/gzip"
	"io"

	"net/http"
	"strings"
	"time"

	"github.com/DeneesK/metrics-alerting/internal/bodyhasher"
	"go.uber.org/zap"
)

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

type compressWriter struct {
	http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w,
		gzip.NewWriter(w),
	}
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *compressWriter) Close() error {
	return c.zw.Close()
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func gzipMiddleware(log *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			acceptEncoding := r.Header.Get("Accept-Encoding")
			supportsGzip := strings.Contains(acceptEncoding, "gzip")
			if supportsGzip {
				cw := newCompressWriter(w)
				cw.Header().Add("Content-Encoding", "gzip")
				w = cw
				defer cw.Close()
			}
			contentEncoding := r.Header.Get("Content-Encoding")
			sendsGzip := strings.Contains(contentEncoding, "gzip")
			if sendsGzip {
				cr, err := newCompressReader(r.Body)
				if err != nil {
					log.Errorf("during compression error ocurred - %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				r.Body = cr
				defer cr.Close()
			}
			next.ServeHTTP(w, r)
		})
	}
}

func withLogging(log *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			responseData := &responseData{
				status: 0,
				size:   0,
			}
			lw := loggingResponseWriter{
				ResponseWriter: w,
				responseData:   responseData,
			}
			start := time.Now()
			h.ServeHTTP(&lw, r)
			duration := time.Since(start)
			log.Infoln(
				"uri", r.RequestURI,
				"method", r.Method,
				"status", responseData.status,
				"duration", duration,
				"size", responseData.size,
			)
		})
	}
}

func checkHash(log *zap.SugaredLogger, key string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				log.Errorf("during reading body error ocurred - %w", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			req.Body.Close()
			var ha string
			if ha = req.Header.Get("HashSHA256"); ha == "" {
				log.Error("empty header HashSHA256")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			hs, err := bodyhasher.CalculateHash(bodyBytes, key)
			if err != nil {
				log.Errorf("hash calculation failed - %w", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if strings.Compare(hs, ha) != 0 {
				log.Errorf("hashes must be equal - %w", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			h.ServeHTTP(w, req)
		})
	}
}
