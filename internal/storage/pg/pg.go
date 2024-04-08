package pg

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
	"github.com/unbeman/av-banner-task/internal/models"
	"github.com/unbeman/av-banner-task/internal/storage"
)

type pgStorage struct {
	connection *pgxpool.Pool
}

func NewPG(ctx context.Context, dsn string) (storage.Database, error) {
	connPool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	pg := &pgStorage{connection: connPool}

	return pg, nil
}

func (p *pgStorage) Shutdown() {
	p.connection.Close()
	log.Info("db conn pool closed")
}

func (p *pgStorage) GetBanner(ctx context.Context, featureId int, tagId int, isActive *bool) (*models.Banner, error) {
	const query = `select b.id, array_agg(bt.tag_id), b.content from "banner" as b 
		join "banner_tags" bt on b.id = bt.banner_id
		where b.feature_id=$1 and bt.tag_id=$2 and b.is_active=COALESCE($3, b.is_active)
		group by b.id
		limit 1`

	banner := &models.Banner{}

	err := p.connection.QueryRow(ctx, query, featureId, tagId, isActive).Scan(&banner.Id, &banner.TagIds, &banner.Content)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%w", storage.ErrNotFound)
	}
	if err != nil {
		return nil, err // todo: обернуть кишки
	}
	return banner, nil
}

func (p *pgStorage) GetBanners(ctx context.Context, featureId int, tagId int, limit int, offset int) (*models.Banners, error) {
	//TODO implement me
	panic("implement me")
}

func (p *pgStorage) CreateBanner(ctx context.Context, featureId int, tagIds []int, isActive int, content struct{}) (*models.Banner, error) {
	//TODO implement me
	panic("implement me")
}

func (p *pgStorage) UpdateBanner(ctx context.Context, bannerId int, featureId int, tagIds []int, isActive int, content struct{}) error {
	//TODO implement me
	panic("implement me")
}

func (p *pgStorage) DeleteBanner(ctx context.Context, bannerId int) error {
	//TODO implement me
	panic("implement me")
}
