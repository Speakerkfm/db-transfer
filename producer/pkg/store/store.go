package store

import (
	"database/sql"
	"github.com/go-redis/redis"
	"time"
)

type Store struct {
	db *sql.DB
	redisClient *redis.Client
}

func NewStore(db *sql.DB, redisClient *redis.Client) *Store {
	return &Store{db: db, redisClient: redisClient}
}

func (s *Store)Set(key string, value string, expiration time.Duration){
	s.redisClient.Set(key, value, expiration)
}

func (s *Store)Get(key string) string{
	return s.redisClient.Get(key).Val()
}