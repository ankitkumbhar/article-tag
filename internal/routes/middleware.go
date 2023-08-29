package routes

import (
	"article-tag/internal/handler"
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap/zapcore"
)

func LogRequest(app *handler.Application) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := handler.GetLogger(app)

			var res map[string]interface{}

			bodyBytes, _ := io.ReadAll(r.Body)
			r.Body.Close() //  must close
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			json.Unmarshal(bodyBytes, &res)

			// log incoming request details
			logger.Debug("request", zapcore.Field{Key: "path", Type: zapcore.StringType, String: r.URL.Path},
				zapcore.Field{Key: "method", Type: zapcore.StringType, String: r.Method},
				zapcore.Field{Key: "body", Type: zapcore.ReflectType, Interface: res})

			next.ServeHTTP(w, r)
		})
	}
}
