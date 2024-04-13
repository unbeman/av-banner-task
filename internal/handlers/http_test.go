package handlers

import (
	"context"
	"fmt"
	"github.com/steinfletcher/apitest"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/unbeman/av-banner-task/internal/config"
	"github.com/unbeman/av-banner-task/internal/controller"
	"github.com/unbeman/av-banner-task/internal/models"
	"github.com/unbeman/av-banner-task/internal/storage/pg"
	"github.com/unbeman/av-banner-task/internal/storage/redis"
	"github.com/unbeman/av-banner-task/internal/utils"
)

type TestConfig struct {
	ServerAddress           string        `env:"TEST_SERVER_ADDRESS"`
	PostgreSqlDSN           string        `env:"TEST_POSTGRES_DSN"`
	JWTPrivateKey           string        `env:"TEST_JWT_PRIVATE_KEY"`
	RedisURl                string        `env:"TEST_REDIS_URL"`
	RedisExpirationDuration time.Duration `env:"TEST_REDIS_EXPIRATION_DURATION"`
	LogLevel                string        `env:"TEST_LOG_LEVEL"`
}

type BannerSuite struct {
	suite.Suite
	jwtManager *utils.JWTManager
	database   *pg.PGStorage
	cache      *redis.RedisManager
	router     *HttpHandler
}

func (s *BannerSuite) SetupSuite() {
	ctx := context.Background()

	//cfg := TestConfig{}
	//err := env.Parse(&cfg)
	//s.Nil(err)

	cfg, err := config.GetConfig()
	s.Nil(err)

	pg, err := pg.NewPG(ctx, cfg.PostgreSqlDSN)
	s.Nil(err)

	s.database = pg

	redisManager, err := redis.NewRedisManager(cfg.RedisURl, cfg.RedisExpirationDuration)
	s.Nil(err)

	s.cache = redisManager

	jwtManager, err := utils.NewJWTManager(cfg.JWTPrivateKey)
	s.Nil(err)

	s.jwtManager = jwtManager

	ctrl, err := controller.NewController(pg, redisManager)
	s.Nil(err)

	h, err := NewHttpHandler(ctrl, jwtManager)

	s.router = h
}

func (s *BannerSuite) TearDownSuite() {
	ctx := context.Background()

	err := s.database.ReleaseBanners(ctx)
	s.Nil(err)

	s.cache.Clear(ctx)
}

const (
	AuthorizationHeader  = "Authorization"
	ContentTypeHeader    = "Content-type"
	JSONContentType      = "application/json"
	FeatureIdParam       = "feature_id"
	TagIdParam           = "tag_id"
	UseLastRevisionParam = "use_last_revision"
)

func (s *BannerSuite) TestGetUserBanner() {
	url := "/user_banner"

	type BannerTestCase struct {
		role           int
		input          models.GetBannerInput
		expectedBanner models.Banner
		expectedStatus int
	}

	activeBanner := BannerTestCase{
		role: USER,
		input: models.GetBannerInput{
			TagId:           1,
			FeatureId:       1,
			UseLastRevision: false,
		},
		expectedBanner: models.Banner{
			FeatureId: 1,
			TagIds:    []int{1, 2},
			Content:   `{"title": "Active Banner"}`,
			IsActive:  true,
		},
		expectedStatus: http.StatusOK,
	}

	inactiveBanner := BannerTestCase{
		role: USER,
		input: models.GetBannerInput{
			TagId:     3,
			FeatureId: 2,
		},
		expectedBanner: models.Banner{
			FeatureId: 2,
			TagIds:    []int{2, 3},
			Content:   `{"title": "Inactive banner"}`,
			IsActive:  false,
		},
		expectedStatus: http.StatusOK,
	}

	activeCashedBanner := BannerTestCase{
		role: USER,
		input: models.GetBannerInput{
			TagId:     3,
			FeatureId: 4,
		},
		expectedBanner: models.Banner{
			FeatureId: 4,
			TagIds:    []int{3, 4},
			Content:   `{"title": "Cached active banner"}`,
			IsActive:  true,
		},
		expectedStatus: http.StatusOK,
	}

	notFoundBanner := BannerTestCase{
		role: USER,
		input: models.GetBannerInput{
			TagId:           4,
			FeatureId:       5,
			UseLastRevision: false,
		},
		expectedBanner: models.Banner{
			FeatureId: 5,
			TagIds:    []int{4, 5},
			Content:   `{"title": "Cached active banner"}`,
			IsActive:  true,
		},
		expectedStatus: http.StatusOK,
	}

	s.Run("Успешное получение активного баннера для пользователя 200 OK", func() {
		testCase := activeBanner
		testCase.role = USER

		banner, err := s.database.CreateBanner(context.Background(), &testCase.expectedBanner)
		defer func() {
			err = s.database.DeleteBanner(context.Background(), banner.Id)
			s.Nil(err)
		}()
		s.Nil(err)

		apitest.
			New().
			Handler(s.router).
			Get(url).
			Header(AuthorizationHeader, s.generateBearerToken(testCase.role)).
			Query(FeatureIdParam, fmt.Sprintf("%d", testCase.input.FeatureId)).
			Query(TagIdParam, fmt.Sprintf("%d", testCase.input.TagId)).
			Expect(s.T()).
			Header(ContentTypeHeader, JSONContentType).
			Status(testCase.expectedStatus).
			Body(testCase.expectedBanner.Content).
			End()

	})

	s.Run("Успешное получение активного баннера для админа 200 OK", func() {
		testCase := activeBanner
		testCase.role = ADMIN

		banner, err := s.database.CreateBanner(context.Background(), &testCase.expectedBanner)
		defer func() {
			err = s.database.DeleteBanner(context.Background(), banner.Id)
			s.Nil(err)
		}()
		s.Nil(err)

		apitest.
			New().
			Handler(s.router).
			Get(url).
			Header(AuthorizationHeader, s.generateBearerToken(testCase.role)).
			Query(FeatureIdParam, fmt.Sprintf("%d", testCase.input.FeatureId)).
			Query(TagIdParam, fmt.Sprintf("%d", testCase.input.TagId)).
			Expect(s.T()).
			Header(ContentTypeHeader, JSONContentType).
			Status(testCase.expectedStatus).
			Body(testCase.expectedBanner.Content).
			End()
	})

	s.Run("Успешное получение неактивного баннера для админа 200 OK", func() {
		testCase := inactiveBanner
		testCase.role = ADMIN

		banner, err := s.database.CreateBanner(context.Background(), &testCase.expectedBanner)
		defer func() {
			err = s.database.DeleteBanner(context.Background(), banner.Id)
			s.Nil(err)
		}()
		s.Nil(err)

		apitest.
			New().
			Handler(s.router).
			Get(url).
			Header(AuthorizationHeader, s.generateBearerToken(testCase.role)).
			Query(FeatureIdParam, fmt.Sprintf("%d", testCase.input.FeatureId)).
			Query(TagIdParam, fmt.Sprintf("%d", testCase.input.TagId)).
			Expect(s.T()).
			Header(ContentTypeHeader, JSONContentType).
			Status(testCase.expectedStatus).
			Body(testCase.expectedBanner.Content).
			End()
	})

	s.Run("Успешное получение баннера напрямую из базы 200 OK", func() {
		testCase := activeBanner
		testCase.input.UseLastRevision = true

		banner, err := s.database.CreateBanner(context.Background(), &testCase.expectedBanner)
		defer func() {
			err = s.database.DeleteBanner(context.Background(), banner.Id)
			s.Nil(err)
		}()
		s.Nil(err)

		apitest.
			New().
			Handler(s.router).
			Get(url).
			Header(AuthorizationHeader, s.generateBearerToken(testCase.role)).
			Query(FeatureIdParam, fmt.Sprintf("%d", testCase.input.FeatureId)).
			Query(TagIdParam, fmt.Sprintf("%d", testCase.input.TagId)).
			Query(UseLastRevisionParam, fmt.Sprintf("%t", testCase.input.UseLastRevision)).
			Expect(s.T()).
			Header(ContentTypeHeader, JSONContentType).
			Status(testCase.expectedStatus).
			Body(testCase.expectedBanner.Content).
			End()
	})

	s.Run("Успешное получение активного баннера из кэша 200 OK", func() {
		testCase := activeCashedBanner

		err := s.cache.SetBanner(
			context.Background(),
			testCase.input.FeatureId,
			testCase.input.TagId,
			&testCase.expectedBanner.Content)
		s.Nil(err)

		apitest.
			New().
			Handler(s.router).
			Get(url).
			Header(AuthorizationHeader, s.generateBearerToken(testCase.role)).
			Query(FeatureIdParam, fmt.Sprintf("%d", testCase.input.FeatureId)).
			Query(TagIdParam, fmt.Sprintf("%d", testCase.input.TagId)).
			Expect(s.T()).
			Header(ContentTypeHeader, JSONContentType).
			Status(testCase.expectedStatus).
			Body(testCase.expectedBanner.Content).
			End()
	})

	s.Run("Неудачное получение баннера - нет параметра в запросе 400 Bad Request", func() {
		testCase := activeBanner
		testCase.expectedStatus = http.StatusBadRequest

		banner, err := s.database.CreateBanner(context.Background(), &testCase.expectedBanner)
		defer func() {
			err = s.database.DeleteBanner(context.Background(), banner.Id)
			s.Nil(err)
		}()
		s.Nil(err)

		apitest.
			New().
			Handler(s.router).
			Get(url).
			Header(AuthorizationHeader, s.generateBearerToken(testCase.role)).
			Query(TagIdParam, fmt.Sprintf("%d", testCase.input.TagId)).
			Expect(s.T()).
			Header(ContentTypeHeader, JSONContentType).
			Status(testCase.expectedStatus).
			End()
	})

	s.Run("Неудачное получение баннера - баннер не найден 404 Not found", func() {
		testCase := notFoundBanner
		testCase.expectedStatus = http.StatusNotFound

		apitest.
			New().
			Handler(s.router).
			Get(url).
			Header(AuthorizationHeader, s.generateBearerToken(testCase.role)).
			Query(FeatureIdParam, fmt.Sprintf("%d", testCase.input.FeatureId)).
			Query(TagIdParam, fmt.Sprintf("%d", testCase.input.TagId)).
			Expect(s.T()).
			Header(ContentTypeHeader, JSONContentType).
			Status(testCase.expectedStatus).
			End()
	})

	s.Run("Неудачное получение баннера пользователем - баннер не активен 404 Not found", func() {
		testCase := inactiveBanner
		testCase.expectedStatus = http.StatusNotFound

		banner, err := s.database.CreateBanner(context.Background(), &testCase.expectedBanner)
		defer func() {
			err = s.database.DeleteBanner(context.Background(), banner.Id)
			s.Nil(err)
		}()
		s.Nil(err)

		apitest.
			New().
			Handler(s.router).
			Get(url).
			Header(AuthorizationHeader, s.generateBearerToken(testCase.role)).
			Query(FeatureIdParam, fmt.Sprintf("%d", testCase.input.FeatureId)).
			Query(TagIdParam, fmt.Sprintf("%d", testCase.input.TagId)).
			Expect(s.T()).
			Header(ContentTypeHeader, JSONContentType).
			Status(testCase.expectedStatus).
			End()
	})

	s.Run("Неудачное получение баннера пользователем - невалидный токен 401 Unauthorized", func() {
		testCase := activeBanner
		testCase.expectedStatus = http.StatusUnauthorized

		banner, err := s.database.CreateBanner(context.Background(), &testCase.expectedBanner)
		defer func() {
			err = s.database.DeleteBanner(context.Background(), banner.Id)
			s.Nil(err)
		}()
		s.Nil(err)

		apitest.
			New().
			Handler(s.router).
			Get(url).
			Header(AuthorizationHeader, "invalid-token").
			Query(FeatureIdParam, fmt.Sprintf("%d", testCase.input.FeatureId)).
			Query(TagIdParam, fmt.Sprintf("%d", testCase.input.TagId)).
			Expect(s.T()).
			Header(ContentTypeHeader, JSONContentType).
			Status(testCase.expectedStatus).
			End()
	})

	s.Run("Неудачное получение баннера пользователем - несуществующий тип пользователя 403 Forbidden", func() {
		testCase := activeBanner
		testCase.role = 42
		testCase.expectedStatus = http.StatusForbidden

		banner, err := s.database.CreateBanner(context.Background(), &testCase.expectedBanner)
		defer func() {
			err = s.database.DeleteBanner(context.Background(), banner.Id)
			s.Nil(err)
		}()
		s.Nil(err)

		apitest.
			New().
			Handler(s.router).
			Get(url).
			Header(AuthorizationHeader, s.generateBearerToken(testCase.role)).
			Query(FeatureIdParam, fmt.Sprintf("%d", testCase.input.FeatureId)).
			Query(TagIdParam, fmt.Sprintf("%d", testCase.input.TagId)).
			Expect(s.T()).
			Header(ContentTypeHeader, JSONContentType).
			Status(testCase.expectedStatus).
			End()
	})

}

func TestSuite(t *testing.T) {
	suite.Run(t, new(BannerSuite))
}

func (s *BannerSuite) generateBearerToken(role int) string {
	token, err := s.jwtManager.Generate(role)
	s.Nil(err)
	return "Bearer " + token
}
