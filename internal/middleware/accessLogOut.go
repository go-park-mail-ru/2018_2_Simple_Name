package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

func AccessLogMiddleware (mux *http.ServeMux, sugar *zap.SugaredLogger) http.HandlerFunc   {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()

		mux.ServeHTTP(w, r)

		sugar.Infow(r.URL.Path,
			"method", r.Method,
			"remote addres", r.RemoteAddr,
			"url", r.URL.Path,
			"work time", time.Since(begin),
			"time now", strconv.Itoa(time.Now().Hour()) + ":" + strconv.Itoa(time.Now().Minute()))
			//"time now", time.Now().UTC())
	})
}

//func ErrorLog (msg string, err error, sugar *zap.SugaredLogger){
//	sugar.Infow(msg,
//		"error", err,
//		"time", strconv.Itoa(time.Now().Hour()) + ":" + strconv.Itoa(time.Now().Minute()))
//}
