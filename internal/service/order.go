package service

import (
	"encoding/json"
	"time"

	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/pkg/cloud"
)

func QueryOrders(accName, provider, ak, regionId string, startTime, endTime time.Time) error {
	cloudCli, err := getProvider(provider, ak, regionId)
	if err != nil {
		return err
	}

	pageNum := 1
	pageSize := 100
	for {
		res, err := cloudCli.GetOrders(cloud.GetOrdersRequest{StartTime: startTime, EndTime: endTime, PageNum: pageNum, PageSize: pageSize})
		if err != nil {
			return err
		}

		if err = SaveOrders(accName, provider, res.Orders); err != nil {
			return err
		}

		if len(res.Orders) < pageSize {
			break
		}
		pageNum += 1
	}
	return nil
}

func SaveOrders(accName, provider string, cloudOrder []cloud.Order) error {
	orderNum := len(cloudOrder)
	if orderNum == 0 {
		return nil
	}

	for _, row := range cloudOrder {
		extend, _ := json.Marshal(row.Extend)
		order := &model.Order{
			AccountName:    accName,
			OrderId:        row.OrderId,
			OrderTime:      row.OrderTime,
			Product:        row.Product,
			Quantity:       row.Quantity,
			UsageStartTime: row.UsageStartTime,
			UsageEndTime:   row.UsageEndTime,
			Provider:       provider,
			RegionId:       row.RegionId,
			ChargeType:     row.ChargeType,
			PayStatus:      row.PayStatus,
			Currency:       row.Currency,
			Cost:           row.Cost,
			Extend:         string(extend[:]),
		}

		if err := model.CreateIgnoreDuplicate(order); err != nil {
			return err
		}
	}

	return nil
}
