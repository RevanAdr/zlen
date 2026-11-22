package httpx

import (
	"crypto/ecdsa"
	"log/slog"
	"net/http"
	"time"

	"os"

	"github.com/golang-jwt/jwt/v5"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

type UserClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type Jwt struct {
	key    *ecdsa.PrivateKey
	token  *jwt.Token
	string string
}

func (rw *statusRecorder) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func MakeToken(userId string, secret []byte) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userId,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
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

func Authentication(key []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("access_token")
			if err != nil {
				WriteToJSON(w, http.StatusUnauthorized, "Unauthorized", nil)
				return
			}

			claims := &jwt.RegisteredClaims{}
			_, tErr := jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (any, error) { return key, nil },
				jwt.WithValidMethods([]string{"HS256"}),
				jwt.WithExpirationRequired(),
			)

			if tErr != nil {
				WriteToJSON(w, http.StatusUnauthorized, "Invalid Token", nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
