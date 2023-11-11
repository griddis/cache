package cache

import (
	"fmt"
	"sync"
	"testing"
	"time"

	logging "github.com/griddis/go-logger"
)

func add(m map[string]map[string]*Item, key1, key2 string) {
	mm, ok := m[key1]
	if !ok {
		mm = make(map[string]*Item)
		m[key1] = mm
	}
	mm[key2] = &Item{
		Value:   "qwerty",
		Created: time.Now(),
	}
}

func Benchmark(b *testing.B) {
	type fields struct {
		config            *Config
		logger            logging.Logger
		s                 sync.RWMutex
		defaultExpiration time.Duration
		cleanupInterval   time.Duration
		items             map[string]map[string]*Item
	}
	var keys Keys
	var field fields
	field.items = make(map[string]map[string]*Item)
	for q := 0; q < 30; q++ {
		for i := 0; i < 100000; i++ {
			keys.NameSpace = fmt.Sprint(q)
			keys.Key = fmt.Sprint(i)
			add(field.items, keys.NameSpace, keys.Key)
		}
	}
	tests := struct {
		name    string
		fields  fields
		args    Keys
		want    []interface{}
		wantErr bool
	}{
		fields: field,
	}
	s := &inmemory{
		config:            tests.fields.config,
		logger:            tests.fields.logger,
		s:                 tests.fields.s,
		defaultExpiration: tests.fields.defaultExpiration,
		cleanupInterval:   tests.fields.cleanupInterval,
		items:             tests.fields.items,
	}

	b.Run("get all by pettern", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s.GetAllByPatternByNamespace(tests.args.NameSpace, tests.args.Key)
		}
	})

	b.Run("get all", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s.GetAll()
		}
	})

	b.Run("get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s.Get(tests.args.NameSpace, tests.args.Key)
		}
	})

	b.Run("get keys", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s.GetKeysByNamespace(tests.args.NameSpace)
		}
	})

	b.Run("get nameSpace", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s.GetNamespaces()
		}
	})

	b.Run("get by nameSpace", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s.GetAllByNamespace(tests.args.NameSpace)
		}
	})

	b.Run("get by key", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s.GetByKey(tests.args.Key)
		}
	})

	b.Run("set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s.Set(tests.args.NameSpace, tests.args.Key, tests.fields.items, 300*time.Second)
		}
	})
}
