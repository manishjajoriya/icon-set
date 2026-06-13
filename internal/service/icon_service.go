package service

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/manishjajoriya/icon-set/internal/types"
	"github.com/manishjajoriya/icon-set/internal/util"
)

type IconService interface {
	GetAll(Ctx *gin.Context) []types.SearchResponse
	SearchIcon(ctx context.Context, query string) []types.SearchResponse
	GetIcon(ctx context.Context, key string) util.Icon
}

type iSvc struct {
	icons map[string]util.Icon
}

func NewIconService(icons map[string]util.Icon) IconService {
	return &iSvc{
		icons: icons,
	}
}

func (s *iSvc) GetAll(Ctx *gin.Context) []types.SearchResponse {
	icons := make([]types.SearchResponse, 0, len(s.icons))
	for _, icon := range s.icons {
		icons = append(icons, types.SearchResponse{
			Key:    icon.Key,
			Name:   icon.Name,
			Prefix: icon.Prefix,
		})
	}
	return icons
}

func (s *iSvc) SearchIcon(ctx context.Context, query string) []types.SearchResponse {
	query = strings.ToLower(query)

	matches := make([]types.SearchResponse, 0, 50)

	for _, icon := range s.icons {
		if !strings.Contains(strings.ToLower(icon.Key), query) {
			continue
		}

		matches = append(matches, types.SearchResponse{
			Key:    icon.Key,
			Name:   icon.Name,
			Prefix: icon.Prefix,
		})

		if len(matches) >= 50 {
			break
		}
	}

	return matches
}

func (s *iSvc) GetIcon(ctx context.Context, key string) util.Icon {
	return s.icons[key]
}
