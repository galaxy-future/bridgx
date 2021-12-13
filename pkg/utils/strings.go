package utils

import (
	"encoding/json"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

const defaultSep = "."

func Interface2String(value interface{}) string {
	key := ""
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	case json.Number:
		key = value.(json.Number).String()
	}

	return key
}

func ObjToJson(obj interface{}) string {
	bytes, _ := jsoniter.Marshal(obj)
	return string(bytes)
}

func StringSliceSplit(slice []string, singleLen int64) [][]string {
	var sliceList [][]string
	sliceLen := int64(len(slice))
	if sliceLen <= singleLen {
		sliceList = append(sliceList, slice)
		return sliceList
	}
	var start, end int64
	cnt := sliceLen / singleLen
	for i := int64(0); i < cnt; i++ {
		start = i * singleLen
		end = (i + 1) * singleLen
		if end > sliceLen { // 不超过最大长度
			end = sliceLen
		}
		newSlice := slice[start:end]
		sliceList = append(sliceList, newSlice)
	}
	//如有剩余元素，单独 append 到 sliceList
	if end < sliceLen {
		start = end
		end = sliceLen
		newSlice := slice[start:end]
		sliceList = append(sliceList, newSlice)
	}
	return sliceList
}

func GetStringSuffix(s string, sep ...string) string {
	if len(s) == 0 {
		return ""
	}
	seperator := defaultSep
	if len(sep) > 0 {
		seperator = sep[0]
	}
	sl := strings.Split(s, seperator)
	if len(sl) == 0 {
		return ""
	}
	return sl[len(sl)-1]
}
