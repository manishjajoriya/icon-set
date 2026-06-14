package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/manishjajoriya/icon-set/internal/config"
	"github.com/manishjajoriya/icon-set/internal/service"
	"github.com/manishjajoriya/icon-set/internal/util"
	"github.com/rs/zerolog/log"
)

type IconHandler struct {
	svc service.IconService
	cfg config.Config
}

func NewIconHandler(svc service.IconService, cfg config.Config) *IconHandler {
	return &IconHandler{
		svc: svc,
		cfg: cfg,
	}
}

func (h *IconHandler) GetAll(c *gin.Context) {
	all := h.svc.GetAll(c)
	c.Header(
		"Cache-Control",
		fmt.Sprintf(
			"public, max-age=%v, immutable",
			h.cfg.Http.CacheControlDay*24*60*60,
		),
	)
	c.JSON(http.StatusOK, all)
}

func (h *IconHandler) GetSorted(c *gin.Context) {
	res, err := util.LoadIconList(&h.cfg)
	if err != nil {
		log.Err(err).Msg("IconHandler GetSorted")
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *IconHandler) SearchIcon(c *gin.Context) {
	query := c.Query("query")
	svg := h.svc.SearchIcon(c, query)
	c.JSON(200, svg)
}

func (h *IconHandler) GetIcon(c *gin.Context) {
	key := c.Param("key")

	icon := h.svc.GetIcon(c, key)

	if icon.Key == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	width := icon.Width
	height := icon.Height

	if width == -1 {
		width = h.cfg.Icon.DefaultSize
	}

	if height == -1 {
		height = h.cfg.Icon.DefaultSize
	}

	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %v %v">%s</svg>`, width, height, icon.Body)

	c.Header(
		"Cache-Control",
		fmt.Sprintf(
			"public, max-age=%v, immutable",
			h.cfg.Http.CacheControlDay*24*60*60,
		),
	)
	c.Data(200, "image/svg+xml; charset=utf-8", []byte(svg))
}
