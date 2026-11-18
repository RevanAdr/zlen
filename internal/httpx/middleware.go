package httpx

import (
	"log/slog"
	"net/http"
	"time"

	"os"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rw *statusRecorder) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		duration := time.Since(start).Seconds()
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		logger.InfoContext(r.Context(), "request",
			"Method", r.Method,
			"Path", r.URL.Path,
			"Code", rw.status,
			"Duration", duration,
		)
	})
}
