package cache

import (
	"fmt"
	"reflect"
	"time"

	"github.com/allegro/bigcache"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/pkg/utils"
	"github.com/vmihailenco/msgpack/v5"
)

var bigLocalCache *bigcache.BigCache

func MustInit() {
	bigLocalCache = initBigCache()
}

func initBigCache() *bigcache.BigCache {
	bigCacheMaxLength := 100000 //max length in big cache
	bigCacheJanitor := 2
	bigCacheExpireTime := 60
	bigCacheHardMaxCacheSize := 1024
	bigCacheConfig := bigcache.Config{
		Shards:             1024, //必须2的次幂
		LifeWindow:         time.Duration(bigCacheExpireTime) * time.Second,
		CleanWindow:        time.Duration(bigCacheJanitor) * time.Second,
		MaxEntriesInWindow: bigCacheMaxLength,
		MaxEntrySize:       4096, //单位Bytes
		Verbose:            true,
		HardMaxCacheSize:   bigCacheHardMaxCacheSize, //单位MB
		OnRemove:           nil,
		OnRemoveWithReason: nil,
	}
	BigLocalCache, initErr := bigcache.NewBigCache(bigCacheConfig)
	if initErr != nil {
		panic("init BigCache failed")
	}
	logs.Logger.Infof("[big_cache_init], bigCache config, bigCacheMaxLength=%v, bigCacheJanitor=%v, bigCacheExpireTime=%v, bigCacheHardMaxCacheSize=%v", bigCacheMaxLength, bigCacheJanitor, bigCacheExpireTime, bigCacheHardMaxCacheSize)
	return BigLocalCache
}

//GetFromBigCache
//inputs:
//	ids should be unique ids
//	out should be a pointer to a slice contains your custom type pointer, which will receive the outputs
//	keyMaker should be a func will format the key of ids
//outputs:
//	[]int64 is the missing ids, need fetch from your distributed cache or database
//	error return error(if it has)
func GetFromBigCache(ids []int64, out interface{}, keyMaker func(int64) string) ([]int64, error) {
	if len(ids) == 0 {
		return ids, nil
	}
	err := checkType(out)
	if err != nil {
		return ids, err
	}
	itemTypePtr := reflect.TypeOf(out).Elem().Elem()

	rv := reflect.ValueOf(out).Elem()
	sliceType := reflect.TypeOf(out).Elem()
	sliceReflect := reflect.MakeSlice(sliceType, 0, len(ids))
	foundIds := make([]int64, 0, len(ids))

	for _, id := range ids {
		key := keyMaker(id)
		val, err := bigLocalCache.Get(key)
		if err != nil {
			if err != bigcache.ErrEntryNotFound {
				logs.Logger.Errorf("get key:%v from big cache failed, err:%v", key, err)
			}
			continue
		}
		if val != nil {
			rp := reflect.New(itemTypePtr).Interface()
			err = msgpack.Unmarshal(val, rp)
			if err != nil {
				continue
			}
			foundIds = append(foundIds, id)
			sliceReflect = reflect.Append(sliceReflect, reflect.ValueOf(rp).Elem())
		}
	}
	rv.Set(sliceReflect)
	needFetchIds := ids
	if len(foundIds) > 0 {
		needFetchIds = utils.Filter(ids, func(input int64) bool {
			for _, id := range foundIds {
				if input == id {
					return false
				}
			}
			return true
		})
	}

	return needFetchIds, nil
}

func checkType(out interface{}) error {
	v := reflect.TypeOf(out)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("type :%v not supported, need ptr", v.Kind())
	}
	if v.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("type :%v not supported, need slice", v.Kind())
	}
	if v.Elem().Elem().Kind() != reflect.Ptr {
		return fmt.Errorf("type :%v not supported, need ptr", v.Kind())
	}
	return nil
}

func SetBigCache(id int64, v interface{}, keyMaker func(int64) string) error {
	b, err := msgpack.Marshal(v)
	if err != nil {
		return err
	}
	return bigLocalCache.Set(keyMaker(id), b)
}

func UserKeyMaker(id int64) string {
	return fmt.Sprintf("u_%v", id)
}
