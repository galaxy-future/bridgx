package tests

import (
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	"github.com/galaxy-future/BridgX/pkg/cloud/aliyun"
)

func TestGetAliyunClient(t *testing.T) {
	c, err := aliyun.New("a", "b", "cn-beijing")
	t.Logf("err:%v\n", err)
	region, err := c.GetRegions()

	t.Logf("err:%v\n", err)
	t.Logf("regions:%v\n", region)

}

func TestQueryOrders(t *testing.T) {
	client, err := bssopenapi.NewClientWithAccessKey("cn-beijing", "a", "b")
	if err != nil {
		t.Log(err.Error())
		return
	}

	endTime := time.Now().UTC()
	duration, _ := time.ParseDuration("-5h")
	startTime := endTime.Add(duration)
	pageNum := 1
	pageSize := 100
	request := bssopenapi.CreateQueryOrdersRequest()
	request.Scheme = "https"
	request.CreateTimeStart = startTime.Format("2006-01-02T15:04:05Z")
	request.CreateTimeEnd = endTime.Format("2006-01-02T15:04:05Z")
	request.PageSize = requests.NewInteger(pageSize)
	detailReq := bssopenapi.CreateGetOrderDetailRequest()
	detailReq.Scheme = "https"
	for {
		request.PageNum = requests.NewInteger(pageNum)
		response, err := client.QueryOrders(request)
		if err != nil {
			t.Log(err.Error())
			return
		}
		if !response.Success {
			t.Log(response.Message)
			return
		}
		pageNum = response.Data.PageNum + 1
		if len(response.Data.OrderList.Order) == 0 {
			break
		}
		cnt := 0
		for _, row := range response.Data.OrderList.Order {
			if row.PretaxAmount == "0" {
				//continue
			}

			t.Log("main", row)

			detailReq.OrderId = row.OrderId
			detailRsp, err := client.GetOrderDetail(detailReq)
			if err != nil {
				t.Log(err.Error())
				return
			}
			if !detailRsp.Success {
				t.Log(detailRsp.Message)
				continue
			}
			if len(detailRsp.Data.OrderList.Order) == 0 {
				continue
			}
			for _, subOrder := range detailRsp.Data.OrderList.Order {
				t.Log(subOrder)
			}

			cnt++
			if cnt > 1 {
				break
			}
		}

		if response.Data.PageNum*response.Data.PageSize >= response.Data.TotalCount {
			break
		}
	}
	t.Log("end,", pageNum)
}

func TestGetOrderDetail(t *testing.T) {
	client, err := bssopenapi.NewClientWithAccessKey("cn-beijing", "a", "b")
	if err != nil {
		t.Log(err.Error())
		return
	}
	request := bssopenapi.CreateGetOrderDetailRequest()
	request.Scheme = "https"
	request.OrderId = "211518876370341"
	response, err := client.GetOrderDetail(request)
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log(response.Data.OrderList)
}
