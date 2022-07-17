package internalhttp

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/alexzvon/hw-zvonlexa/hw12_13_14_15_calendar/internal/concat"
	"github.com/alexzvon/hw-zvonlexa/hw12_13_14_15_calendar/internal/logger"
)

func loggingMiddleware(logg logger.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next(w, r)

		strInfo := concat.ConCat(
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

		if err := logg.LogHTTPInfo(strInfo); err != nil {
			log.Fatalln(err)
		}
	}
}
