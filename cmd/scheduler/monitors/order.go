package monitors

import (
	"time"

	"github.com/galaxy-future/BridgX/internal/bcc"
	"github.com/galaxy-future/BridgX/internal/constants"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/internal/service"
	"github.com/spf13/cast"
)

type QueryOrderJobs struct {
}

func (p *QueryOrderJobs) Run() {
	accounts := make([]model.Account, 0)
	if err := model.QueryAll(map[string]interface{}{}, &accounts, ""); err != nil {
		logs.Logger.Error("query account,", err)
		return
	}
	if len(accounts) == 0 {
		return
	}

	startTime, endTime, err := getQueryTimeRange()
	if err != nil {
		logs.Logger.Error("getQueryTimeRange,", err)
		return
	}
	logs.Logger.Debug("start,", startTime, endTime)
	for _, row := range accounts {
		region, err := model.GetOneRegionByAccKey(row.AccountKey)
		if err != nil {
			logs.Logger.Error("GetOneRegionByAccKey,", err)
			return
		}
		if err = service.QueryOrders(row.AccountName, row.Provider, row.AccountKey, region.RegionId,
			startTime, endTime); err != nil {
			logs.Logger.Error("queryOrders,", err)
			return
		}
	}

	updateQueryStartTime(endTime)
	logs.Logger.Debug("end,", endTime)
}

func getQueryTimeRange() (time.Time, time.Time, error) {
	var startTime, endTime time.Time
	startTimeStr, err := bcc.GetConfig(constants.CostConfigGroup, constants.QueryOrderStartTime)
	if err != nil {
		return startTime, endTime, err
	}
	logs.Logger.Debug("getConfig:", startTimeStr)
	if startTimeStr == "" {
		endTime = time.Now().UTC()
		duration, _ := time.ParseDuration("-" + cast.ToString(constants.DefaultQueryOrderInterval) + "s")
		startTime = endTime.Add(duration)
		return startTime, endTime, nil
	}

	startTime, err = time.Parse("2006-01-02 15:04:05", startTimeStr)
	duration, _ := time.ParseDuration(cast.ToString(constants.DefaultQueryOrderInterval) + "s")
	endTime = startTime.Add(duration)
	return startTime, endTime, nil
}

func updateQueryStartTime(startTime time.Time) {
	startTimeStr := startTime.Format("2006-01-02 15:04:05")
	err := bcc.PublishConfig(constants.CostConfigGroup, constants.QueryOrderStartTime, startTimeStr)
	if err != nil {
		logs.Logger.Error(err)
		return
	}
}
