package app

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	iamAPI "github.com/ChopX4/raketka/iam/internal/api"
	"github.com/ChopX4/raketka/iam/internal/config"
	iamMigrator "github.com/ChopX4/raketka/iam/internal/migrator"
	"github.com/ChopX4/raketka/iam/internal/repository"
	sessionRepository "github.com/ChopX4/raketka/iam/internal/repository/session"
	userRepository "github.com/ChopX4/raketka/iam/internal/repository/user"
	"github.com/ChopX4/raketka/iam/internal/service"
	iamService "github.com/ChopX4/raketka/iam/internal/service/iam"
	"github.com/ChopX4/raketka/platform/pkg/cache"
	platformRedis "github.com/ChopX4/raketka/platform/pkg/cache/redis"
	"github.com/ChopX4/raketka/platform/pkg/closer"
	"github.com/ChopX4/raketka/platform/pkg/logger"
	"github.com/ChopX4/raketka/platform/pkg/pgxtx"
	auth_v1 "github.com/ChopX4/raketka/shared/pkg/proto/auth/v1"
)

type diContainer struct {
	authV1API auth_v1.AuthServiceServer

	iamService service.IamService

	sessionRepository repository.SessionRepository
	userRepository    repository.UserRepository
	txManager         pgxtx.TxManager

	postgresPool *pgxpool.Pool
	redisPool    *redigo.Pool
	redisClient  cache.RedisClient
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) AuthV1API(ctx context.Context) auth_v1.AuthServiceServer {
	if d.authV1API == nil {
		d.authV1API = iamAPI.NewIamApi(d.IamService(ctx))
	}

	return d.authV1API
}

func (d *diContainer) IamService(ctx context.Context) service.IamService {
	if d.iamService == nil {
		d.iamService = iamService.NewIamService(
			d.SessionRepository(ctx),
			d.UserRepository(ctx),
			d.TxManager(ctx),
			config.AppConfig().Session.TTL(),
		)
	}

	return d.iamService
}

func (d *diContainer) SessionRepository(ctx context.Context) repository.SessionRepository {
	if d.sessionRepository == nil {
		d.sessionRepository = sessionRepository.NewSessionRepository(d.RedisClient(ctx))
	}

	return d.sessionRepository
}

func (d *diContainer) UserRepository(ctx context.Context) repository.UserRepository {
	if d.userRepository == nil {
		d.userRepository = userRepository.NewUserRepository(d.PostgreSQLPool(ctx))
	}

	return d.userRepository
}

func (d *diContainer) TxManager(ctx context.Context) pgxtx.TxManager {
	if d.txManager == nil {
		d.txManager = pgxtx.NewTxManager(d.PostgreSQLPool(ctx))
	}

	return d.txManager
}

func (d *diContainer) PostgreSQLPool(ctx context.Context) *pgxpool.Pool {
	if d.postgresPool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().Postgre.URI())
		if err != nil {
			logger.Error(ctx, "failed to create postgres pool", zap.Error(err))
			panic(fmt.Sprintf("failed to create pgx pool: %v", err))
		}

		if err = pool.Ping(ctx); err != nil {
			logger.Error(ctx, "failed to ping postgres", zap.Error(err))
			pool.Close()
			panic(fmt.Sprintf("failed to ping postgres: %v", err))
		}

		closer.AddNamed("postgres pool", func(context.Context) error {
			pool.Close()
			return nil
		})

		d.postgresPool = pool
	}

	return d.postgresPool
}

func (d *diContainer) RedisPool(ctx context.Context) *redigo.Pool {
	if d.redisPool == nil {
		pool := &redigo.Pool{
			MaxIdle:     config.AppConfig().Redis.MaxIdle(),
			IdleTimeout: config.AppConfig().Redis.IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", config.AppConfig().Redis.Address())
			},
			TestOnBorrowContext: func(ctx context.Context, conn redigo.Conn, _ time.Time) error {
				_, err := conn.Do("PING")
				return err
			},
		}

		conn, err := pool.GetContext(ctx)
		if err != nil {
			logger.Error(ctx, "failed to get redis connection", zap.Error(err))
			panic(fmt.Sprintf("failed to create redis pool: %v", err))
		}

		if _, err = conn.Do("PING"); err != nil {
			logger.Error(ctx, "failed to ping redis", zap.Error(err))

			if closeErr := conn.Close(); closeErr != nil {
				logger.Error(ctx, "failed to close redis connection after ping error", zap.Error(closeErr))
			}

			if closeErr := pool.Close(); closeErr != nil {
				logger.Error(ctx, "failed to close redis pool after ping error", zap.Error(closeErr))
			}

			panic(fmt.Sprintf("failed to ping redis: %v", err))
		}

		if err = conn.Close(); err != nil {
			logger.Error(ctx, "failed to close redis connection", zap.Error(err))
			panic(fmt.Sprintf("failed to close redis connection: %v", err))
		}

		closer.AddNamed("redis pool", func(context.Context) error {
			return pool.Close()
		})

		d.redisPool = pool
	}

	return d.redisPool
}

func (d *diContainer) RedisClient(ctx context.Context) cache.RedisClient {
	if d.redisClient == nil {
		client := platformRedis.NewClient(
			d.RedisPool(ctx),
			logger.Logger(),
			config.AppConfig().Redis.ConnectionTimeout(),
		)

		if err := client.Ping(ctx); err != nil {
			logger.Error(ctx, "failed to ping redis via client", zap.Error(err))
			panic(fmt.Sprintf("failed to ping redis via client: %v", err))
		}

		d.redisClient = client
	}

	return d.redisClient
}

func (d *diContainer) RunMigrations(ctx context.Context) {
	db, err := sql.Open("pgx", config.AppConfig().Postgre.URI())
	if err != nil {
		logger.Error(ctx, "failed to connect to database for migrations", zap.Error(err))
		panic(fmt.Sprintf("failed to connect to database for migrations: %v", err))
	}

	migratorRunner := iamMigrator.NewMigrator(db, config.AppConfig().Postgre.MigrationsPath())
	if err = migratorRunner.Up(); err != nil {
		logger.Error(ctx, "failed to run database migrations", zap.Error(err))
		if closeErr := db.Close(); closeErr != nil {
			logger.Error(ctx, "failed to close migrations database connection after migration error", zap.Error(closeErr))
		}
		panic(fmt.Sprintf("failed to migrate database: %v", err))
	}

	if err = db.Close(); err != nil {
		logger.Error(ctx, "failed to close migrations database connection", zap.Error(err))
		panic(fmt.Sprintf("failed to close migrations database connection: %v", err))
	}
}
