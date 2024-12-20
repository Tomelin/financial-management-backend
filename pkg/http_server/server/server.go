package http_server

import (
	"net/http"

	"github.com/Tomelin/financial-management-backend/internal/core/entity"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/http2"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var cpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "cpu_temperature_celsius",
	Help: "Current temperature of the CPU.",
})

var routerPath *gin.RouterGroup

func init() {
	prometheus.MustRegister(cpuTemp)
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func (s *RestAPI) Run(handle http.Handler) error {

	srv := http.Server{
		Addr:    ":" + s.Config.Port,
		Handler: s.Route.Handler(),
	}

	http2.ConfigureServer(&srv, &http2.Server{})
	s.Route.Use(s.ValidateToken)
	s.Route.Use(s.MiddlewareHeader)

	return srv.ListenAndServe()
}

func (s *RestAPI) RunTLS() error {
	return nil
}

func setHeader(c *gin.Context) {

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	c.Next()
}

func (s *RestAPI) MiddlewareHeader(c *gin.Context) {
	const SkipMiddlewareKey = "skipMiddleware"

	if c.GetBool(SkipMiddlewareKey) {
		c.Next()
		return
	}

	if c.GetHeader("Authorization") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": entity.Error("invalid token", "token", "MiddlewareHeader", entity.ApplicationLayerHandler, entity.ResponseCodeUnauthorized)})
		c.Writer.Flush()
		c.Abort()
		return
	}
	c.Next()
}

func (s *RestAPI) ValidateToken(c *gin.Context) {

	c.Next()
}
