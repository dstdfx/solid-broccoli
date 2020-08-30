package v1

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	RequestIDHeader     = "x-request-id"
	RefererHeader       = "referer"
	UserAgentHeader     = "user-agent"
	XRealIPHeader       = "x-real-ip"
	XForwardedForHeader = "x-forwarded-for"
)

type ctxKey int

const (
	ctxRequestID ctxKey = iota
	ctxLogger
)

// SetRequestID middleware creates a new request ID and saves it into request context.
func SetRequestID(log *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u, err := uuid.NewV4()
			if err != nil {
				log.Error("failed ti generate UUID", zap.Error(err))
				http.Error(w, "", http.StatusInternalServerError)

				return
			}

			w.Header().Set(RequestIDHeader, u.String())
			ctx := context.WithValue(r.Context(), ctxRequestID, u.String())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetRequestID gets a request ID from context or returns an empty string.
func GetRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(ctxRequestID).(string); ok {
		return reqID
	}

	return ""
}

// RequestLogger handles logging of additional information about every request.
func RequestLogger(log *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Clone logger to not use fields that will be set in other handlers
			logClone := *log

			// Retrieve basic request information
			method := r.Method
			path := r.URL.Path
			requestID := GetRequestID(r.Context())
			ipAddr := requestGetRemoteAddress(r)
			userAgent := r.Header.Get(UserAgentHeader)
			referer := r.Header.Get(RefererHeader)

			// Retrieve request body contents
			var bodyBytes []byte
			if r.Body != nil {
				bodyBytes, _ = ioutil.ReadAll(r.Body)
			}

			// Restore the io.ReadCloser to its original state
			r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

			bodyString := string(bodyBytes)

			ww := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Save start time
			start := time.Now().UTC()

			// Proceed to next middleware
			next.ServeHTTP(ww, r)

			defer func() {
				latency := time.Since(start)
				statusCode := ww.Status()
				msg := "Request completed"

				fields := []zapcore.Field{
					zap.Int("status", statusCode),
					zap.Duration("latency", latency),
					zap.String("ip", ipAddr),
					zap.String("method", method),
					zap.String("path", path),
					zap.String(UserAgentHeader, userAgent),
					zap.String(RefererHeader, referer),
					zap.String(RequestIDHeader, requestID),
					zap.String("request_body", bodyString),
					zap.Int("response_bytes_written", ww.BytesWritten()),
				}

				switch {
				case statusCode > 499:
					logClone.Error(msg, fields...)
				case statusCode > 399:
					logClone.Warn(msg, fields...)
				default:
					logClone.Info(msg, fields...)
				}
			}()
		})
	}
}

// SetContextLogger populates additional field of the provided logger and saves it into context.
func SetContextLogger(log *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := log.With(
				zap.String(RequestIDHeader, GetRequestID(r.Context())),
			)

			ctx := context.WithValue(r.Context(), ctxLogger, logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetContextLogger retrieves logger from the provided context.
func GetContextLogger(ctx context.Context) (*zap.Logger, error) {
	if log, ok := ctx.Value(ctxLogger).(*zap.Logger); ok {
		return log, nil
	}

	return nil, errors.New("no logger in request context")
}

// requestGetRemoteAddress returns ip address of the client making the request, taking into account http proxies.
func requestGetRemoteAddress(r *http.Request) string {
	headerRealIP := r.Header.Get(XRealIPHeader)
	headerForwardedFor := r.Header.Get(XForwardedForHeader)

	if headerRealIP == "" && headerForwardedFor == "" {
		return ipAddrFromRemoteAddr(r.RemoteAddr)
	}

	if headerForwardedFor != "" {
		// X-Forwarded-For is potentially a list of addresses separated with ","
		parts := strings.Split(headerForwardedFor, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}

		return parts[0]
	}

	return headerRealIP
}

// Request.RemoteAddress contains port, which we want to remove i.e.: "[::1]:58292" => "[::1]".
func ipAddrFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}

	return s[:idx]
}
