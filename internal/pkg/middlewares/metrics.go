package middlewares

import (
	"github.com/prometheus/client_golang/prometheus"
	mylog "github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"strings"
)

var hits = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "hits",
	Help: "count hits"}, []string{"status", "path"})
var errors = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "errors",
	Help: "count errors"}, []string{"status", "path"})
var durationHist = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "durationHist",
	Help:    "histogram of requests duration in seconds",
	Buckets: []float64{1e-6, 1e-5, 1e-4, 1e-3, 5e-3, 0.01, 0.025, 0.1, 0.5, 1, 2, 5, 10}}, []string{"path"})

type MetricsProm struct {
}

func RegisterMetrics () {
	prometheus.MustRegister(hits, errors, durationHist)
}

func InitMetrics() *MetricsProm {
	return &MetricsProm{}
}

func (m *MetricsProm) Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := r.RequestURI
		i:= strings.Index(uri, "?")
		if i > 0 {
			uri = uri[:i]+"param"
		}
		urlSplit := strings.Split(uri, "/")
		_, err := strconv.Atoi(urlSplit[len(urlSplit)-1])
		if err == nil {
			urlSplit[len(urlSplit)-1] = "id"
		}
		url := strings.Join(urlSplit, "/")


		customWriter := &statusWriter{
			ResponseWriter: w,
		}
		//start time
		timerHist := prometheus.NewTimer(durationHist.WithLabelValues(url))
		next.ServeHTTP(customWriter, r)
		//stop time
		if s := timerHist.ObserveDuration().Seconds(); s < 0 {
			mylog.Debug().Msg("negative request duration")
		}
		//record hits and errors
		hits.WithLabelValues(strconv.Itoa(customWriter.Status), url).Inc()
		if customWriter.Status >= 400 {
			errors.WithLabelValues(strconv.Itoa(customWriter.Status), url).Inc()
		}
	})
}
type statusWriter struct {
	http.ResponseWriter
	Status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.Status = status
	w.ResponseWriter.WriteHeader(http.StatusOK)
}

