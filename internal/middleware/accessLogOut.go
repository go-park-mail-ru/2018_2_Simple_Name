package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

func AccessLogMiddleware(mux *http.ServeMux, sugar *zap.SugaredLogger, HitStat *prometheus.GaugeVec) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()

		mux.ServeHTTP(w, r)

		HitStat.With(prometheus.Labels{
			"url":    r.URL.Path,
			"method": r.Method,
			"code":   w.Header().Get("Status-Code"),
		}).Inc()

		sugar.Infow(r.URL.Path,
			"method", r.Method,
			"remote addres", r.RemoteAddr,
			"url", r.URL.Path,
			"work time", time.Since(begin),
			"time now", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
		//"time now", time.Now().UTC())
	})
}

//func ErrorLog (msg string, err error, sugar *zap.SugaredLogger){
//	sugar.Infow(msg,
//		"error", err,
//		"time", strconv.Itoa(time.Now().Hour()) + ":" + strconv.Itoa(time.Now().Minute()))
//}
