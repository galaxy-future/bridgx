package model

import "time"

type Order struct {
	Base
	AccountName    string    `json:"account_name"`
	OrderId        string    `json:"order_id"`
	OrderTime      time.Time `json:"order_time"`
	Product        string    `json:"product"`
	Quantity       int32     `json:"quantity"`
	UsageStartTime time.Time `json:"usage_start_time"`
	UsageEndTime   time.Time `json:"usage_end_time"`
	Provider       string    `json:"provider"`
	RegionId       string    `json:"region_id"`
	ChargeType     string    `json:"charge_type"`
	PayStatus      int8      `json:"pay_status"`
	Currency       string    `json:"currency"`
	Cost           float32   `json:"cost"`
	Extend         string    `json:"extend"`
}

func (Order) TableName() string {
	return "order_202101"
}
