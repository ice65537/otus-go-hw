package internalhttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
)

type midWareData struct {
	log     *logger.Logger
	user    string
	latency time.Duration
}

func midWare(log *logger.Logger, next http.Handler) http.Handler {
	mwd := new(midWareData)
	mwd.log = log
	return mwLog(mwd,
		mwLatency(mwd,
			mwAuth(mwd, next),
		))
}

func mwLog(mwd *midWareData, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		//
		mwd.log.Debug("MW.Log",
			fmt.Sprintf("%s req [%s] from %s to %s by %s",
				r.Proto, r.Method, r.RemoteAddr, r.URL.Path, r.UserAgent()),
			4,
		)
		mwd.log.Debug("MW.Log",
			fmt.Sprintf("Resp {%v} to %s with latency %s", w.Header(), mwd.user, mwd.latency),
			4,
		)
	})
}

func mwLatency(mwd *midWareData, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t0 := time.Now()
		//
		next.ServeHTTP(w, r)
		//
		mwd.latency = time.Since(t0)
	})
}

func mwAuth(mwd *midWareData, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user, ok := (*r).Header["User"]; ok {
			mwd.user = user[0]
		}
		if mwd.user == "" {
			w.WriteHeader(http.StatusForbidden)
		} else {
			//
			next.ServeHTTP(w, r)
		}
	})
}
