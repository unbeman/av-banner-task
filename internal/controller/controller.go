package controller

import (
	"context"
	"errors"
	"github.com/unbeman/av-banner-task/internal/models"
	"github.com/unbeman/av-banner-task/internal/storage"
)

type Controller struct {
	database storage.Database
	cache    storage.Cache
}

func NewController(db storage.Database, cache storage.Cache) (*Controller, error) {
	ctrl := &Controller{database: db, cache: cache}
	return ctrl, nil
}

func (c *Controller) GetBanner(ctx context.Context, input *models.GetBannerInput, isActive *bool) (*models.GetBannerOutput, error) {
	var bannerContent *models.GetBannerOutput

	if input.UseLastRevision {
		banner, err := c.database.GetBanner(ctx, input.FeatureId, input.TagId, isActive)
		if err != nil {
			return nil, err
		}
		bannerContent = (*models.GetBannerOutput)(&banner.Content)

		return bannerContent, nil
	} else {
		content, err := c.cache.GetBanner(ctx, input.FeatureId, input.TagId)

		if errors.Is(err, storage.ErrNotFound) {
			banner, err := c.database.GetBanner(ctx, input.FeatureId, input.TagId, isActive)
			if err != nil {
				return nil, err
			}
			bannerContent = (*models.GetBannerOutput)(&banner.Content)

			if banner.IsActive { // добавляем в кэш только активные баннеры
				if err = c.cache.SetBanner(ctx, input.FeatureId, input.TagId, &banner.Content); err != nil {
					return nil, err
				}
			}

			return bannerContent, nil
		}
		if err != nil {
			return nil, err
		}
		bannerContent = (*models.GetBannerOutput)(content)

		return bannerContent, nil
	}
}

func (c *Controller) GetBanners(ctx context.Context, input *models.GetBannersInput) (*models.Banners, error) {
	return c.database.GetBanners(ctx, input.FeatureId, input.TagId, input.Limit, input.Offset)
}

func (c *Controller) CreateBanner(ctx context.Context, input *models.Banner) (*models.CreateBannerOutput, error) {
	banner, err := c.database.CreateBanner(ctx, input)
	if err != nil {
		return nil, err
	}
	bannerOut := &models.CreateBannerOutput{BannerId: banner.Id}
	return bannerOut, nil
}

func (c *Controller) UpdateBanner(ctx context.Context, input *models.UpdateBannerInput) error {
	return c.database.UpdateBanner(ctx, input)
}

func (c *Controller) DeleteBanner(ctx context.Context, bannerId int) error {
	return c.database.DeleteBanner(ctx, bannerId)
}
