package middleware

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type Middleware interface {
	Logging(next http.Handler) http.Handler
}
type middleware struct {
	logger logrus.FieldLogger
	// aiClient    appinsights.Client
	// redisCache  cache.Cacher
	// amqpAdapter amqpadapter.Adapter
}

func NewMiddleware(logger logrus.FieldLogger) Middleware {
	return &middleware{
		logger: logger,
	}
}

// Logging logs the incoming HTTP request & its duration.
func (m *middleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// App insights log
		// m.aiClient.TrackEvent(event)

		// Container log
		// _ = m.logger.Info(
		// 	"id"
		// )
	})
}
