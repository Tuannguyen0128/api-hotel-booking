package cache

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func Test_util_Get(t *testing.T) {
	data := dataWithId{id: "id1"}

	type fields struct {
		timeout time.Duration
		data    map[string]item[dataWithId]
		keyList []string
		size    int
	}
	type args struct {
		k string
	}
	tests := map[string]struct {
		name   string
		fields fields
		args   args
		want   dataWithId
		found  bool
	}{
		"success/found": {
			fields: fields{
				timeout: time.Minute,
				data:    map[string]item[dataWithId]{data.id: item[dataWithId]{timeout: time.Now().Add(time.Minute), v: data}},
				keyList: []string{data.id},
				size:    10,
			},
			args:  args{k: data.id},
			want:  data,
			found: true,
		},
		"success/not-found": {
			fields: fields{
				timeout: time.Minute,
				data:    map[string]item[dataWithId]{},
				keyList: []string{},
				size:    10,
			},
			args:  args{k: data.id},
			want:  dataWithId{},
			found: false,
		},
		"success/timeout": {
			fields: fields{
				timeout: time.Minute,
				data:    map[string]item[dataWithId]{data.id: item[dataWithId]{timeout: time.Now().Add(-time.Minute), v: data}},
				keyList: []string{data.id},
				size:    10,
			},
			args:  args{k: data.id},
			want:  dataWithId{},
			found: false,
		},
	}
	for n, tt := range tests {
		tt.name = n
		t.Run(tt.name, func(t *testing.T) {
			u := &util[dataWithId]{
				timeout: tt.fields.timeout,
				data:    tt.fields.data,
				keyList: tt.fields.keyList,
				size:    tt.fields.size,
				lock:    sync.Mutex{},
			}
			got, got1 := u.Get(tt.args.k)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.found {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.fields)
			}
		})
	}
}
