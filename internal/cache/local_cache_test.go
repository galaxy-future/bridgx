package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/galaxy-future/BridgX/internal/logs"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack/v5"
)

type Base struct {
	Id       int64      `json:"id" gorm:"primary_key"`
	CreateAt *time.Time `json:"-"`
	UpdateAt *time.Time `json:"-"`
}

type User struct {
	Base
	Username   string `json:"username"`
	Password   string `json:"password"`
	UserType   int8   `json:"user_type"`
	UserStatus string `json:"user_status"`
	OrgId      int64  `json:"org_id"`
	CreateBy   string `json:"create_by"`
}

func TestGetThroughBigCache(t *testing.T) {

	logs.Init()
	MustInit()

	b, _ := msgpack.Marshal(User{
		Base:       Base{Id: 1},
		Username:   "root",
		Password:   "123",
		UserType:   1,
		UserStatus: "ddd",
		OrgId:      1231283123120998812,
		CreateBy:   "root",
	})
	bigLocalCache.Set("1", b)
	type args struct {
		ids      []int64
		out      interface{}
		keyMaker func(int64) string
		delegate func([]int64, interface{}) error
	}
	users := []*User{}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"GET",
			args{
				ids: []int64{1, 2},
				out: &users,
				keyMaker: func(i int64) string {
					return fmt.Sprintf("%d", i)
				},
				delegate: nil,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var needFetch []int64
			var err error
			if needFetch, err = GetFromBigCache(tt.args.ids, tt.args.out, tt.args.keyMaker); (err != nil) != tt.wantErr {
				t.Errorf("GetFromBigCache() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Len(t, needFetch, 1)
			b, _ := jsoniter.MarshalToString(tt.args.out)
			t.Logf("out:%v", b)
		})
	}
}

func TestSlice(t *testing.T) {
	input := make([]string, 0)
	input = append(input, "1234")
	changeValue(input)
	fmt.Println(input)
}

func changeValue(input []string) {
	input[0] = "xx"
}
