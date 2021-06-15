package server

import (
	"net/http"
)

// adapter wraps an http.Handler with additional functionality.
type adapter func(http.Handler) http.Handler

// adapt will iterate over all adapters, calling them one
// by one (in reverse order) in a chained manner, returning the result
// of the first adapter.
func adapt(h http.Handler, adapters ...adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func (s *server) logging() adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.logger.Infof("%s %s", r.Method, r.RequestURI)
			h.ServeHTTP(w, r)
		})
	}
}
