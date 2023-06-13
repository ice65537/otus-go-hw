package internalhttp

import (
	"fmt"
	"net/http"

	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
)

func midWare(log *logger.Logger, next http.Handler) http.Handler {
	return mwLog(log, mwAuth(log, next))
}

func mwLog(log *logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("MW.Log", fmt.Sprintf("%v", (*r).Body), 5)
		next.ServeHTTP(w, r)
	})
}

func mwAuth(log *logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("MW.Auth", fmt.Sprintf("%v", (*r).Header), 5)
		next.ServeHTTP(w, r)
	})
}
