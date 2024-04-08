package storage

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("not found")

type Database interface {
	GetBanner(ctx context.Context, featureId int, tagId int, isActive bool)
	GetBanners(ctx context.Context, featureId int, tagId int, limit int, offset int)
	CreateBanner(ctx context.Context, featureId int, tagIds []int, isActive int, content struct{})
	UpdateBanner(ctx context.Context, bannerId int, featureId int, tagIds []int, isActive int, content struct{})
	DeleteBanner(ctx context.Context, bannerId int)
}
