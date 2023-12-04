package cache

import (
	"reflect"
	"testing"
	"time"

	"golang.org/x/exp/slices"
)

func TestNew(t *testing.T) {
	type args struct {
		defaultExpiry   time.Duration
		cleanupInterval time.Duration
	}
	tests := []struct {
		name string
		args args
		want *Cache
	}{
		{
			name: "default",
			args: args{0, 0},
			want: &Cache{
				defaultExpiry:   0,
				cleanupInterval: 0,
				items:           make(map[string]Item),
			},
		},
		{
			name: "with expiry and cleanup",
			args: args{1, 1},
			want: &Cache{
				defaultExpiry:   1,
				cleanupInterval: 1,
				items:           make(map[string]Item),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.defaultExpiry, tt.args.cleanupInterval); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Set(t *testing.T) {
	type fields struct {
		defaultExpiry   time.Duration
		cleanupInterval time.Duration
		items           map[string]Item
	}
	type args struct {
		key      string
		value    any
		duration time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "first",
			fields: fields{0, 0, make(map[string]Item)},
			args: args{
				key:      "item",
				value:    "item",
				duration: 0,
			},
			want: "item",
		},
		{
			name:   "second",
			fields: fields{0, 0, make(map[string]Item)},
			args: args{
				key:      "item2",
				value:    "item2",
				duration: 0,
			},
			want: "item2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				defaultExpiry:   tt.fields.defaultExpiry,
				cleanupInterval: tt.fields.cleanupInterval,
				items:           tt.fields.items,
			}
			c.Set(tt.args.key, tt.args.value, tt.args.duration)
			got, found := c.Get(tt.args.key)
			if !found || got != tt.want {
				t.Errorf("Got %s want %s, found %v", got, tt.want, found)
			}
		})
	}
}

func TestCache_Get(t *testing.T) {
	type args struct {
		key    string
		value  string
		search string
	}
	tests := []struct {
		name  string
		args  args
		want  any
		want1 bool
	}{
		{
			name: "first",
			args: args{
				key:    "item",
				value:  "item",
				search: "item",
			},
			want:  "item",
			want1: true,
		},
		{
			name: "second",
			args: args{
				key:    "item",
				value:  "item",
				search: "item1",
			},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(0, 0)
			c.Set(tt.args.key, tt.args.value, 0)
			got, got1 := c.Get(tt.args.search)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cache.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Cache.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestCache_Delete(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "first",
			args: args{
				key:   "item",
				value: "item",
			},
			wantErr: true,
		},
		{
			name: "second",
			args: args{
				key:   "item",
				value: "item",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(0, 0)
			if !tt.wantErr {
				c.Set(tt.args.key, tt.args.value, 0)
			}
			if err := c.Delete(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Cache.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCache_expiredKeys(t *testing.T) {
	tests := []struct {
		name   string
		expiry time.Duration
		want   []string
		keys   []string
	}{
		{
			name:   "first",
			expiry: time.Nanosecond,
			want:   []string{"item", "item1", "item2"},
			keys:   []string{"item", "item1", "item2"},
		},
		{
			name:   "second",
			expiry: time.Hour,
			want:   []string{},
			keys:   []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.expiry, 0)
			for _, key := range tt.keys {
				c.Set(key, "", tt.expiry)
			}
			time.Sleep(10 * time.Nanosecond)
			got := c.expiredKeys()
			slices.Sort(got)
			if !reflect.DeepEqual(got, tt.want) && len(got) != len(tt.want) {
				t.Errorf("Cache.expiredKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_clearItems(t *testing.T) {
	type args struct {
		keys  []string
		found bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "first",
			args: args{
				keys:  []string{"item", "item1", "item2"},
				found: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(0, 0)
			for _, key := range tt.args.keys {
				c.Set(key, "", 0)
			}
			c.clearItems(tt.args.keys)
			for _, key := range tt.args.keys {
				_, found := c.Get(key)
				if found != tt.args.found {
					t.Errorf("got %v, want %v", found, tt.args.found)
				}
			}
		})
	}
}
