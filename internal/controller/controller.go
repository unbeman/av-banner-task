package controller

import (
	"context"
	"github.com/unbeman/av-banner-task/internal/models"
	"github.com/unbeman/av-banner-task/internal/storage"
)

type Controller struct {
	database storage.Database
}

func NewController(db storage.Database) (*Controller, error) {
	ctrl := &Controller{database: db}
	return ctrl, nil
}

func (c *Controller) GetBannerForUser(ctx context.Context, input *models.GetBannerInput) (*models.GetBannerOutput, error) {
	isActive := true
	banner, err := c.database.GetBanner(ctx, input.FeatureId, input.TagId, &isActive)
	if err != nil {
		return nil, err
	}
	return (*models.GetBannerOutput)(&banner.Content), nil
}

func (c *Controller) GetBannerForAdmin(ctx context.Context, input *models.GetBannerInput) (*models.GetBannerOutput, error) {
	banner, err := c.database.GetBanner(ctx, input.FeatureId, input.TagId, nil)
	if err != nil {
		return nil, err
	}
	return (*models.GetBannerOutput)(&banner.Content), nil
}
