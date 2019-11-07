package store

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"time"
)

type Store struct {
	gorm *gorm.DB
	redisClient *redis.Client
}

func NewStore(db *gorm.DB, redisClient *redis.Client) *Store {
	return &Store{gorm: db, redisClient: redisClient}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func (s *Store)Set(key string, value string, expiration time.Duration){
	s.redisClient.Set(key, value, expiration)
}

func (s *Store)Get(key string) string{
	return s.redisClient.Get(key).Val()
}