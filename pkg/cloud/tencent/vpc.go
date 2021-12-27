package tencent

import (
	"errors"
	"strconv"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"

	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func (p *TencentCloud) CreateVPC(req cloud.CreateVpcRequest) (cloud.CreateVpcResponse, error) {
	request := vpc.NewCreateVpcRequest()
	request.VpcName = &req.VpcName
	request.CidrBlock = &req.CidrBlock
	response, err := p.vpcClient.CreateVpc(request)
	if err != nil {
		logs.Logger.Errorf("CreateVPC TencentCloud failed.err: [%v], req[%v]", err, req)
		return cloud.CreateVpcResponse{}, err
	}
	if response == nil {
		logs.Logger.Errorf("CreateVPC TencentCloud failed, response is nil, req[%v]", req)
		return cloud.CreateVpcResponse{}, _errResponseIsNil
	}
	res := cloud.CreateVpcResponse{
		VpcId:     *response.Response.Vpc.VpcId,
		RequestId: *response.Response.RequestId,
	}
	return res, nil
}

func (p *TencentCloud) GetVPC(req cloud.GetVpcRequest) (cloud.GetVpcResponse, error) {
	request := vpc.NewDescribeVpcsRequest()
	request.VpcIds = append(request.VpcIds, &req.VpcId)
	//每次请求的vpc实例Id的上限为100。参数不支持同时指定VpcIds和Filters
	response, err := p.vpcClient.DescribeVpcs(request)
	if err != nil {
		logs.Logger.Errorf("GetVPC TencentCloud failed.err: [%v], req[%v]", err, req)
		return cloud.GetVpcResponse{}, err
	}
	if response == nil {
		logs.Logger.Errorf("GetVPC TencentCloud failed, response is nil, req[%v]", req)
		return cloud.GetVpcResponse{}, errors.New("response is nil")
	}
	if *response.Response.TotalCount != 1 {
		logs.Logger.Errorf("GetVPC TencentCloud failed, totalCount isn't one, req[%v]", req)
		return cloud.GetVpcResponse{}, _errIsNotOne
	}
	vpc := response.Response.VpcSet[0]
	switches, err := p.DescribeSwitches(cloud.DescribeSwitchesRequest{VpcId: *vpc.VpcId})
	if err != nil {
		return cloud.GetVpcResponse{}, err
	}
	switchIds := make([]string, 0, len(switches.Switches))
	for _, row := range switches.Switches {
		switchIds = append(switchIds, row.SwitchId)
	}
	res := cloud.GetVpcResponse{
		Vpc: vpcInfo2CloudVpc(switchIds, req.RegionId, vpc),
	}
	return res, nil
}

func (p *TencentCloud) DescribeVpcs(req cloud.DescribeVpcsRequest) (cloud.DescribeVpcsResponse, error) {
	vpcs := make([]*vpc.Vpc, 0, _pageSize)
	request := vpc.NewDescribeVpcsRequest()
	offset := _offset
	request.Limit = common.StringPtr(strconv.Itoa(_pageSize))
	for {
		request.Offset = common.StringPtr(strconv.Itoa(offset))
		response, err := p.vpcClient.DescribeVpcs(request)
		if err != nil {
			logs.Logger.Errorf("DescribeVpcs TencentCloud failed.err: [%v], req[%v]", err, req)
			return cloud.DescribeVpcsResponse{}, err
		}
		if response == nil {
			logs.Logger.Errorf("DescribeVpcs TencentCloud failed, response is nil, req[%v]", req)
			return cloud.DescribeVpcsResponse{}, _errResponseIsNil
		}
		vpcs = append(vpcs, response.Response.VpcSet...)
		if len(response.Response.VpcSet) < _pageSize {
			break
		}
		offset = (offset - 1) * _pageSize
	}
	cloudVpcs := make([]cloud.VPC, 0, len(vpcs))
	for _, vpc := range vpcs {
		switches, err := p.DescribeSwitches(cloud.DescribeSwitchesRequest{VpcId: *vpc.VpcId})
		if err != nil {
			return cloud.DescribeVpcsResponse{}, err
		}
		switchIds := make([]string, 0, len(switches.Switches))
		for _, row := range switches.Switches {
			switchIds = append(switchIds, row.SwitchId)
		}
		cloudVpcs = append(cloudVpcs, vpcInfo2CloudVpc(switchIds, req.RegionId, vpc))
	}
	return cloud.DescribeVpcsResponse{Vpcs: cloudVpcs}, nil
}

func vpcInfo2CloudVpc(switchIds []string, regionId string, vpc *vpc.Vpc) cloud.VPC {
	return cloud.VPC{
		VpcId:     *vpc.VpcId,
		VpcName:   *vpc.VpcName,
		CidrBlock: *vpc.CidrBlock,
		SwitchIds: switchIds,
		RegionId:  regionId,
		Status:    cloud.VPCStatusAvailable,
		CreateAt:  *vpc.CreatedTime,
	}
}

func (p *TencentCloud) CreateSwitch(req cloud.CreateSwitchRequest) (cloud.CreateSwitchResponse, error) {
	request := vpc.NewCreateSubnetRequest()
	request.VpcId = &req.VpcId
	request.SubnetName = &req.VSwitchName
	request.CidrBlock = &req.CidrBlock
	request.Zone = &req.ZoneId
	response, err := p.vpcClient.CreateSubnet(request)
	if err != nil {
		logs.Logger.Errorf("CreateSwitch TencentCloud failed.err: [%v], req[%v]", err, req)
		return cloud.CreateSwitchResponse{}, err
	}
	if response == nil {
		logs.Logger.Errorf("CreateSwitch TencentCloud failed, response is nil, req[%v]", req)
		return cloud.CreateSwitchResponse{}, _errResponseIsNil
	}
	return cloud.CreateSwitchResponse{SwitchId: *response.Response.Subnet.SubnetId, RequestId: *response.Response.RequestId}, nil
}

func (p *TencentCloud) GetSwitch(req cloud.GetSwitchRequest) (cloud.GetSwitchResponse, error) {
	request := vpc.NewDescribeSubnetsRequest()
	//子网实例ID查询。形如：subnet-pxir56ns。每次请求的实例的上限为100。参数不支持同时指定SubnetIds和Filters
	request.SubnetIds = append(request.SubnetIds, &req.SwitchId)
	response, err := p.vpcClient.DescribeSubnets(request)
	if err != nil {
		logs.Logger.Errorf("GetSwitch TencentCloud failed.err: [%v], req[%v]", err, req)
		return cloud.GetSwitchResponse{}, err
	}
	if response == nil {
		logs.Logger.Errorf("GetSwitch TencentCloud failed, response is nil, req[%v]", req)
		return cloud.GetSwitchResponse{}, _errResponseIsNil
	}
	if *response.Response.TotalCount != 1 {
		logs.Logger.Errorf("GetSwitch TencentCloud failed, totalCount isn't one, req[%v]", req)
		return cloud.GetSwitchResponse{}, _errIsNotOne
	}
	subnet := response.Response.SubnetSet[0]
	res := cloud.GetSwitchResponse{Switch: subnetInfo2CloudSwitch(subnet)}
	return res, nil
}

func (p *TencentCloud) DescribeSwitches(req cloud.DescribeSwitchesRequest) (cloud.DescribeSwitchesResponse, error) {
	subnets := make([]*vpc.Subnet, 0, _pageSize)
	vpsIds := make([]*string, 0, 1)
	vpsIds = append(vpsIds, &req.VpcId)
	request := vpc.NewDescribeSubnetsRequest()
	filter := vpc.Filter{Name: common.StringPtr(_subnetFilterVpcId), Values: vpsIds}
	request.Filters = append(request.Filters, &filter)
	offset := _offset
	request.Limit = common.StringPtr(strconv.Itoa(_pageSize))
	for {
		request.Offset = common.StringPtr(strconv.Itoa(offset))
		response, err := p.vpcClient.DescribeSubnets(request)
		if err != nil {
			logs.Logger.Errorf("DescribeSwitches TencentCloud failed.err: [%v], req[%v]", err, req)
			return cloud.DescribeSwitchesResponse{}, err
		}
		if response == nil {
			logs.Logger.Errorf("DescribeSwitches TencentCloud failed, response is nil, req[%v]", req)
			return cloud.DescribeSwitchesResponse{}, _errResponseIsNil
		}
		subnets = append(subnets, response.Response.SubnetSet...)
		if len(response.Response.SubnetSet) < _pageSize {
			break
		}
		offset = (offset - 1) * _pageSize
	}
	switches := make([]cloud.Switch, 0, len(subnets))
	for _, subnet := range subnets {
		switches = append(switches, subnetInfo2CloudSwitch(subnet))
	}
	res := cloud.DescribeSwitchesResponse{Switches: switches}
	return res, nil
}

func subnetInfo2CloudSwitch(subnet *vpc.Subnet) cloud.Switch {
	var isDefault int
	if *subnet.IsDefault {
		isDefault = 1
	}
	return cloud.Switch{
		VpcId:                   *subnet.VpcId,
		SwitchId:                *subnet.SubnetId,
		Name:                    *subnet.SubnetName,
		IsDefault:               isDefault,
		AvailableIpAddressCount: int(*subnet.AvailableIpAddressCount),
		VStatus:                 cloud.SubnetAvailable,
		CreateAt:                *subnet.CreatedTime,
		ZoneId:                  *subnet.Zone,
		CidrBlock:               *subnet.CidrBlock,
		//GatewayIp:               "",
	}
}
