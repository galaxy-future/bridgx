package tests

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/galaxy-future/BridgX/pkg/cloud/alibaba"
)

func TestGetAlibabaCloudClient(t *testing.T) {
	c, err := alibaba.New("a", "b", "cn-beijing")
	t.Logf("err:%v\n", err)
	region, err := c.GetRegions()

	t.Logf("err:%v\n", err)
	t.Logf("regions:%v\n", region)

}

func TestQueryOrders(t *testing.T) {
	cloudCli, err := alibaba.New("a", "b", "cn-beijing")
	if err != nil {
		t.Log(err.Error())
		return
	}

	//endTime := time.Now().UTC()
	//duration, _ := time.ParseDuration("-5h")
	//startTime := endTime.Add(duration)
	startTime, _ := time.Parse("2006-01-02 15:04:05", "2021-11-19 11:40:02")
	endTime, _ := time.Parse("2006-01-02 15:04:05", "2021-11-19 11:45:02")
	pageNum := 1
	pageSize := 100
	for {
		res, err := cloudCli.GetOrders(cloud.GetOrdersRequest{StartTime: startTime, EndTime: endTime,
			PageNum: pageNum, PageSize: pageSize})
		if err != nil {
			t.Log(err.Error())
			return
		}
		cnt := 0
		t.Log("len:", len(res.Orders))
		for _, row := range res.Orders {
			cnt += 1
			if cnt > alibaba.SubOrderNumPerMain {
				t.Log("---------------")
				break
			}
			rowStr, _ := json.Marshal(row)
			t.Log(string(rowStr))
		}
		if len(res.Orders) < pageSize {
			break
		}
		pageNum += 1
	}
	t.Log(pageNum)
}

func TestGetOrderDetail(t *testing.T) {
	client, err := bssopenapi.NewClientWithAccessKey("cn-beijing", "a", "b")
	if err != nil {
		t.Log(err.Error())
		return
	}
	request := bssopenapi.CreateGetOrderDetailRequest()
	request.Scheme = "https"
	request.OrderId = "211577282350149"
	response, err := client.GetOrderDetail(request)
	if err != nil {
		t.Log(err.Error())
		return
	}

	orders, err := json.Marshal(response.Data.OrderList)
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log(string(orders))
}
