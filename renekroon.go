package cache

/*
import (
	"context"
	"fmt"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	logging "github.com/griddis/go-logger"
)

type renekroon struct {
	config *Config
	ttl    time.Duration
	cache  *ttlcache.Cache
	logger logging.Logger
}

func NewServiceReneKroon(ctx context.Context, cfg *Config) Service {
	logger := logging.FromContext(ctx)
	logger = logger.With("svc", "netbox")

	var cache *ttlcache.Cache = ttlcache.NewCache()
	checkExpirationCallback := func(key string, value interface{}) bool {
		logger.Debug("checkExpireCallback", "key", key)
		// all other values are allowed to expire
		return true
	}

	expirationCallback := func(key string, reason ttlcache.EvictionReason, value interface{}) {
		fmt.Printf("This key(%s) has expired because of %s\n", key, reason)
	}
	expireCallback := func(key string, value interface{}) {
		logger.Debug("expireCallback", "key", key)
	}
	cache.SetCheckExpirationCallback(checkExpirationCallback)
	cache.SetExpirationReasonCallback(expirationCallback)
	cache.SetExpirationCallback(expireCallback)
	return &renekroon{
		config: cfg,
		cache:  cache,
		logger: logger,
	}
}

func (s *renekroon) Get(key string) (interface{}, error) {
	return s.cache.Get(key)
}

func (s *renekroon) GetKeys() []string {
	return s.cache.GetKeys()
}

func (s *renekroon) GetAll() []interface{} {
	var ret []interface{}
	return ret
}

func (s *renekroon) GetAllItems() []*Item {
	var ret []*Item
	return ret
}

func (s *renekroon) GetAllByPattern(search string) ([]interface{}, error) {
	var ret []interface{}
	return ret, nil
}

func (s *renekroon) Set(key string, value interface{}, duration time.Duration) error {
	//s.cache.SetWithTTL(key, value, s.ttl)
	s.cache.Set(key, value)
	return nil
}

func (s *renekroon) SetGlobalTTL(ttl time.Duration) error {
	s.ttl = ttl
	s.cache.SetTTL(ttl)
	return nil
}
*/
