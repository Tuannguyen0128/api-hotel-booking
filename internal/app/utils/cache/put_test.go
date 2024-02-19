package cache

import (
	"sync"
	"testing"
	"time"
)

func Test_util_Put(t *testing.T) {
	data := dataWithId{id: "id1"}

	type fields struct {
		timeout time.Duration
		data    map[string]item[dataWithId]
		keyList []string
		size    int
	}
	type args struct {
		v dataWithId
	}
	tests := map[string]struct {
		name   string
		fields fields
		args   args
	}{
		"success/new": {
			fields: fields{
				timeout: time.Minute,
				data:    map[string]item[dataWithId]{},
				keyList: []string{},
				size:    10,
			},
			args: args{v: data},
		},
		"success/replace": {
			fields: fields{
				timeout: time.Minute,
				data:    map[string]item[dataWithId]{data.id: item[dataWithId]{timeout: time.Now().Add(time.Minute), v: data}},
				keyList: []string{data.id},
				size:    10,
			},
			args: args{v: data},
		},
		"success/oversize": {
			fields: fields{
				timeout: time.Minute,
				data:    map[string]item[dataWithId]{"id0": item[dataWithId]{timeout: time.Now().Add(time.Minute), v: data}},
				keyList: []string{"id0"},
				size:    1,
			},
			args: args{v: data},
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
			u.Put(tt.args.v)
		})
	}
}
