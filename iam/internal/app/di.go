package app

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	redigo "github.com/gomodule/redigo/redis"
	"go.uber.org/zap"

	authv1 "iam/internal/api/auth/v1"
	userv1 "iam/internal/api/user/v1"
	"iam/internal/config"
	userrepo "iam/internal/repository/user"
	"iam/internal/repository/session"
	"iam/internal/service/auth"
	userservice "iam/internal/service/user"
	"platform/pkg/cache/redis"
	"platform/pkg/logger"
	"platform/pkg/migrator/pg"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type DI struct {
	Config      *config.Config
	DB          *sql.DB
	UserRepo    userrepo.Repository
	SessionRepo session.Repository
	UserService userservice.Service
	AuthService auth.Service
	UserAPI     *userv1.API
	AuthAPI     *authv1.API
}

func NewDI(cfg *config.Config) (*DI, error) {
	logger.Init(cfg.Logger.Level(), cfg.Logger.AsJSON())
	log := logger.Logger()
	ctx := context.Background()

	di := &DI{Config: cfg}

	// 1. PostgreSQL
	db, err := sql.Open("pgx", cfg.Postgres.DSN())
	if err != nil {
		return nil, fmt.Errorf("db open: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}
	di.DB = db
	log.Info(ctx, "Connected to PostgreSQL", zap.String("db", cfg.Postgres.Database()))

	// 2. Миграции через platform
	m := pg.New(db, migrationsFS, "migrations")
	if err := m.Up(ctx); err != nil {
		return nil, fmt.Errorf("migrations up: %w", err)
	}
	log.Info(ctx, "Migrations applied successfully")

	// 3. Redis
	redisPool := &redigo.Pool{
		MaxIdle:     cfg.Redis.MaxIdle(),
		IdleTimeout: cfg.Redis.IdleTimeout(),
		Dial: func() (redigo.Conn, error) {
			return redigo.Dial("tcp", cfg.Redis.Address())
		},
	}
	redisClient := redis.NewClient(redisPool, log, cfg.Redis.ConnectionTimeout())
	log.Info(ctx, "Connected to Redis", zap.String("addr", cfg.Redis.Address()))

	// 4. Репозитории
	di.UserRepo = userrepo.NewPostgresRepository(db)
	di.SessionRepo = session.NewRedisRepository(redisClient, cfg.Session.TTL())

	// 5. Сервисы
	di.UserService = userservice.NewService(di.UserRepo)
	di.AuthService = auth.NewService(di.UserRepo, di.SessionRepo, cfg.Session.TTL())

	// 6. API
	di.UserAPI = userv1.NewAPI(di.UserService)
	di.AuthAPI = authv1.NewAPI(di.AuthService)

	log.Info(ctx, "IAMService DI initialized")
	return di, nil
}