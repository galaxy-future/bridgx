package helper

import (
	"fmt"
	"strconv"

	"github.com/galaxy-future/BridgX/internal/constants"
	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
	"github.com/gin-gonic/gin"
)

func GetPagerParamFromQuery(c *gin.Context) (pageNumber int, pageSize int) {
	pageNumberContent := c.Query("page_number")
	PageSizeContent := c.Query("page_size")

	if pageNumberContent != "" {
		value, err := strconv.ParseInt(pageNumberContent, 10, 60)
		if err == nil {
			pageNumber = int(value)
		}
	}
	if PageSizeContent != "" {
		value, err := strconv.ParseInt(PageSizeContent, 10, 60)
		if err == nil {
			pageSize = int(value)
		}
	}

	if pageNumber < 1 {
		pageNumber = 1
	}
	if pageSize < 1 || pageSize > constants.DefaultPageSize {
		pageSize = constants.DefaultPageSize
	}
	return pageNumber, pageSize
}

func GetUserToken(c *gin.Context) (string, error) {
	value := c.Value(gf_cluster.HeaderTokenName)
	if value == nil {
		return "", fmt.Errorf("获取登录信息失败")
	}

	return value.(string), nil
}
