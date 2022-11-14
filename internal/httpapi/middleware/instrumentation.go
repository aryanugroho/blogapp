package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aryanugroho/blogapp/internal/statsd"
)

func InstrumentStatsD() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
			t1 := time.Now().UnixNano() / int64(time.Millisecond)

			next.ServeHTTP(wr, req)

			t2 := time.Now().UnixNano() / int64(time.Millisecond)
			diff := t2 - t1
			_ = statsd.Gauge(fmt.Sprintf("%s.%s", req.Method, req.URL.EscapedPath()), float64(diff))
		})
	}
}
