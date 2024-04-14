package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"

	"github.com/unbeman/av-banner-task/internal/models"
	"github.com/unbeman/av-banner-task/internal/storage"
)

var (
	getBannerQuery = `select b.content, b.is_active from "banner" as b 
		inner join "banner_feature_tags" bft on b.id = bft.banner_id
		where bft.feature_id=$1 and bft.tag_id=$2 and (($3::bool is NULL) or (b.is_active=$3))`

	getBannersWithFilterByTagId = `select
    bft.banner_id,
    bft.feature_id,
    binfo.is_active,
    binfo.content,
    binfo.created_at,
    binfo.updated_at,
    lbft.tags
from banner_feature_tags as bft
    inner join banner as binfo on binfo.id = bft.banner_id
    cross join lateral
        (
            select array
            (
                select sbft.tag_id from banner_feature_tags as sbft
                where sbft.banner_id=binfo.id and sbft.feature_id=binfo.feature_id
            ) as tags
        ) as lbft
where ((@feature_id::integer is NULL) or (bft.feature_id=@feature_id)) and bft.tag_id=@tag_id 
limit @limit 
offset @offset`

	getBanners = `select
    binfo.id,
    binfo.feature_id,
    binfo.is_active,
    binfo.content,
    binfo.created_at,
    binfo.updated_at,
    lbft.tags
from banner as binfo
    cross join lateral
        (
            select array
            (
                select sbft.tag_id from banner_feature_tags as sbft
                where sbft.banner_id=binfo.id and sbft.feature_id=binfo.feature_id
            ) as tags
        ) as lbft
where ((@feature_id::integer is NULL) or (binfo.feature_id=$1))
limit @limit 
offset @offset`

	insertBanner = `insert into banner(feature_id, is_active, content) values (@feature_id, @is_active, @content) returning id`

	updateBannerActiveQuery  = `update banner set is_active=$1 where id=$2`
	updateBannerContentQuery = `update banner set content=$1 where id=$2`

	updateBannerFeatureQuery      = `update banner set feature_id=$1 where id=$2`
	updateBannerFeatureInBFTQuery = `update banner_feature_tags set feature_id=$1 where banner_id=$2`

	deleteBannerTagsQuery = `delete from banner_feature_tags where banner_id=$1`

	deleteBannerByIdQuery = `delete from banner where id=$1`
)

type PGStorage struct {
	connection *pgxpool.Pool
}

func NewPG(ctx context.Context, dsn string) (*PGStorage, error) {
	connPool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	pg := &PGStorage{connection: connPool}

	return pg, nil
}

func (p *PGStorage) Shutdown() {
	p.connection.Close()
	log.Info("postgresql connection pool closed")
}

func (p *PGStorage) GetBanner(ctx context.Context, featureId int, tagId int, isActive *bool) (*models.Banner, error) {
	banner := &models.Banner{}

	err := p.connection.QueryRow(ctx, getBannerQuery, featureId, tagId, isActive).Scan(&banner.Content, &banner.IsActive)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("banner with given feature_id (%d) and tag_id (%d): %w", featureId, tagId, storage.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("couldn't get banner: %w", err)
	}
	return banner, nil
}

func (p *PGStorage) GetBanners(ctx context.Context, featureId *int, tagId *int, limit *int, offset *int) (*models.Banners, error) {
	var query string

	if tagId != nil {
		query = getBannersWithFilterByTagId
	} else {
		query = getBanners
	}

	banners := models.Banners{}
	rows, err := p.connection.Query(ctx, query, pgx.NamedArgs{"feature_id": featureId, "tag_id": tagId, "limit": limit, "offset": offset})
	if err != nil {
		return nil, fmt.Errorf("couldn't get banners: %w", err)
	}

	for rows.Next() {
		banner := models.Banner{}
		err = rows.Scan(&banner.Id, &banner.FeatureId, &banner.IsActive, &banner.Content, &banner.CreatedAt, &banner.UpdateAt, &banner.TagIds)
		if err != nil {
			return nil, fmt.Errorf("couldn't scan banner: %w", err)
		}
		banners = append(banners, &banner)
	}

	return &banners, nil
}

func (p *PGStorage) CreateBanner(ctx context.Context, banner *models.Banner) (*models.Banner, error) {
	tx, err := p.connection.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	err = tx.QueryRow(
		ctx,
		insertBanner,
		pgx.NamedArgs{
			"feature_id": banner.FeatureId,
			"is_active":  banner.IsActive,
			"content":    banner.Content,
		},
	).Scan(&banner.Id)
	if err != nil {
		return nil, fmt.Errorf("couldn't insert banner: %w", err)
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"banner_feature_tags"},
		[]string{"banner_id", "feature_id", "tag_id"},
		pgx.CopyFromSlice(len(banner.TagIds), func(i int) ([]interface{}, error) {
			return []interface{}{banner.Id, banner.FeatureId, banner.TagIds[i]}, nil
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("couldn't insert given set of feature and tags : %w", checkConflictErr(err))

	}

	return banner, err
}

func (p *PGStorage) updateBannerFeaturesAndTags(ctx context.Context, tx pgx.Tx, banner *models.UpdateBannerInput) error {
	switch {
	case banner.FeatureId != nil && banner.TagIds != nil: // обновить и фичу и тэги
		_, err := tx.Exec(ctx, updateBannerFeatureQuery, banner.FeatureId, banner.Id) // обвновляем фичу в баннере
		if err != nil {
			return fmt.Errorf("can't update banner feature_id: %w", err)
		}

		_, err = tx.Exec(ctx, deleteBannerTagsQuery, banner.Id) // удаляем тэги баннера
		if err != nil {
			return fmt.Errorf("can't delete banner tags: %w", err)
		}

		_, err = tx.CopyFrom( // вставляем новые теги
			ctx,
			pgx.Identifier{"banner_feature_tags"},
			[]string{"banner_id", "feature_id", "tag_id"},
			pgx.CopyFromSlice(len(*banner.TagIds), func(i int) ([]interface{}, error) {
				return []interface{}{banner.Id, banner.FeatureId, (*banner.TagIds)[i]}, nil
			}),
		)
		if err != nil {
			return fmt.Errorf("couldn't insert new tags: %w", checkConflictErr(err))
		}

	case banner.FeatureId == nil && banner.TagIds != nil: // обновить только тэги
		_, err := tx.Exec(ctx, deleteBannerTagsQuery, banner.Id) // удаляем тэги баннера
		if err != nil {
			return fmt.Errorf("can't delete banner tags: %w", err)
		}

		_, err = tx.CopyFrom( // вставляем новые теги
			ctx,
			pgx.Identifier{"banner_feature_tags"},
			[]string{"banner_id", "feature_id", "tag_id"},
			pgx.CopyFromSlice(len(*banner.TagIds), func(i int) ([]interface{}, error) {
				return []interface{}{banner.Id, banner.FeatureId, (*banner.TagIds)[i]}, nil
			}),
		)
		if err != nil {
			return fmt.Errorf("couldn't insert new tags: %w", checkConflictErr(err))
		}

	case banner.FeatureId != nil && banner.TagIds == nil: // обновить только фичу
		_, err := tx.Exec(ctx, updateBannerFeatureQuery, banner.FeatureId, banner.Id) // обвновляем фичу в баннере
		if err != nil {
			return fmt.Errorf("can't update banner feature_id: %w", err)
		}
		_, err = tx.Exec(ctx, updateBannerFeatureInBFTQuery, banner.FeatureId, banner.Id) // обвновляем фичу связке баннера с тэгами
		if err != nil {
			return fmt.Errorf("couldn't update banner feature_id in relation: %w", checkConflictErr(err))
		}

	}
	return nil
}

func (p *PGStorage) updateBannerInfo(ctx context.Context, tx pgx.Tx, banner *models.UpdateBannerInput) error {
	if banner.IsActive != nil {
		_, err := tx.Exec(ctx, updateBannerActiveQuery, *banner.IsActive, banner.Id)
		if err != nil {
			return fmt.Errorf("can't update banner active: %w", err)
		}
	}

	if banner.Content != nil {
		_, err := tx.Exec(ctx, updateBannerContentQuery, *banner.Content, banner.Id)
		if err != nil {
			return fmt.Errorf("can't update banner content: %w", err)
		}
	}
	return nil
}

func (p *PGStorage) UpdateBanner(ctx context.Context, banner *models.UpdateBannerInput) error {
	tx, err := p.connection.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	err = p.updateBannerInfo(ctx, tx, banner)
	if err != nil {
		return err
	}

	err = p.updateBannerFeaturesAndTags(ctx, tx, banner)
	if err != nil {
		return err
	}

	return nil
}

func (p *PGStorage) DeleteBanner(ctx context.Context, bannerId int) error {
	tx, err := p.connection.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	result, err := tx.Exec(ctx, deleteBannerTagsQuery, bannerId)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("banner with given id (%d): %w", bannerId, storage.ErrNotFound)
	}
	result, err = tx.Exec(ctx, deleteBannerByIdQuery, bannerId)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("banner with given id (%d): %w", bannerId, storage.ErrNotFound)
	}
	return nil
}

func checkConflictErr(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.ForeignKeyViolation {
			return storage.ErrNotFound
		}
		if pgErr.Code == pgerrcode.UniqueViolation {
			return storage.ErrConflict
		}
	}
	return err
}

func (p PGStorage) Ping(ctx context.Context) error {
	return p.connection.Ping(ctx)
}
