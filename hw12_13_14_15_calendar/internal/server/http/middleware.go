package internalhttp

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/logger"
)

func loggingMiddleware(logg logger.Logger, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next(w, r)

		strInfo := []string{
			r.RemoteAddr,
			" [",
			start.Format(time.RFC822),
			"] ",
			r.Method,
			" ",
			r.URL.String(),
			" ",
			strconv.Itoa(r.Response.StatusCode),
			" ",
			fmt.Sprintf("%f", time.Since(start).Seconds()),
			" \"",
			r.UserAgent(),
			"\" \n",
		}

		logg.LogHTTPInfo(strInfo)
	})
}
