package tests

import (
	"errors"
	"testing"
	"time"

	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/galaxy-future/BridgX/pkg/cloud/alibaba"
	"github.com/galaxy-future/BridgX/pkg/cloud/huawei"
	"github.com/galaxy-future/BridgX/pkg/cloud/tencent"
	jsoniter "github.com/json-iterator/go"
)

func getCloudClient() (client cloud.Provider, err error) {
	provider := cloud.TencentCloud
	switch provider {
	case cloud.AlibabaCloud:
		client, err = alibaba.New("ak", "sk", "regionId")
	case cloud.HuaweiCloud:
		client, err = huawei.New("ak", "sk", "regionId")
	case cloud.TencentCloud:
		client, err = tencent.New("ak", "sk", "regionId")
	default:
		return nil, errors.New("invalid provider")
	}
	if err != nil {
		return nil, err
	}
	return client, nil
}

func TestCreateIns(t *testing.T) {
	client, err := getCloudClient()
	if err != nil {
		t.Log(err)
		return
	}

	param := cloud.Params{
		InstanceType: "c6s.large.2",
		ImageId:      "",
		Network: &cloud.Network{
			VpcId:                   "",
			SubnetId:                "",
			SecurityGroup:           "",
			InternetChargeType:      cloud.BandwidthPayByTraffic,
			InternetMaxBandwidthOut: 0,
			InternetIpType:          "5_bgp",
		},
		Disks: &cloud.Disks{
			SystemDisk: cloud.DiskConf{Size: 40, Category: "SSD"},
			DataDisk:   []cloud.DiskConf{},
		},
		Charge: &cloud.Charge{
			ChargeType: cloud.InstanceChargeTypePostPaid,
			Period:     1,
			PeriodUnit: "Month",
		},
		Password: "xxx",
		Tags: []cloud.Tag{
			{
				Key:   cloud.TaskId,
				Value: "12345",
			},
			{
				Key:   cloud.ClusterName,
				Value: "cluster2",
			},
		},
		DryRun: true,
	}
	res, err := client.BatchCreate(param, 1)
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log(res)
}

func TestShowIns(t *testing.T) {
	client, err := getCloudClient()
	if err != nil {
		t.Log(err)
		return
	}

	var res interface{}
	var resStr string
	ids := []string{""}
	res, err = client.GetInstances(ids)
	if err != nil {
		t.Log(err)
		return
	}
	resStr, _ = jsoniter.MarshalToString(res)
	t.Log(resStr)

	tags := []cloud.Tag{{Key: cloud.TaskId, Value: "12345"}}
	res, err = client.GetInstancesByTags("", tags)
	if err != nil {
		t.Log(err)
		return
	}
	resStr, _ = jsoniter.MarshalToString(res)
	t.Log(resStr)
}

func TestCtlIns(t *testing.T) {
	client, err := getCloudClient()
	if err != nil {
		t.Log(err)
		return
	}

	ids := []string{""}

	err = client.StopInstances(ids)
	if err != nil {
		t.Log(err.Error())
	}

	time.Sleep(time.Duration(60) * time.Second)
	err = client.StartInstances(ids)
	if err != nil {
		t.Log(err.Error())
	}

	time.Sleep(time.Duration(60) * time.Second)
	err = client.BatchDelete(ids, "")
	if err != nil {
		t.Log(err.Error())
	}
}

func TestGetResource(t *testing.T) {
	client, err := getCloudClient()
	if err != nil {
		t.Log(err)
		return
	}

	var res interface{}
	var resStr string
	res, err = client.GetRegions()
	if err != nil {
		t.Log(err.Error())
		return
	}
	resStr, _ = jsoniter.MarshalToString(res)
	t.Log(resStr)

	res, err = client.GetZones(cloud.GetZonesRequest{})
	if err != nil {
		t.Log(err.Error())
		return
	}
	resStr, _ = jsoniter.MarshalToString(res)
	t.Log(resStr)

	res, err = client.DescribeAvailableResource(cloud.DescribeAvailableResourceRequest{})
	if err != nil {
		t.Log(err.Error())
		return
	}
	resStr, _ = jsoniter.MarshalToString(res)
	t.Log(resStr)

	res, err = client.DescribeInstanceTypes(cloud.DescribeInstanceTypesRequest{TypeName: []string{"1"}})
	if err != nil {
		t.Log(err.Error())
		return
	}
	resStr, _ = jsoniter.MarshalToString(res)
	t.Log(resStr)

	res, err = client.DescribeImages(cloud.DescribeImagesRequest{InsType: "c6s.large.2"})
	if err != nil {
		t.Log(err.Error())
		return
	}
	resStr, _ = jsoniter.MarshalToString(res)
	t.Log(resStr)
}

func TestCreateSecGrp(t *testing.T) {
	client, err := getCloudClient()
	if err != nil {
		t.Log(err)
		return
	}

	req := cloud.CreateSecurityGroupRequest{
		SecurityGroupName: "test2",
		VpcId:             "",
	}
	res, err := client.CreateSecurityGroup(req)
	if err != nil {
		t.Log(err.Error())
		return
	}
	resStr, _ := jsoniter.MarshalToString(res)
	t.Log(resStr)
}

func TestAddSecGrpRule(t *testing.T) {
	client, err := getCloudClient()
	if err != nil {
		t.Log(err)
		return
	}

	req := cloud.AddSecurityGroupRuleRequest{
		SecurityGroupId: "",
		IpProtocol:      "udp",
		PortFrom:        8894,
		PortTo:          8895,
		CidrIp:          "192.168.1.1/24",
	}
	err = client.AddIngressSecurityGroupRule(req)
	if err != nil {
		t.Log(err.Error())
		return
	}

	req = cloud.AddSecurityGroupRuleRequest{
		SecurityGroupId: "",
		IpProtocol:      "tcp",
		PortFrom:        1000,
		PortTo:          1000,
		CidrIp:          "192.168.1.1/24",
	}
	err = client.AddEgressSecurityGroupRule(req)
	if err != nil {
		t.Log(err.Error())
		return
	}
}

func TestShowSecGrp(t *testing.T) {
	client, err := getCloudClient()
	if err != nil {
		t.Log(err)
		return
	}

	var res interface{}
	var resStr string
	res, err = client.DescribeSecurityGroups(cloud.DescribeSecurityGroupsRequest{
		VpcId: "",
	})
	if err != nil {
		t.Log(err.Error())
		return
	}
	resStr, _ = jsoniter.MarshalToString(res)
	t.Log(resStr)

	res, err = client.DescribeGroupRules(cloud.DescribeGroupRulesRequest{
		SecurityGroupId: "",
	})
	if err != nil {
		t.Log(err.Error())
		return
	}
	resStr, _ = jsoniter.MarshalToString(res)
	t.Log(resStr)
}

func TestCreateVpc(t *testing.T) {
	client, err := getCloudClient()
	if err != nil {
		t.Log(err)
		return
	}

	var resStr string
	vpc, err := client.CreateVPC(cloud.CreateVpcRequest{
		VpcName:   "vpc1",
		CidrBlock: "10.8.0.0/16",
	})
	if err != nil {
		t.Log(err.Error())
		return
	}
	resStr, _ = jsoniter.MarshalToString(vpc)
	t.Log(resStr)
}

func TestCreateSubnet(t *testing.T) {
	client, err := getCloudClient()
	if err != nil {
		t.Log(err)
		return
	}

	var res interface{}
	var resStr string

	vpcId := ""
	res, err = client.CreateSwitch(cloud.CreateSwitchRequest{
		ZoneId:      "",
		CidrBlock:   "10.8.0.0/18",
		VSwitchName: "subnet1",
		VpcId:       vpcId,
		GatewayIp:   "10.8.63.254",
	})
	if err != nil {
		t.Log(err.Error())
		return
	}
	resStr, _ = jsoniter.MarshalToString(res)
	t.Log(resStr)
}

func TestShowVpc(t *testing.T) {
	client, err := getCloudClient()
	if err != nil {
		t.Log(err)
		return
	}

	var res interface{}
	var resStr string
	vpcId := ""
	swId := ""
	res, err = client.GetVPC(cloud.GetVpcRequest{
		VpcId: vpcId,
	})
	if err != nil {
		t.Log(err.Error())
		return
	}
	resStr, _ = jsoniter.MarshalToString(res)
	t.Log(resStr)

	res, err = client.DescribeVpcs(cloud.DescribeVpcsRequest{})
	if err != nil {
		t.Log(err.Error())
		return
	}
	resStr, _ = jsoniter.MarshalToString(res)
	t.Log(resStr)

	res, err = client.GetSwitch(cloud.GetSwitchRequest{
		SwitchId: swId,
	})
	if err != nil {
		t.Log(err.Error())
		return
	}
	resStr, _ = jsoniter.MarshalToString(res)
	t.Log(resStr)

	res, err = client.DescribeSwitches(cloud.DescribeSwitchesRequest{
		VpcId: vpcId,
	})
	if err != nil {
		t.Log(err.Error())
		return
	}
	resStr, _ = jsoniter.MarshalToString(res)
	t.Log(resStr)
}

func TestQueryOrders(t *testing.T) {
	cli, err := getCloudClient()
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
		res, err := cli.GetOrders(cloud.GetOrdersRequest{StartTime: startTime, EndTime: endTime,
			PageNum: pageNum, PageSize: pageSize})
		if err != nil {
			t.Log(err.Error())
			return
		}
		cnt := 0
		t.Log("len:", len(res.Orders))
		for _, row := range res.Orders {
			cnt += 1
			if cnt > 3 {
				t.Log("---------------")
				break
			}
			rowStr, _ := jsoniter.MarshalToString(row)
			t.Log(rowStr)
		}
		if len(res.Orders) < pageSize {
			break
		}
		pageNum += 1
	}
	t.Log(pageNum)
}
