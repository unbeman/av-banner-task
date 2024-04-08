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

func (c *Controller) GetBanner(ctx context.Context, input *models.GetBannerInput) (*models.GetBannerOutput, error) {
	panic("implement me")
	return nil, nil
}
