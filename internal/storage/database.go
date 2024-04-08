package storage

import (
	"context"
	"errors"
	"github.com/unbeman/av-banner-task/internal/models"
)

var ErrNotFound = errors.New("not found")

type Database interface {
	GetBanner(ctx context.Context, featureId int, tagId int, isActive *bool) (*models.Banner, error)
	GetBanners(ctx context.Context, featureId int, tagId int, limit int, offset int) (*models.Banners, error)
	CreateBanner(ctx context.Context, featureId int, tagIds []int, isActive int, content struct{}) (*models.Banner, error)
	UpdateBanner(ctx context.Context, bannerId int, featureId int, tagIds []int, isActive int, content struct{}) error
	DeleteBanner(ctx context.Context, bannerId int) error
}
