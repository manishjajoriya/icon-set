package main

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/manishjajoriya/icon-set/internal/config"
	"github.com/manishjajoriya/icon-set/internal/handler"
	"github.com/manishjajoriya/icon-set/internal/middleware"
	"github.com/manishjajoriya/icon-set/internal/service"
	"github.com/manishjajoriya/icon-set/internal/util"
	"github.com/rs/zerolog/log"
)

type app struct {
	icons map[string]util.Icon
	cfg   *config.Config
}

func (a *app) Init() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
		c.Next()
	})
	r.Use(middleware.ZerologMiddleware())
	r.Use(gin.Recovery())

	iconService := service.NewIconService(a.icons)
	iconHandler := handler.NewIconHandler(iconService, *a.cfg)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		v1.GET("/usage", func(c *gin.Context) {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			c.JSON(200, gin.H{
				"Alloc":     fmt.Sprintf("%.2f MB", float64(m.Alloc)/1024/1024),
				"HeapAlloc": fmt.Sprintf("%.2f MB", float64(m.HeapAlloc)/1024/1024),
				"HeapSys":   fmt.Sprintf("%.2f MB", float64(m.HeapSys)/1024/1024),
				"NumGC":     fmt.Sprintf("%d", m.NumGC),
			})

		})
		v1.GET("/all", iconHandler.GetAll)
		v1.GET("/search", iconHandler.SearchIcon)
		v1.GET("/icon/:key", iconHandler.GetIcon)
	}

	return r
}

func (a *app) Run(r *gin.Engine) {
	log.Info().Msgf("Running server on port %v", a.cfg.Http.Port)
	err := r.Run(":" + a.cfg.Http.Port)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
