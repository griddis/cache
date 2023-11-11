package cache

import (
	"context"
	"regexp"
	"sync"
	"time"

	"github.com/pkg/errors"

	logging "github.com/griddis/go-logger"
)

type inmemory struct {
	config            *Config
	logger            logging.Logger
	s                 sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	items             map[string]map[string]*Item
}

type Keys struct {
	NameSpace string
	Key       string
}

type Item struct {
	Value      interface{}
	Created    time.Time
	LastUpdate time.Time
	Expiration int64
}

func NewServiceInMemory(ctx context.Context, cfg *Config) Service {
	logger := logging.FromContext(ctx)
	logger = logger.With("svc", "cache")
	logger = logger.SetDefaultFieldName("msg")

	items := map[string]map[string]*Item{}
	cache := inmemory{
		config:            cfg,
		logger:            logger,
		items:             items,
		defaultExpiration: cfg.DefaultExpiration,
		cleanupInterval:   cfg.CleanupInterval,
	}

	if cfg.CleanupInterval > 0 {
		cache.StartGC()
	}

	return &cache
}

//1189 ns
func (s *inmemory) Get(nameSpace, key string) (interface{}, error) {
	s.s.RLock()
	item, found := s.items[nameSpace][key]
	s.s.RUnlock()

	if !found {
		return nil, errors.New("cache not found")
	}

	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return nil, errors.New("cache is expired, but is stored")
		}
	}
	return item.Value, nil
}

//287M ns
func (s *inmemory) GetKeysByNamespace(nameSpace string) (keys []string) {
	//reflect.ValueOf(s.items).MapKeys()
	s.s.RLock()
	for key := range s.items[nameSpace] {
		keys = append(keys, key)
	}
	s.s.RUnlock()
	return keys
}

//680 ns
func (s *inmemory) GetNamespaces() (nameSpace []string) {
	nameSpace = make([]string, 0, len(s.items))
	s.s.RLock()
	for ns := range s.items {
		nameSpace = append(nameSpace, ns)
	}
	s.s.RUnlock()
	return nameSpace
}

func (s *inmemory) GetItems(nameSpace, key string) *Item {
	s.s.RLock()
	item := s.items[nameSpace][key]
	s.s.RUnlock()
	return item
}

//188M ns
func (s *inmemory) GetAll() []interface{} {
	var ret []interface{}
	s.s.RLock()
	for _, item := range s.items {
		for _, it := range item {
			if it.Expiration > 0 {
				if time.Now().UnixNano() > it.Expiration {
					continue
				}
			}
			ret = append(ret, it.Value)
		}
	}
	s.s.RUnlock()
	return ret
}

//779 ns
func (s *inmemory) GetAllByPatternByNamespace(nameSpace, key string) ([]interface{}, error) { //namespase
	r, err := regexp.Compile(key)
	if err != nil {
		return nil, errors.Wrap(err, "GetAllByPattern error compile regexp")
	}
	var ret []interface{}
	s.s.RLock()
	if item, ok := s.items[nameSpace]; ok {
		for k, it := range item {
			if !r.MatchString(k) {
				continue
			}
			if it.Expiration > 0 {
				if time.Now().UnixNano() > it.Expiration {
					continue
				}
			}
			ret = append(ret, it.Value)
		}
	}
	s.s.RUnlock()
	return ret, nil
}

//23.29 ns
func (s *inmemory) GetAllByNamespace(nameSpace string) []interface{} {
	var ret []interface{}
	s.s.RLock()
	if item, ok := s.items[nameSpace]; ok {
		for _, it := range item {
			if it.Expiration > 0 {
				if time.Now().UnixNano() > it.Expiration {
					continue
				}
			}
			ret = append(ret, it.Value)
		}
	}
	s.s.RUnlock()
	return ret
}

//23.68 ns
func (s *inmemory) GetByKey(key string) interface{} {
	var ret interface{}
	s.s.RLock()
	for _, item := range s.items {
		if it, ok := item[key]; ok {
			if it.Expiration > 0 {
				if time.Now().UnixNano() > it.Expiration {
					continue
				}
				ret = it.Value
			}
		}
	}
	s.s.RUnlock()
	return ret
}

//223 ns
func (s *inmemory) Set(nameSpace, key string, value interface{}, duration time.Duration) error {
	var expiration int64
	if duration == 0 {
		duration = s.defaultExpiration
	}

	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}
	s.s.Lock()
	item, ok := s.items[nameSpace]
	if !ok {
		item = map[string]*Item{}
	}
	item[key] = &Item{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}
	s.items[nameSpace] = item
	s.s.Unlock()

	return nil
}

func (s *inmemory) SetGlobalTTL(ttl time.Duration) error {
	return nil
}

func (s *inmemory) StartGC() {
	go s.gc()
}

func (s *inmemory) gc() {
	for {
		<-time.After(s.cleanupInterval)
		s.logger.Debug("start gc")

		if s.items == nil {
			s.logger.Debug("gc completed, items is zero record")
			return
		}

		if keys := s.expiredKeys(); len(keys) != 0 {
			s.clearItems(keys)
		}

	}
}

func (s *inmemory) expiredKeys() (keys []Keys) {
	var it *Item
	var item map[string]*Item
	var key Keys
	s.s.RLock()
	for key.NameSpace, item = range s.items {
		for key.Key, it = range item {
			if time.Now().UnixNano() > it.Expiration && it.Expiration > 0 {
				keys = append(keys, key)
			}
		}
	}
	s.s.RUnlock()

	return
}

func (s *inmemory) clearItems(keys []Keys) {
	s.s.Lock()
	for _, k := range keys {
		s.logger.Debug("clean cache", "nameSpace", k.NameSpace, "key", k.Key)
		delete(s.items[k.NameSpace], k.Key)
	}
	s.s.Unlock()
}
