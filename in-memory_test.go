package cache

import (
	"reflect"
	"sync"
	"testing"
	"time"

	logging "github.com/griddis/go-logger"
)

func Test_inmemory_GetAllByPattern(t *testing.T) {

	type fields struct {
		config            *Config
		logger            logging.Logger
		s                 sync.RWMutex
		defaultExpiration time.Duration
		cleanupInterval   time.Duration
		items             map[string]map[string]*Item
	}

	tests := []struct {
		name    string
		fields  fields
		args    Keys
		want    []interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &inmemory{
				config:            tt.fields.config,
				logger:            tt.fields.logger,
				s:                 tt.fields.s,
				defaultExpiration: tt.fields.defaultExpiration,
				cleanupInterval:   tt.fields.cleanupInterval,
				items:             tt.fields.items,
			}
			got, err := s.GetAllByPatternByNamespace(tt.args.NameSpace, tt.args.Key)
			if (err != nil) != tt.wantErr {
				t.Errorf("inmemory.GetAllByPattern() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("inmemory.GetAllByPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}
