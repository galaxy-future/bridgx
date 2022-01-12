package huawei

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v2/model"
	"github.com/spf13/cast"
)

// CreateVPC 返回缺少RequestId
func (p *HuaweiCloud) CreateVPC(req cloud.CreateVpcRequest) (cloud.CreateVpcResponse, error) {
	request := &model.CreateVpcRequest{}
	vpcbody := &model.CreateVpcOption{
		Cidr: &req.CidrBlock,
		Name: &req.VpcName,
	}
	request.Body = &model.CreateVpcRequestBody{
		Vpc: vpcbody,
	}
	response, err := p.vpcClient.CreateVpc(request)
	if err != nil {
		return cloud.CreateVpcResponse{}, err
	}
	if response.HttpStatusCode != http.StatusOK {
		return cloud.CreateVpcResponse{}, fmt.Errorf("httpcode %d", response.HttpStatusCode)
	}

	res := cloud.CreateVpcResponse{
		VpcId:     response.Vpc.Id,
		RequestId: "",
	}
	return res, nil
}

func (p *HuaweiCloud) GetVPC(req cloud.GetVpcRequest) (cloud.GetVpcResponse, error) {
	request := &model.ShowVpcRequest{
		VpcId: req.VpcId,
	}
	response, err := p.vpcClient.ShowVpc(request)
	if err != nil {
		return cloud.GetVpcResponse{}, err
	}
	if response.HttpStatusCode != http.StatusOK {
		return cloud.GetVpcResponse{}, fmt.Errorf("httpcode %d", response.HttpStatusCode)
	}

	vpc := vpcInfo2CloudVpc([]model.Vpc{*response.Vpc}, req.RegionId)
	return cloud.GetVpcResponse{Vpc: vpc[0]}, nil
}

func (p *HuaweiCloud) DescribeVpcs(req cloud.DescribeVpcsRequest) (cloud.DescribeVpcsResponse, error) {
	vpcs := make([]model.Vpc, 0, 16)
	request := &model.ListVpcsRequest{}
	limitRequest := int32(_pageSize)
	request.Limit = &limitRequest
	markerRequest := ""
	for {
		if markerRequest != "" {
			request.Marker = &markerRequest
		}
		response, err := p.vpcClient.ListVpcs(request)
		if err != nil {
			return cloud.DescribeVpcsResponse{}, err
		}
		if response.HttpStatusCode != http.StatusOK {
			return cloud.DescribeVpcsResponse{}, fmt.Errorf("httpcode %d", response.HttpStatusCode)
		}

		vpcs = append(vpcs, *response.Vpcs...)
		vpcNum := len(*response.Vpcs)
		if vpcNum < _pageSize {
			break
		}
		markerRequest = (*response.Vpcs)[vpcNum-1].Id
	}

	return cloud.DescribeVpcsResponse{Vpcs: vpcInfo2CloudVpc(vpcs, req.RegionId)}, nil
}

// CreateSwitch add GatewayIp,miss RequestId
func (p *HuaweiCloud) CreateSwitch(req cloud.CreateSwitchRequest) (cloud.CreateSwitchResponse, error) {
	request := &model.CreateSubnetRequest{}
	subnetbody := &model.CreateSubnetOption{
		Name:      req.VSwitchName,
		Cidr:      req.CidrBlock,
		VpcId:     req.VpcId,
		GatewayIp: req.GatewayIp,
	}
	if req.ZoneId != "" {
		subnetbody.AvailabilityZone = &req.ZoneId
	}
	request.Body = &model.CreateSubnetRequestBody{
		Subnet: subnetbody,
	}
	response, err := p.vpcClient.CreateSubnet(request)
	if err != nil {
		return cloud.CreateSwitchResponse{}, err
	}
	if response.HttpStatusCode != http.StatusOK {
		return cloud.CreateSwitchResponse{}, fmt.Errorf("httpcode %d", response.HttpStatusCode)
	}

	return cloud.CreateSwitchResponse{SwitchId: response.Subnet.Id}, nil
}

func (p *HuaweiCloud) GetSwitch(req cloud.GetSwitchRequest) (cloud.GetSwitchResponse, error) {
	request := &model.ShowSubnetRequest{
		SubnetId: req.SwitchId,
	}
	response, err := p.vpcClient.ShowSubnet(request)
	if err != nil {
		return cloud.GetSwitchResponse{}, err
	}
	if response.HttpStatusCode != http.StatusOK {
		return cloud.GetSwitchResponse{}, fmt.Errorf("httpcode %d", response.HttpStatusCode)
	}

	usedIpNum, err := p.getUsedIpNum([]string{req.SwitchId})
	if err != nil {
		return cloud.GetSwitchResponse{}, err
	}
	s := subnetInfo2CloudSwitch([]model.Subnet{*response.Subnet}, usedIpNum)
	return cloud.GetSwitchResponse{Switch: s[0]}, nil
}

func (p *HuaweiCloud) DescribeSwitches(req cloud.DescribeSwitchesRequest) (cloud.DescribeSwitchesResponse, error) {
	subnets := make([]model.Subnet, 0, _pageSize)
	swIds := make([]string, 0, _pageSize)
	request := &model.ListSubnetsRequest{}
	limitRequest := int32(_pageSize)
	request.Limit = &limitRequest
	vpcIdRequest := req.VpcId
	request.VpcId = &vpcIdRequest
	markerRequest := ""
	for {
		if markerRequest != "" {
			request.Marker = &markerRequest
		}
		response, err := p.vpcClient.ListSubnets(request)
		if err != nil {
			return cloud.DescribeSwitchesResponse{}, err
		}
		if response.HttpStatusCode != http.StatusOK {
			return cloud.DescribeSwitchesResponse{}, fmt.Errorf("httpcode %d", response.HttpStatusCode)
		}

		subnets = append(subnets, *response.Subnets...)
		for _, subnet := range *response.Subnets {
			swIds = append(swIds, subnet.Id)
		}
		netNum := len(*response.Subnets)
		if netNum < _pageSize {
			break
		}
		markerRequest = (*response.Subnets)[netNum-1].Id
	}

	usedIpNum, err := p.getUsedIpNum(swIds)
	if err != nil {
		return cloud.DescribeSwitchesResponse{}, err
	}
	return cloud.DescribeSwitchesResponse{Switches: subnetInfo2CloudSwitch(subnets, usedIpNum)}, nil
}

//miss CreateAt
func vpcInfo2CloudVpc(vpcInfo []model.Vpc, regionId string) []cloud.VPC {
	vpcs := make([]cloud.VPC, 0, len(vpcInfo))
	for _, vpc := range vpcInfo {
		stat, _ := vpc.Status.MarshalJSON()
		vpcs = append(vpcs, cloud.VPC{
			VpcId:     vpc.Id,
			VpcName:   vpc.Name,
			CidrBlock: vpc.Cidr,
			RegionId:  regionId,
			Status:    _vpcStatus[string(stat)],
		})
	}
	return vpcs
}

//miss IsDefault,CreateAt
func subnetInfo2CloudSwitch(subnetInfo []model.Subnet, UsedIpNum map[string]int) []cloud.Switch {
	switchs := make([]cloud.Switch, 0, len(subnetInfo))
	for _, subnet := range subnetInfo {
		stat, _ := subnet.Status.MarshalJSON()
		totalHostNum := getSubnetTotalHostNum(subnet.Cidr)

		switchs = append(switchs, cloud.Switch{
			VpcId:                   subnet.VpcId,
			SwitchId:                subnet.Id,
			Name:                    subnet.Name,
			AvailableIpAddressCount: totalHostNum - 1 - UsedIpNum[subnet.Id],
			VStatus:                 _subnetStatus[string(stat)],
			ZoneId:                  subnet.AvailabilityZone,
			CidrBlock:               subnet.Cidr,
			GatewayIp:               subnet.GatewayIp,
		})
	}
	return switchs
}

func (p *HuaweiCloud) getUsedIpNum(switchIds []string) (map[string]int, error) {
	resMap := make(map[string]int, len(switchIds))
	request := &model.ListPrivateipsRequest{}
	for _, switchId := range switchIds {
		request.SubnetId = switchId
		response, err := p.vpcClient.ListPrivateips(request)
		if err != nil {
			return nil, err
		}
		if response.HttpStatusCode != http.StatusOK {
			logs.Logger.Errorf("%s, httpcode %d", switchId, response.HttpStatusCode)
			continue
		}

		resMap[switchId] = len(*response.Privateips)
	}

	return resMap, nil
}

func getSubnetTotalHostNum(cidr string) int {
	index := strings.Index(cidr, "/")
	if index < 0 {
		return 0
	}
	num := cast.ToInt(cidr[index+1:])
	if num < 1 || num > 31 {
		return 0
	}

	return 1<<(32-num) - 2
}
