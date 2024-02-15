package cache

import (
	"context"
	"time"
)

type Config struct {
	Storage           string
	DefaultExpiration time.Duration
	CleanupInterval   time.Duration
}

type Service interface {
	SetGlobalTTL(ttl time.Duration) error
	Set(primaryKey, key string, value interface{}, duration time.Duration) error
	GetKeysByPrimaryKey(primaryKey string) []string
	GetPrimaryKeys() []string
	Get(primaryKey, key string) (interface{}, error)
	GetAll() []interface{}
	GetAllByPatternByPrimaryKey(primaryKey, key string) ([]interface{}, error)
	GetItems(primaryKey, key string) *Item
	GetAllByPrimaryKey(primaryKey string) []interface{}
	GetByKey(key string) interface{}
	// Count() int
}

func NewService(ctx context.Context, cfg *Config) Service {
	return NewServiceInMemory(ctx, cfg)
}
