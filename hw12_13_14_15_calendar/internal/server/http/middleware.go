package internalhttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
)

func getReqSession(r *http.Request) logger.Session {
	ctx := r.Context()
	sss := logger.GetCtxSession(ctx)
	return sss
}

func reqWithSession(r *http.Request, sss logger.Session) *http.Request {
	ctx := logger.PushCtxSession(r.Context(), sss)
	return r.WithContext(ctx)
}

func midWareHandler(log *logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sss := getReqSession(r) // generate uuid and fix start timestamp
		r = reqWithSession(r, sss)
		log.Debug(r.Context(), "MidWare.Handler",
			fmt.Sprintf("New request [%s] via [%s] from [%s] to [%s] by [%s]",
				r.Method, r.Proto, r.RemoteAddr, r.URL.Path, r.UserAgent()),
			5,
		)

		var ok bool
		sss.User, ok = mwAuth(r)
		if !ok {
			w.WriteHeader(http.StatusForbidden)
			midWareAfterResponse(log, r, http.StatusForbidden, []byte(""))
			return
		}
		r = reqWithSession(r, sss)
		log.Debug(r.Context(), "MW.Auth", fmt.Sprintf("User [%s] auth success", sss.User), 2)

		wrw := newWrappedRW(w)
		next.ServeHTTP(wrw, r)
		midWareAfterResponse(log, r, wrw.statusCode, wrw.writeBuf)
	})
}

func mwAuth(r *http.Request) (string, bool) {
	user, ok := (*r).Header["User"]
	if !ok {
		return "", false
	}
	return user[0], true
}

func midWareAfterResponse(log *logger.Logger, r *http.Request, retCode int, response []byte) {
	sss := getReqSession(r)
	log.Debug(r.Context(), "MidWare.AfterResponse",
		fmt.Sprintf("Responce status [%d] with latency %s", retCode, time.Since(sss.Start)),
		3,
	)
	log.Debug(r.Context(), "MidWare.AfterResponse",
		fmt.Sprintf("Response body [%s]", response),
		5,
	)
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
