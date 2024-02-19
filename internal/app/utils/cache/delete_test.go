package cache

import (
	"sync"
	"testing"
	"time"
)

type dataWithId struct {
	id string
}

func (d dataWithId) GetId() string {
	return d.id
}

func Test_util_Delete(t *testing.T) {
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
		want   bool
	}{
		"success/no-data": {
			fields: fields{
				timeout: time.Minute,
				data:    map[string]item[dataWithId]{},
				keyList: []string{},
				size:    10,
			},
			args: args{k: "id1"},
			want: false,
		},
		"success/have-data": {
			fields: fields{
				timeout: time.Minute,
				data:    map[string]item[dataWithId]{"id1": item[dataWithId]{timeout: time.Now().Add(time.Minute), v: dataWithId{id: "id1"}}},
				keyList: []string{"id1"},
				size:    10,
			},
			args: args{k: "id1"},
			want: true,
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
			if got := u.Delete(tt.args.k); got != tt.want {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}
