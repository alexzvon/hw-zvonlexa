package internalhttp

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/alexzvon/hw-zvonlexa/hw12_13_14_15_calendar/internal/logger"
	"github.com/alexzvon/hw-zvonlexa/hw12_13_14_15_calendar/internal/myutils"
)

func loggingMiddleware(logg logger.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next(w, r)

		strInfo := myutils.ConCat(
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
		)

		logg.LogHTTPInfo(strInfo)
	}
}
