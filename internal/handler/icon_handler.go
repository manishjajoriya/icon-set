package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/manishjajoriya/icon-set/internal/config"
	"github.com/manishjajoriya/icon-set/internal/service"
)

type IconHandler struct {
	svc service.IconService
	cfg config.IconConfig
}

func NewIconHandler(svc service.IconService, cfg config.IconConfig) *IconHandler {
	return &IconHandler{
		svc: svc,
		cfg: cfg,
	}
}

func (h *IconHandler) GetAll(c *gin.Context) {
	all := h.svc.GetAll(c)
	c.JSON(http.StatusOK, all)
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

	if width <= 0 {
		width = h.cfg.DefaultSize
	}

	if height <= 0 {
		height = h.cfg.DefaultSize
	}

	if icon.Height == -1 {
		height = width
	}

	if icon.Width == -1 {
		width = height
	}

	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %v %v">%s</svg>`, width, height, icon.Body)

	c.Data(200, "image/svg+xml; charset=utf-8", []byte(svg))
}
