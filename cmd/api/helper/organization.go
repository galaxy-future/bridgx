package helper

import (
	"github.com/galaxy-future/BridgX/cmd/api/response"
	"github.com/galaxy-future/BridgX/internal/model"
)

func ConvertToOrgInfo(org *model.Org) response.OrgInfo {
	return response.OrgInfo{
		Id:       org.Id,
		OrgName:  org.OrgName,
		CreateAt: org.CreateAt.Format("2006-01-02 15:04:05"),
		UpdateAt: org.UpdateAt.Format("2006-01-02 15:04:05"),
	}
}
