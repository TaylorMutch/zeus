package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func New() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(sloggin.NewWithConfig(
		slog.Default(),
		sloggin.Config{
			DefaultLevel:       slog.LevelInfo,
			ClientErrorLevel:   slog.LevelWarn,
			ServerErrorLevel:   slog.LevelError,
			WithUserAgent:      true,
			WithRequestID:      true,
			WithRequestBody:    false,
			WithRequestHeader:  false,
			WithResponseBody:   false,
			WithResponseHeader: false,
			WithSpanID:         false,
			WithTraceID:        true,
			Filters:            []sloggin.Filter{sloggin.IgnorePath("/liveness", "/readiness")},
		},
	))
	router.Use(otelgin.Middleware("zeus"))
	router.Use(gin.Recovery())
	router.GET("/liveness", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	router.GET("/readiness", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	return router
}
