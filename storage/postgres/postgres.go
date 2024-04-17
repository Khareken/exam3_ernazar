package postgres

import (
	"context"
	"exam/config"
	"exam/pkg/logger"
	"exam/storage"
	"exam/storage/redis"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)


type Store struct {
	Pool *pgxpool.Pool
	logger logger.ILogger
	cfg config.Config
	redis  storage.IRedisStorage
}


func New(ctx context.Context,cfg config.Config,log logger.ILogger,redis storage.IRedisStorage) (storage.IStorage,error) {
	url := fmt.Sprintf(`host=%s port=%v user=%s password=%s database=%s sslmode=disable`,
	cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDatabase)

	pgPoolConfig,err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil,err
	}

	pgPoolConfig.MaxConns = 100
	pgPoolConfig.MaxConnLifetime = time.Hour

	ctx,cancel:= context.WithTimeout(ctx,config.TimewithContex)
	defer cancel()

	newPool,err := pgxpool.NewWithConfig(ctx,pgPoolConfig)
	if err != nil {
       fmt.Println("error while connecting to db",err.Error())
	   return nil,err
	}
	return Store{
		Pool: newPool,
		cfg: cfg,
		logger: log,
		redis: redis,
	},nil
}

func (s Store) CloseDB() {
	s.Pool.Close()
}

func (s Store) Customer() storage.ICustomerStorage {
	newCustomer := NewCustomer(s.Pool,s.logger)
  
	return &newCustomer
}

func (s Store) Redis() storage.IRedisStorage {
	newRedis := redis.New(s.cfg)

	return newRedis
}