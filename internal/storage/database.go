package storage

import (
	"context"
	"github.com/unbeman/av-banner-task/internal/models"
)

type Database interface {
	GetBanner(ctx context.Context, featureId int, tagId int, isActive *bool) (*models.Banner, error)
	GetBanners(ctx context.Context, featureId *int, tagId *int, limit *int, offset *int) (*models.Banners, error)
	CreateBanner(ctx context.Context, banner *models.Banner) (*models.Banner, error)
	UpdateBanner(ctx context.Context, banner *models.UpdateBannerInput) error
	DeleteBanner(ctx context.Context, bannerId int) error

	Shutdown()
}
