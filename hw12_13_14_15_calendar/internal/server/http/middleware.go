package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
)

func midWarePreProc(log *logger.Logger, next http.Handler) http.Handler {
	return mwBeforeStart(log,
		mwAuth(
			mwWrapRW(log, next),
		))
}

func midWarePostProc(log *logger.Logger, r *http.Request) {
	mwd := getMWData(r)
	mwAfterStop(log, &mwd)
}

type midWareData struct {
	uuid    string
	user    string
	start   time.Time
	stop    time.Time
	latency time.Duration
}

type midWareDataKey string

const keyMWData midWareDataKey = "keyMWData"

func getMWData(r *http.Request) midWareData {
	ctx := r.Context()
	value, ok := ctx.Value(keyMWData).(midWareData)
	if !ok {
		value = midWareData{user: "Anonimus"}
	}
	return value
}

func withMWData(r *http.Request, mwd midWareData) *http.Request {
	ctx := context.WithValue(r.Context(), keyMWData, mwd)
	return r.WithContext(ctx)
}

func mwBeforeStart(log *logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mwd := getMWData(r)
		mwd.start = time.Now()
		mwd.uuid = uuid.New().String()
		log.Debug("MidWare.BeforeStart",
			fmt.Sprintf("{%s}: New request [%s] via %s  from %s to %s by %s",
				mwd.uuid, r.Method, r.Proto, r.RemoteAddr, r.URL.Path, r.UserAgent()),
			5,
		)
		next.ServeHTTP(w, withMWData(r, mwd))
	})
}

func mwAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := (*r).Header["User"]
		if !ok {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		mwd := getMWData(r)
		mwd.user = user[0]
		next.ServeHTTP(w, withMWData(r, mwd))
	})
}

type wrappedRW struct {
	http.ResponseWriter
	statusCode int
	writeBuf   []byte
}

func newWrappedRW(w http.ResponseWriter) *wrappedRW {
	return &wrappedRW{w, -1, []byte{}}
}

func (wrw *wrappedRW) WriteHeader(code int) {
	wrw.statusCode = code
	wrw.ResponseWriter.WriteHeader(code)
}

func (wrw *wrappedRW) Write(buf []byte) (int, error) {
	wrw.writeBuf = buf
	return wrw.ResponseWriter.Write(buf)
}

func mwWrapRW(log *logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrw := newWrappedRW(w)
		next.ServeHTTP(wrw, r)
		mwd := getMWData(r)
		log.Debug("MidWare.WrapRW",
			fmt.Sprintf("{%s}: User[%s] response status [%d]", mwd.uuid, mwd.user, wrw.statusCode),
			3,
		)
		log.Debug("MidWare.WrapRW",
			fmt.Sprintf("{%s}: User[%s] response body [%s]", mwd.uuid, mwd.user, wrw.writeBuf),
			5,
		)
	})
}

func mwAfterStop(log *logger.Logger, mwd *midWareData) {
	mwd.stop = time.Now()
	mwd.latency = mwd.stop.Sub(mwd.start)
	log.Debug("MidWare.AfterStop",
		fmt.Sprintf("{%s}: User[%s] responce latency %s", mwd.uuid, mwd.user, mwd.latency),
		3,
	)
}
