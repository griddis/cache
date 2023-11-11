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
	Set(nameSpace, key string, value interface{}, duration time.Duration) error
	GetKeysByNamespace(nameSpace string) []string
	GetNamespaces() []string
	Get(nameSpace, key string) (interface{}, error)
	GetAll() []interface{}
	GetAllByPatternByNamespace(nameSpace, key string) ([]interface{}, error)
	GetItems(nameSpace, key string) *Item
	GetAllByNamespace(nameSpace string) []interface{}
	GetByKey(key string) interface{}
	// Count() int
}

func NewService(ctx context.Context, cfg *Config) Service {
	return NewServiceInMemory(ctx, cfg)
}
