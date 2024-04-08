package pg

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	log "github.com/sirupsen/logrus"
	"github.com/unbeman/av-banner-task/internal/storage"
)

type Statements struct {
	GetBanner    *sql.Stmt
	GetBanners   *sql.Stmt
	CreateBanner *sql.Stmt
	DeleteBanner *sql.Stmt
}

func NewStatements(conn *sql.DB) (Statements, error) {
	var err error
	s := Statements{}
	s.GetBanners, err = conn.Prepare(
		`select b.id, b.content 
				from banner as b 
    			join banner_tags bt on b.id = bt.banner_id 
                where feature_id=$1 and bt.tag_id=$2 and b.is_active=$3;`)
	if err != nil {
		return s, err
	}
	s.GetBanners, err = conn.Prepare("")
	if err != nil {
		return s, err
	}
	s.CreateBanner, err = conn.Prepare("")
	if err != nil {
		return s, err
	}
	s.DeleteBanner, err = conn.Prepare("")
	if err != nil {
		return s, err
	}
	return s, nil
}

type pgStorage struct {
	connection *sql.DB
	statements Statements
}

func NewPG(dsn string) (storage.Database, error) {
	connection, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	pg := &pgStorage{connection: connection}

	pg.statements, err = NewStatements(connection)
	if err != nil {
		return nil, err
	}
	return pg, nil
}

func (p *pgStorage) Shutdown() error {
	err := p.connection.Close()
	log.Infoln("db conn closed")
	return err
}

func (p *pgStorage) GetBanner(ctx context.Context, featureId int, tagId int, isActive bool) {
	//TODO implement me
	panic("implement me")
}

func (p *pgStorage) GetBanners(ctx context.Context, featureId int, tagId int, limit int, offset int) {
	//TODO implement me
	panic("implement me")
}

func (p *pgStorage) CreateBanner(ctx context.Context, featureId int, tagIds []int, isActive int, content struct{}) {
	//TODO implement me
	panic("implement me")
}

func (p *pgStorage) UpdateBanner(ctx context.Context, bannerId int, featureId int, tagIds []int, isActive int, content struct{}) {
	//TODO implement me
	panic("implement me")
}

func (p *pgStorage) DeleteBanner(ctx context.Context, bannerId int) {
	//TODO implement me
	panic("implement me")
}
