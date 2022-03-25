package alibaba

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	ecsClient "github.com/alibabacloud-go/ecs-20140526/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	vpcClient "github.com/alibabacloud-go/vpc-20160428/v2/client"
	sdkErr "github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/galaxy-future/BridgX/pkg/utils"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
)

const (
	Instancetype   = "InstanceType"
	AcceptLanguage = "zh-CN"
)

type AlibabaCloud struct {
	client    *ecs.Client
	vpcClient *vpcClient.Client
	ecsClient *ecsClient.Client
	bssClient *bssopenapi.Client
	lock      sync.Mutex
}

func New(AK, SK, region string) (*AlibabaCloud, error) {
	client, err := ecs.NewClientWithAccessKey(region, AK, SK)
	if err != nil {
		return nil, err
	}
	conf := openapi.Config{
		AccessKeyId:     tea.String(AK),
		AccessKeySecret: tea.String(SK),
		RegionId:        tea.String(region),
	}
	vpcClt, err := vpcClient.NewClient(&conf)
	if err != nil {
		return nil, err
	}
	ecsClt, err := ecsClient.NewClient(&conf)
	if err != nil {
		return nil, err
	}
	bssCtl, err := bssopenapi.NewClientWithAccessKey(region, AK, SK)
	if err != nil {
		return nil, err
	}
	return &AlibabaCloud{client: client, vpcClient: vpcClt, ecsClient: ecsClt, bssClient: bssCtl}, nil
}

// BatchCreate the maximum of 'num' is 100
func (p *AlibabaCloud) BatchCreate(m cloud.Params, num int) (instanceIds []string, err error) {
	request := ecs.CreateRunInstancesRequest()
	request.Scheme = "https"

	request.RegionId = m.Region
	request.ImageId = m.ImageId
	request.ZoneId = m.Zone
	request.InstanceType = m.InstanceType
	request.SecurityGroupId = m.Network.SecurityGroup
	request.VSwitchId = m.Network.SubnetId
	if m.Network.InternetMaxBandwidthOut != 0 {
		request.InternetChargeType = m.Network.InternetChargeType
		request.InternetMaxBandwidthOut = requests.NewInteger(m.Network.InternetMaxBandwidthOut)
	}
	request.Password = m.Password

	request.SystemDiskCategory = m.Disks.SystemDisk.Category
	request.SystemDiskSize = strconv.Itoa(m.Disks.SystemDisk.Size)
	dataDisks := make([]ecs.RunInstancesDataDisk, 0)
	for _, disk := range m.Disks.DataDisk {
		dataDisks = append(dataDisks, ecs.RunInstancesDataDisk{Size: strconv.Itoa(disk.Size), Category: disk.Category, PerformanceLevel: disk.PerformanceLevel})
	}
	request.DataDisk = &dataDisks
	request.Amount = requests.NewInteger(num)
	request.MinAmount = requests.NewInteger(num)
	if m.Charge.ChargeType == cloud.InstanceChargeTypePrePaid {
		request.InstanceChargeType = _inEcsChargeType[m.Charge.ChargeType]
		request.PeriodUnit = m.Charge.PeriodUnit
		request.Period = requests.NewInteger(m.Charge.Period)
	}
	if len(m.Tags) > 0 {
		tags := make([]ecs.RunInstancesTag, 0)
		for _, tag := range m.Tags {
			rTag := ecs.RunInstancesTag{
				Key:   tag.Key,
				Value: tag.Value,
			}
			tags = append(tags, rTag)
		}
		request.Tag = &tags
	}
	if m.DryRun {
		request.DryRun = "true"
	}

	response, err := p.client.RunInstances(request)
	if m.DryRun && err != nil {
		realErr := err.(*sdkErr.ServerError)
		if realErr.ErrorCode() == "DryRunOperation" {
			return []string{}, nil
		}
		return []string{}, err
	}
	if err != nil {
		return []string{}, err
	}
	return response.InstanceIdSets.InstanceIdSet, nil
}

func (p *AlibabaCloud) GetInstances(ids []string) (instances []cloud.Instance, err error) {
	batchIds := utils.StringSliceSplit(ids, 50)
	cloudInstance := make([]ecs.Instance, 0, len(ids))
	for _, onceIds := range batchIds {
		request := ecs.CreateDescribeInstancesRequest()
		request.Scheme = "https"
		var idsStr []byte
		var response *ecs.DescribeInstancesResponse
		idsStr, err = jsoniter.Marshal(onceIds)
		request.InstanceIds = string(idsStr)
		request.PageSize = requests.NewInteger(50)
		response, err = p.client.DescribeInstances(request)
		cloudInstance = append(cloudInstance, response.Instances.Instance...)
	}
	instances = generateInstances(cloudInstance)
	return
}

// BatchDelete 出现InvalidInstanceId.NotFound错误后，request不能复用，每次循环需重新创建
func (p *AlibabaCloud) BatchDelete(ids []string, regionId string) (err error) {
	batchIds := utils.StringSliceSplit(ids, 50)
	var response *ecs.DeleteInstancesResponse
	for _, onceIds := range batchIds {
		for {
			request := ecs.CreateDeleteInstancesRequest()
			request.Scheme = "https"
			request.RegionId = regionId
			request.Force = requests.NewBoolean(true)
			request.InstanceId = &onceIds
			response, err = p.client.DeleteInstances(request)
			if err == nil {
				logs.Logger.Infof("[BatchDelete] requestId: %s", response.RequestId)
				break
			}
			if realErr, ok := err.(*sdkErr.ServerError); ok {
				if realErr.ErrorCode() == "InvalidInstanceIds.NotFound" {
					break
				} else if realErr.ErrorCode() == "InvalidInstanceId.NotFound" {
					invalidIds := getInvalidIds(realErr.Message())
					onceIds = utils.StringSliceDiff(onceIds, invalidIds)
					if len(onceIds) > 0 {
						continue
					}
					break
				}
			}
			return err
		}
	}
	return nil
}

func (p *AlibabaCloud) StartInstances(ids []string) error {
	batchIds := utils.StringSliceSplit(ids, _maxNumEcsPerOperation)
	request := ecs.CreateStartInstancesRequest()
	request.Scheme = "https"
	for _, onceIds := range batchIds {
		request.InstanceId = &onceIds
		res, err := p.client.StartInstances(request)
		if err != nil {
			return err
		}
		logs.Logger.Debug(res)
	}
	return nil
}

func (p *AlibabaCloud) StopInstances(ids []string) error {
	batchIds := utils.StringSliceSplit(ids, _maxNumEcsPerOperation)
	request := ecs.CreateStopInstancesRequest()
	request.Scheme = "https"
	for _, onceIds := range batchIds {
		request.InstanceId = &onceIds
		res, err := p.client.StopInstances(request)
		if err != nil {
			return err
		}
		logs.Logger.Debug(res)
	}
	return nil
}

func (p *AlibabaCloud) GetInstancesByTags(region string, tags []cloud.Tag) (instances []cloud.Instance, err error) {
	eTag := make([]*ecsClient.ListTagResourcesRequestTag, 0, len(tags))
	for _, tag := range tags {
		eTag = append(eTag, &ecsClient.ListTagResourcesRequestTag{
			Key:   tea.String(tag.Key),
			Value: tea.String(tag.Value),
		})
	}
	request := &ecsClient.ListTagResourcesRequest{
		RegionId:     tea.String(region),
		ResourceType: tea.String("instance"),
		Tag:          eTag,
	}

	instanceIds := make([]string, 0, _pageSize)
	for {
		response, err := p.ecsClient.ListTagResources(request)
		if err != nil {
			return nil, err
		}

		for _, resource := range response.Body.TagResources.TagResource {
			instanceIds = append(instanceIds, tea.StringValue(resource.ResourceId))
		}
		nextToken := tea.StringValue(response.Body.NextToken)
		if nextToken == "" {
			break
		}
		request.NextToken = tea.String(nextToken)
	}

	return p.GetInstances(instanceIds)
}

func generateInstances(cloudInstance []ecs.Instance) (instances []cloud.Instance) {
	for _, instance := range cloudInstance {
		ipOuter := ""
		if len(instance.PublicIpAddress.IpAddress) > 0 {
			ipOuter = instance.PublicIpAddress.IpAddress[0]
		}
		expireAt, err := time.Parse("2006-01-02T15:04Z", instance.ExpiredTime)
		var expireAtPtr *time.Time
		if err == nil {
			expireAtPtr = &expireAt
		}
		instances = append(instances, cloud.Instance{
			Id:       instance.InstanceId,
			CostWay:  instance.InstanceChargeType,
			Provider: cloud.AlibabaCloud,
			IpInner:  strings.Join(instance.VpcAttributes.PrivateIpAddress.IpAddress, ","),
			IpOuter:  ipOuter,
			ImageId:  instance.ImageId,
			ExpireAt: expireAtPtr,
			Network: &cloud.Network{
				VpcId:                   instance.VpcAttributes.VpcId,
				SubnetId:                instance.VpcAttributes.VSwitchId,
				SecurityGroup:           strings.Join(instance.SecurityGroupIds.SecurityGroupId, ","),
				InternetChargeType:      _bandwidthChargeType[instance.InternetChargeType],
				InternetMaxBandwidthOut: instance.InternetMaxBandwidthOut,
			},
			Status: _ecsStatus[instance.Status],
		})
	}
	return
}

func (p *AlibabaCloud) GetInstancesByCluster(regionId, clusterName string) (instances []cloud.Instance, err error) {
	return p.GetInstancesByTags(regionId, []cloud.Tag{{
		Key:   cloud.ClusterName,
		Value: clusterName,
	}})
}

func (p *AlibabaCloud) CreateVPC(req cloud.CreateVpcRequest) (cloud.CreateVpcResponse, error) {
	request := &vpcClient.CreateVpcRequest{
		RegionId:  &req.RegionId,
		CidrBlock: &req.CidrBlock,
		VpcName:   &req.VpcName,
	}

	response, err := p.vpcClient.CreateVpc(request)
	if err != nil {
		logs.Logger.Errorf("CreateVPC AlibabaCloud failed.err: [%v], req[%v]", err, req)
		return cloud.CreateVpcResponse{}, err
	}
	if response != nil && response.Body != nil {
		return cloud.CreateVpcResponse{
			VpcId:     *response.Body.VpcId,
			RequestId: *response.Body.RequestId,
		}, nil
	}
	return cloud.CreateVpcResponse{}, nil
}

func (p *AlibabaCloud) GetVPC(req cloud.GetVpcRequest) (cloud.GetVpcResponse, error) {
	request := &vpcClient.DescribeVpcAttributeRequest{
		VpcId:    tea.String(req.VpcId),
		RegionId: tea.String(req.RegionId),
	}

	response, err := p.vpcClient.DescribeVpcAttribute(request)
	if err != nil {
		logs.Logger.Errorf("GetVPC AlibabaCloud failed.err: [%v], req[%v]", err, req)
		return cloud.GetVpcResponse{}, err
	}
	if response != nil && response.Body != nil {
		res := cloud.GetVpcResponse{
			Vpc: cloud.VPC{
				VpcId:     *response.Body.VpcId,
				VpcName:   *response.Body.VpcName,
				CidrBlock: *response.Body.CidrBlock,
				RegionId:  req.RegionId,
				Status:    _vpcStatus[*response.Body.Status],
				CreateAt:  *response.Body.CreationTime,
			},
		}
		return res, nil
	}

	return cloud.GetVpcResponse{}, err
}

func (p *AlibabaCloud) DescribeVpcs(req cloud.DescribeVpcsRequest) (cloud.DescribeVpcsResponse, error) {
	var page int32 = 1
	vpcs := make([]cloud.VPC, 0, 128)
	for {
		request := &vpcClient.DescribeVpcsRequest{
			RegionId:   tea.String(req.RegionId),
			PageSize:   tea.Int32(50),
			PageNumber: tea.Int32(page),
		}
		response, err := p.vpcClient.DescribeVpcs(request)
		if err != nil {
			logs.Logger.Errorf("DescribeVpcs AlibabaCloud failed.err: [%v], req[%v]", err, req)
			return cloud.DescribeVpcsResponse{}, err
		}
		if response != nil && response.Body != nil && response.Body.Vpcs != nil {
			for _, vpc := range response.Body.Vpcs.Vpc {
				vpcs = append(vpcs, cloud.VPC{
					VpcId:     *vpc.VpcId,
					VpcName:   *vpc.VpcName,
					CidrBlock: *vpc.CidrBlock,
					RegionId:  *vpc.RegionId,
					Status:    *vpc.Status,
					CreateAt:  *vpc.CreationTime,
				})
			}
			if *response.Body.TotalCount > page*50 {
				page++
			} else {
				break
			}
		}
		if err != nil {
			logs.Logger.Errorf("DescribeVpcs failed,error: %v pageNumber:%d pageSize:%d region:%s", err, page, 50, req.RegionId)
		}
	}
	return cloud.DescribeVpcsResponse{Vpcs: vpcs}, nil
}

func (p *AlibabaCloud) CreateSwitch(req cloud.CreateSwitchRequest) (cloud.CreateSwitchResponse, error) {
	request := &vpcClient.CreateVSwitchRequest{
		ZoneId:      tea.String(req.ZoneId),
		RegionId:    tea.String(req.RegionId),
		CidrBlock:   tea.String(req.CidrBlock),
		VpcId:       tea.String(req.VpcId),
		VSwitchName: tea.String(req.VSwitchName),
	}

	response, err := p.vpcClient.CreateVSwitch(request)
	if err != nil {
		logs.Logger.Errorf("CreateSwitch AlibabaCloud failed.err: [%v], req[%v]", err, req)
		return cloud.CreateSwitchResponse{}, err
	}
	if response != nil && response.Body != nil {
		return cloud.CreateSwitchResponse{
			SwitchId:  *response.Body.VSwitchId,
			RequestId: *response.Body.RequestId,
		}, err
	}
	return cloud.CreateSwitchResponse{}, err
}

func (p *AlibabaCloud) GetSwitch(req cloud.GetSwitchRequest) (cloud.GetSwitchResponse, error) {
	request := &vpcClient.DescribeVSwitchAttributesRequest{
		VSwitchId: tea.String(req.SwitchId),
	}
	response, err := p.vpcClient.DescribeVSwitchAttributes(request)
	if err != nil {
		logs.Logger.Errorf("GetSwitch AlibabaCloud failed.err: [%v], req[%v]", err, req)
		return cloud.GetSwitchResponse{}, err
	}
	if response != nil && response.Body != nil {
		var isDefault int
		if *response.Body.IsDefault {
			isDefault = 1
		}
		return cloud.GetSwitchResponse{
			Switch: cloud.Switch{
				VpcId:                   *response.Body.VpcId,
				SwitchId:                *response.Body.VSwitchId,
				Name:                    *response.Body.VSwitchName,
				IsDefault:               isDefault,
				AvailableIpAddressCount: int(*response.Body.AvailableIpAddressCount),
				VStatus:                 _subnetStatus[*response.Body.Status],
				CreateAt:                *response.Body.CreationTime,
				CidrBlock:               *response.Body.CidrBlock,
			},
		}, nil
	}
	return cloud.GetSwitchResponse{}, nil
}

func (p *AlibabaCloud) DescribeSwitches(req cloud.DescribeSwitchesRequest) (cloud.DescribeSwitchesResponse, error) {
	var page int32 = 1
	switches := make([]cloud.Switch, 0, 128)
	for {
		request := &vpcClient.DescribeVSwitchesRequest{
			VpcId:      tea.String(req.VpcId),
			PageSize:   tea.Int32(50),
			PageNumber: tea.Int32(page),
		}
		response, err := p.vpcClient.DescribeVSwitches(request)
		if err != nil {
			logs.Logger.Errorf("DescribeSwitches AlibabaCloud failed.err: [%v], req[%v]", err, req)
			return cloud.DescribeSwitchesResponse{}, err
		}
		if response != nil && response.Body != nil && response.Body.VSwitches != nil {
			for _, vswitch := range response.Body.VSwitches.VSwitch {
				var isDefault int
				if *vswitch.IsDefault {
					isDefault = 1
				}
				switches = append(switches, cloud.Switch{
					VpcId:                   *vswitch.VpcId,
					SwitchId:                *vswitch.VSwitchId,
					Name:                    *vswitch.VSwitchName,
					IsDefault:               isDefault,
					AvailableIpAddressCount: int(*vswitch.AvailableIpAddressCount),
					VStatus:                 _subnetStatus[*vswitch.Status],
					CreateAt:                *vswitch.CreationTime,
					CidrBlock:               *vswitch.CidrBlock,
					ZoneId:                  *vswitch.ZoneId,
				})
			}
			if *response.Body.TotalCount > page*50 {
				page++
			} else {
				break
			}
		}
		if err != nil {
			logs.Logger.Errorf("DescribeSwitches failed,error: %v pageNumber:%d pageSize:%d vpcId:%s", err, page, 50, req.VpcId)
		}
	}
	return cloud.DescribeSwitchesResponse{Switches: switches}, nil
}

func (p *AlibabaCloud) CreateSecurityGroup(req cloud.CreateSecurityGroupRequest) (cloud.CreateSecurityGroupResponse, error) {
	request := &ecsClient.CreateSecurityGroupRequest{
		RegionId:          tea.String(req.RegionId),
		SecurityGroupName: tea.String(req.SecurityGroupName),
		VpcId:             tea.String(req.VpcId),
		SecurityGroupType: tea.String(req.SecurityGroupType),
	}

	response, err := p.ecsClient.CreateSecurityGroup(request)
	if err != nil {
		logs.Logger.Errorf("CreateSecurityGroup AlibabaCloud failed.err: [%v], req[%v]", err, req)
		return cloud.CreateSecurityGroupResponse{}, err
	}
	if response != nil && response.Body != nil {
		return cloud.CreateSecurityGroupResponse{
			SecurityGroupId: *response.Body.SecurityGroupId,
			RequestId:       *response.Body.RequestId,
		}, nil
	}
	return cloud.CreateSecurityGroupResponse{}, err
}

func (p *AlibabaCloud) AddIngressSecurityGroupRule(req cloud.AddSecurityGroupRuleRequest) error {
	portRange := getPortRange(req.PortFrom, req.PortTo, req.IpProtocol)
	if req.GroupId == "" && req.CidrIp == "" {
		req.CidrIp = "0.0.0.0/0"
	}
	request := &ecsClient.AuthorizeSecurityGroupRequest{
		RegionId:           tea.String(req.RegionId),
		SecurityGroupId:    tea.String(req.SecurityGroupId),
		IpProtocol:         tea.String(_protocol[req.IpProtocol]),
		PortRange:          tea.String(portRange),
		SourceGroupId:      tea.String(req.GroupId),
		SourceCidrIp:       tea.String(req.CidrIp),
		SourcePrefixListId: tea.String(req.PrefixListId),
	}

	_, err := p.ecsClient.AuthorizeSecurityGroup(request)
	if err != nil {
		logs.Logger.Errorf("AddIngressSecurityGroupRule AlibabaCloud failed.err: [%v], req[%v]", err, req)
		return err
	}
	return nil
}

func (p *AlibabaCloud) AddEgressSecurityGroupRule(req cloud.AddSecurityGroupRuleRequest) error {
	portRange := getPortRange(req.PortFrom, req.PortTo, req.IpProtocol)
	if req.GroupId == "" && req.CidrIp == "" {
		req.CidrIp = "0.0.0.0/0"
	}
	request := &ecsClient.AuthorizeSecurityGroupEgressRequest{
		RegionId:         tea.String(req.RegionId),
		SecurityGroupId:  tea.String(req.SecurityGroupId),
		IpProtocol:       tea.String(_protocol[req.IpProtocol]),
		PortRange:        tea.String(portRange),
		DestGroupId:      tea.String(req.GroupId),
		DestCidrIp:       tea.String(req.CidrIp),
		DestPrefixListId: tea.String(req.PrefixListId),
	}

	_, err := p.ecsClient.AuthorizeSecurityGroupEgress(request)
	if err != nil {
		logs.Logger.Errorf("AddEgressSecurityGroupRule AlibabaCloud failed.err: [%v], req[%v]", err, req)
		return err
	}
	return nil
}

func (p *AlibabaCloud) DescribeSecurityGroups(req cloud.DescribeSecurityGroupsRequest) (cloud.DescribeSecurityGroupsResponse, error) {
	var page int32 = 1
	groups := make([]cloud.SecurityGroup, 0, 128)

	for {
		request := &ecsClient.DescribeSecurityGroupsRequest{
			RegionId:   tea.String(req.RegionId),
			VpcId:      tea.String(req.VpcId),
			PageSize:   tea.Int32(50),
			PageNumber: tea.Int32(page),
		}
		response, err := p.ecsClient.DescribeSecurityGroups(request)
		if err != nil {
			logs.Logger.Errorf("GetSecurityGroup AlibabaCloud failed.err: [%v], req[%v]", err, req)
			return cloud.DescribeSecurityGroupsResponse{}, err
		}
		if response != nil && response.Body != nil && response.Body.SecurityGroups != nil {
			for _, group := range response.Body.SecurityGroups.SecurityGroup {
				groups = append(groups, cloud.SecurityGroup{
					SecurityGroupId:   *group.SecurityGroupId,
					SecurityGroupType: *group.SecurityGroupType,
					SecurityGroupName: *group.SecurityGroupName,
					CreateAt:          *group.CreationTime,
					VpcId:             *group.VpcId,
					RegionId:          req.RegionId,
				})
			}
			if *response.Body.TotalCount > page*50 {
				page++
			} else {
				break
			}
		}
		if err != nil {
			logs.Logger.Errorf("GetSecurityGroup failed,error: %v pageNumber:%d pageSize:%d vpcId:%s", err, page, 50, req.VpcId)
		}
	}
	return cloud.DescribeSecurityGroupsResponse{Groups: groups}, nil
}

func (p *AlibabaCloud) GetRegions() (cloud.GetRegionsResponse, error) {
	response, err := p.vpcClient.DescribeRegions(&vpcClient.DescribeRegionsRequest{
		AcceptLanguage: tea.String(AcceptLanguage),
	})
	if err != nil {
		logs.Logger.Errorf("GetRegions AlibabaCloud failed.err: [%v]", err)
		return cloud.GetRegionsResponse{}, err
	}
	if response != nil && response.Body != nil {
		regions := make([]cloud.Region, 0, 100)
		for _, region := range response.Body.Regions.Region {
			regions = append(regions, cloud.Region{
				RegionId:  *region.RegionId,
				LocalName: *region.LocalName,
			})
		}
		return cloud.GetRegionsResponse{
			Regions: regions,
		}, nil
	}
	return cloud.GetRegionsResponse{}, nil
}

func (p *AlibabaCloud) GetZones(req cloud.GetZonesRequest) (cloud.GetZonesResponse, error) {
	response, err := p.vpcClient.DescribeZones(&vpcClient.DescribeZonesRequest{
		RegionId: tea.String(req.RegionId),
	})
	if err != nil {
		logs.Logger.Errorf("GetZones AlibabaCloud failed.err: [%v] req[%v]", err, req)
		return cloud.GetZonesResponse{}, err
	}
	if response != nil && response.Body != nil {
		zones := make([]cloud.Zone, 0, 100)
		for _, region := range response.Body.Zones.Zone {
			zones = append(zones, cloud.Zone{
				ZoneId:    *region.ZoneId,
				LocalName: *region.LocalName,
			})
		}
		return cloud.GetZonesResponse{
			Zones: zones,
		}, nil
	}
	return cloud.GetZonesResponse{}, err
}

// DescribeAvailableResource response miss InstanceChargeType
func (p *AlibabaCloud) DescribeAvailableResource(req cloud.DescribeAvailableResourceRequest) (cloud.DescribeAvailableResourceResponse, error) {
	request := &ecsClient.DescribeAvailableResourceRequest{
		RegionId:            tea.String(req.RegionId),
		DestinationResource: tea.String(Instancetype),
		NetworkCategory:     tea.String("vpc"),
	}
	if req.ZoneId != "" {
		request.ZoneId = tea.String(req.ZoneId)
	}

	zoneInsType := make(map[string][]cloud.InstanceType, 8)
	insTypeChargeType := []string{"PrePaid", "PostPaid"}
	for _, chargeType := range insTypeChargeType {
		request.InstanceChargeType = tea.String(chargeType)
		response, err := p.ecsClient.DescribeAvailableResource(request)
		if err != nil {
			logs.Logger.Errorf("DescribeAvailableResource AlibabaCloud failed.err: [%v] req[%v]", err, req)
			return cloud.DescribeAvailableResourceResponse{}, err
		}
		if response == nil || response.Body == nil || response.Body.AvailableZones == nil {
			return cloud.DescribeAvailableResourceResponse{}, errors.New("response is null")
		}

		for _, zone := range response.Body.AvailableZones.AvailableZone {
			if zone.AvailableResources == nil || len(zone.AvailableResources.AvailableResource) < 1 ||
				tea.StringValue(zone.StatusCategory) != "WithStock" {
				continue
			}

			zoneId := *zone.ZoneId
			insTypes, err := p.getResourceDetail(zone.AvailableResources.AvailableResource, chargeType)
			if err != nil {
				logs.Logger.Errorf("zoneId[%v] getResourceDetail failed: [%v]", zoneId, err)
				continue
			}

			zoneInsType[zoneId] = append(zoneInsType[zoneId], insTypes...)
		}
	}
	return cloud.DescribeAvailableResourceResponse{InstanceTypes: zoneInsType}, nil
}

// DescribeInstanceTypes Up to 10 at once
func (p *AlibabaCloud) DescribeInstanceTypes(req cloud.DescribeInstanceTypesRequest) (cloud.DescribeInstanceTypesResponse, error) {
	ecsInsTypes := make([]*ecsClient.DescribeInstanceTypesResponseBodyInstanceTypesInstanceType, 0, len(req.TypeName))
	var onceNum int64 = 10
	batchIds := utils.StringSliceSplit(req.TypeName, onceNum)
	for _, onceIds := range batchIds {
		request := &ecsClient.DescribeInstanceTypesRequest{
			InstanceTypes: tea.StringSlice(onceIds),
		}
		response, err := p.ecsClient.DescribeInstanceTypes(request)
		if err != nil {
			logs.Logger.Errorf("DescribeInstanceTypes AlibabaCloud failed.err: [%v] req[%v]", err, req)
			return cloud.DescribeInstanceTypesResponse{}, err
		}
		if response == nil || response.Body == nil || response.Body.InstanceTypes == nil {
			return cloud.DescribeInstanceTypesResponse{}, errors.New("response is null")
		}
		for _, info := range response.Body.InstanceTypes.InstanceType {
			ecsInsTypes = append(ecsInsTypes, info)
		}
	}
	insTypes := ecsInsType2CloudInsType(ecsInsTypes)
	return cloud.DescribeInstanceTypesResponse{Infos: insTypes}, nil
}

func (p *AlibabaCloud) DescribeImages(req cloud.DescribeImagesRequest) (cloud.DescribeImagesResponse, error) {
	var page int32 = 1
	var pageSize int32 = 50
	images := make([]cloud.Image, 0)
	request := &ecsClient.DescribeImagesRequest{
		RegionId:        tea.String(req.RegionId),
		PageSize:        tea.Int32(pageSize),
		ImageOwnerAlias: tea.String(_imageType[req.ImageType]),
	}
	if req.ImageType == cloud.ImageGlobal && req.InsType != "" {
		request.InstanceType = tea.String(req.InsType)
	}
	for {
		request.PageNumber = tea.Int32(page)
		response, err := p.ecsClient.DescribeImages(request)
		if err != nil {
			return cloud.DescribeImagesResponse{}, fmt.Errorf("pageNumber:%d pageSize:%d region:%s, %v", page, 50, req.RegionId, err)
		}
		if response != nil && response.Body != nil && response.Body.Images != nil {
			for _, img := range response.Body.Images.Image {
				images = append(images, cloud.Image{
					Platform:  *img.Platform,
					OsType:    _osType[*img.OSType],
					OsName:    *img.OSName,
					Size:      int(tea.Int32Value(img.Size)),
					ImageId:   *img.ImageId,
					ImageName: *img.ImageName,
				})
			}

			if page*pageSize >= tea.Int32Value(response.Body.TotalCount) {
				break
			}
			page++
		}
	}
	return cloud.DescribeImagesResponse{Images: images}, nil
}

func (*AlibabaCloud) ProviderType() string {
	return cloud.AlibabaCloud
}

func (p *AlibabaCloud) DescribeGroupRules(req cloud.DescribeGroupRulesRequest) (cloud.DescribeGroupRulesResponse, error) {
	rules := make([]cloud.SecurityGroupRule, 0, 128)
	request := &ecsClient.DescribeSecurityGroupAttributeRequest{
		RegionId:        tea.String(req.RegionId),
		SecurityGroupId: tea.String(req.SecurityGroupId),
	}
	response, err := p.ecsClient.DescribeSecurityGroupAttribute(request)
	if err != nil {
		logs.Logger.Errorf("DescribeGroupRules AlibabaCloud failed.err: [%v], req[%v]", err, req)
		return cloud.DescribeGroupRulesResponse{}, err
	}
	if response != nil && response.Body != nil && response.Body.Permissions != nil {
		for _, rule := range response.Body.Permissions.Permission {
			var otherGroupId, cidrIp, prefixListId string
			switch _secGrpRuleDirection[*rule.Direction] {
			case cloud.SecGroupRuleIn:
				otherGroupId = *rule.SourceGroupId
				cidrIp = *rule.SourceCidrIp
				if cidrIp == "" {
					cidrIp = tea.StringValue(rule.Ipv6SourceCidrIp)
				}
				prefixListId = *rule.SourcePrefixListId
			case cloud.SecGroupRuleOut:
				otherGroupId = *rule.DestGroupId
				cidrIp = *rule.DestCidrIp
				if cidrIp == "" {
					cidrIp = tea.StringValue(rule.Ipv6DestCidrIp)
				}
				prefixListId = *rule.DestPrefixListId
			}

			from, to := portRange2Int(*rule.PortRange)
			rules = append(rules, cloud.SecurityGroupRule{
				VpcId:           *response.Body.VpcId,
				SecurityGroupId: *response.Body.SecurityGroupId,
				PortFrom:        from,
				PortTo:          to,
				Protocol:        _outProtocol[*rule.IpProtocol],
				Direction:       _secGrpRuleDirection[*rule.Direction],
				GroupId:         otherGroupId,
				CidrIp:          cidrIp,
				PrefixListId:    prefixListId,
				CreateAt:        *rule.CreateTime,
			})
		}
	}

	return cloud.DescribeGroupRulesResponse{Rules: rules}, nil
}

func (p *AlibabaCloud) GetOrders(req cloud.GetOrdersRequest) (cloud.GetOrdersResponse, error) {
	request := bssopenapi.CreateQueryOrdersRequest()
	request.Scheme = "https"
	request.CreateTimeStart = req.StartTime.Format("2006-01-02T15:04:05Z")
	request.CreateTimeEnd = req.EndTime.Format("2006-01-02T15:04:05Z")
	request.PageNum = requests.NewInteger(req.PageNum)
	request.PageSize = requests.NewInteger(req.PageSize)
	response, err := p.bssClient.QueryOrders(request)
	if err != nil {
		return cloud.GetOrdersResponse{}, err
	}
	if !response.Success {
		return cloud.GetOrdersResponse{}, errors.New(response.Message)
	}
	orderNum := len(response.Data.OrderList.Order)
	if orderNum == 0 {
		return cloud.GetOrdersResponse{}, nil
	}

	orders := make([]cloud.Order, 0, orderNum*_subOrderNumPerMain)
	detailReq := bssopenapi.CreateGetOrderDetailRequest()
	detailReq.Scheme = "https"
	for _, row := range response.Data.OrderList.Order {
		detailReq.OrderId = row.OrderId
		detailRsp, err := p.bssClient.GetOrderDetail(detailReq)
		if err != nil {
			return cloud.GetOrdersResponse{}, err
		}
		if !detailRsp.Success {
			return cloud.GetOrdersResponse{}, errors.New(detailRsp.Message)
		}
		if len(detailRsp.Data.OrderList.Order) == 0 {
			continue
		}

		for _, subOrder := range detailRsp.Data.OrderList.Order {
			orderTime, _ := time.Parse("2006-01-02T15:04:05Z", subOrder.CreateTime)
			var usageStartTime, usageEndTime time.Time
			if subOrder.UsageStartTime == "" {
				usageStartTime = orderTime
			} else {
				usageStartTime, _ = time.Parse("2006-01-02T15:04:05Z", subOrder.UsageStartTime)
			}
			if subOrder.UsageEndTime == "" {
				usageEndTime = orderTime
			} else {
				usageEndTime, _ = time.Parse("2006-01-02T15:04:05Z", subOrder.UsageEndTime)
			}
			if _orderChargeType[subOrder.SubscriptionType] == cloud.OrderPostPaid && usageEndTime.Sub(usageStartTime).Hours() > 24*365*20 {
				usageEndTime, _ = time.Parse("2006-01-02 15:04:05", "2038-01-01 00:00:00")
			}

			orders = append(orders, cloud.Order{
				OrderId:        subOrder.SubOrderId,
				OrderTime:      orderTime,
				Product:        subOrder.ProductCode,
				Quantity:       cast.ToInt32(subOrder.Quantity),
				UsageStartTime: usageStartTime,
				UsageEndTime:   usageEndTime,
				RegionId:       subOrder.Region,
				ChargeType:     _orderChargeType[subOrder.SubscriptionType],
				PayStatus:      _payStatus[subOrder.PaymentStatus],
				Currency:       subOrder.Currency,
				Cost:           cast.ToFloat32(subOrder.PretaxAmount),
				Extend: map[string]interface{}{
					"main_order_id": subOrder.OrderId,
					"order_type":    subOrder.OrderType,
				},
			})
		}
	}
	return cloud.GetOrdersResponse{Orders: orders}, nil
}

//miss ChargeType,Status
func ecsInsType2CloudInsType(ecsInsType []*ecsClient.DescribeInstanceTypesResponseBodyInstanceTypesInstanceType) []cloud.InstanceType {
	insType := make([]cloud.InstanceType, 0, len(ecsInsType))
	for _, info := range ecsInsType {
		mem := int(tea.Float32Value(info.MemorySize))
		if mem < 1 {
			continue
		}
		isGpu := false
		if tea.Int32Value(info.GPUAmount) > 0 {
			isGpu = true
		}
		insType = append(insType, cloud.InstanceType{
			IsGpu:       isGpu,
			Core:        int(tea.Int32Value(info.CpuCoreCount)),
			Memory:      mem,
			Family:      tea.StringValue(info.InstanceTypeFamily),
			InsTypeName: tea.StringValue(info.InstanceTypeId),
		})
	}
	return insType
}

func (p *AlibabaCloud) getResourceDetail(availableResource []*ecsClient.DescribeAvailableResourceResponseBodyAvailableZonesAvailableZoneAvailableResourcesAvailableResource,
	chargeType string) ([]cloud.InstanceType, error) {
	insTypeStat := make(map[string]string, 100)
	insTypeIds := make([]string, 0, 100)
	for _, resource := range availableResource {
		if resource.SupportedResources == nil {
			continue
		}
		for _, ins := range resource.SupportedResources.SupportedResource {
			if ins == nil || _insTypeStat[tea.StringValue(ins.StatusCategory)] != cloud.InsTypeAvailable {
				continue
			}
			insTypeStat[*ins.Value] = _insTypeStat[tea.StringValue(ins.StatusCategory)]
			insTypeIds = append(insTypeIds, *ins.Value)
		}
	}

	res, err := p.DescribeInstanceTypes(cloud.DescribeInstanceTypesRequest{TypeName: insTypeIds})
	if err != nil {
		return nil, err
	}
	for i, info := range res.Infos {
		res.Infos[i].ChargeType = _insTypeChargeType[chargeType]
		res.Infos[i].Status = insTypeStat[info.InsTypeName]
	}
	return res.Infos, nil
}

func getPortRange(from, to int, protocol string) (portRange string) {
	if from < 1 || !(protocol == cloud.ProtocolUdp || protocol == cloud.ProtocolTcp) {
		return "-1/-1"
	}

	return fmt.Sprintf("%d/%d", from, to)
}

func portRange2Int(portRange string) (from, to int) {
	if portRange == "-1/-1" {
		return 0, 0
	}

	idx := strings.Index(portRange, "/")
	if idx == -1 {
		from = cast.ToInt(portRange)
		to = from
	} else {
		from = cast.ToInt(portRange[:idx])
		to = cast.ToInt(portRange[idx+1:])
	}
	return
}

func getInvalidIds(msg string) []string {
	invalidIds := make([]string, 0)
	msg = msg[strings.Index(msg, "(")+1 : strings.Index(msg, ")")]
	for {
		end := strings.Index(msg, ";")
		if end == -1 {
			invalidIds = append(invalidIds, msg)
			break
		}
		invalidIds = append(invalidIds, msg[:end])
		msg = msg[end+1:]
	}
	return invalidIds
}
