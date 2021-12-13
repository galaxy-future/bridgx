package helper

import (
	"reflect"
	"runtime"

	"github.com/galaxy-future/BridgX/pkg/utils"
	"github.com/gin-gonic/gin"
)

type logReader interface {
	GetOperation(handler string) string
	GetOperationDetail(info string) string
}

var (
	logReaderMap  = make(map[string]logReader)
	defaultReader = defaultLogReader{}
)

func RegisterHandlerLogReader(handlerFunc gin.HandlerFunc, reader logReader) {
	fName := reflectFuncName(handlerFunc)
	logReaderMap[fName] = reader
}

func reflectFuncName(f gin.HandlerFunc) string {
	fName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	return utils.GetStringSuffix(fName)
}

func GetLogReader(fName string) logReader {
	reader, ok := logReaderMap[fName]
	if !ok {
		return defaultReader
	}
	return reader
}

type defaultLogReader struct{}

func (d defaultLogReader) GetOperation(handler string) string {
	return handler
}

func (d defaultLogReader) GetOperationDetail(info string) string {
	return info
}
