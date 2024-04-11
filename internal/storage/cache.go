package storage

import "context"

type Cache interface {
	GetBanner(ctx context.Context, featureId, tagId int) (*string, error)
	SetBanner(ctx context.Context, featureId, tagId int, bannerContent *string) error

	Shutdown()
}
