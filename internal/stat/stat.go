package stat

import (
	"github.com/prometheus/client_golang/prometheus"
)

func NewPrometheus() *prometheus.GaugeVec {
	HitStat := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "Simple_Game",
			Subsystem: "hit_stat",
			Name:      "HitStat",
			Help:      "Hit info.",
		},
		[]string{
			"url",
			"method",
			"code",
		},
	)
	prometheus.MustRegister(HitStat)
	return HitStat
}
