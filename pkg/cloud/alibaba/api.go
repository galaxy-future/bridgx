package alibaba

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	ecsClient "github.com/alibabacloud-go/ecs-20140526/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	vpcClient "github.com/alibabacloud-go/vpc-20160428/v2/client"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/galaxy-future/BridgX/internal/constants"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/galaxy-future/BridgX/pkg/utils"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
)

const (
	DirectionIn    = "ingress"
	DirectionOut   = "egress"
	Instancetype   = "InstanceType"
	AcceptLanguage = "zh-CN"
	PrePaid        = "PrePaid"
)

type AlibabaCloud struct {
	client    *ecs.Client
	vpcClient *vpcClient.Client
	ecsClient *ecsClient.Client
	bssClient *bssopenapi.Client
	lock      sync.Mutex
}

func (p *AlibabaCloud) GetInstancesByTags(region string, tags []cloud.Tag) (instances []cloud.Instance, err error) {
	request := ecs.CreateDescribeInstancesRequest()
	request.Scheme = "https"

	eTag := make([]ecs.DescribeInstancesTag, 0)
	for _, tag := range tags {
		eTag = append(eTag, ecs.DescribeInstancesTag{
			Key:   tag.Key,
			Value: tag.Value,
		})
	}
	request.Tag = &eTag
	pageNumber := 1
	request.PageSize = requests.NewInteger(50)
	cloudInstance := make([]ecs.Instance, 0)
	response, err := p.client.DescribeInstances(request)
	cloudInstance = append(cloudInstance, response.Instances.Instance...)
	maxPage := math.Ceil(float64(response.TotalCount / 50))
	for pageNumber < int(maxPage) {
		pageNumber++
		request.PageNumber = requests.NewInteger(pageNumber)
		response, err = p.client.DescribeInstances(request)
		cloudInstance = append(cloudInstance, response.Instances.Instance...)
	}
	instances = generateInstances(cloudInstance)
	return
}

func generateInstances(cloudInstance []ecs.Instance) (instances []cloud.Instance) {
	for _, instance := range cloudInstance {
		ipOuter := ""
		if len(instance.PublicIpAddress.IpAddress) > 0 {
			ipOuter = instance.PublicIpAddress.IpAddress[0]
		}
		instances = append(instances, cloud.Instance{
			Id:       instance.InstanceId,
			CostWay:  instance.InstanceChargeType,
			Provider: CloudName,
			IpInner:  strings.Join(instance.VpcAttributes.PrivateIpAddress.IpAddress, ","),
			IpOuter:  ipOuter,
			ImageId:  instance.ImageId,
			Network: &cloud.Network{
				VpcId:                   instance.VpcAttributes.VpcId,
				SubnetId:                instance.VpcAttributes.VSwitchId,
				SecurityGroup:           strings.Join(instance.SecurityGroupIds.SecurityGroupId, ","),
				InternetChargeType:      instance.InternetChargeType,
				InternetMaxBandwidthOut: instance.InternetMaxBandwidthOut,
			},
			Status: instance.Status,
		})
	}
	return
}

const (
	CloudName = "AlibabaCloud"
)

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
	return &AlibabaCloud{client: client, vpcClient: vpcClt, ecsClient: ecsClt, bssClient: bssCtl}, err
}

// BatchCreate the maximum of 'num' is 100
func (p *AlibabaCloud) BatchCreate(m cloud.Params, num int) (instanceIds []string, err error) {
	request := ecs.CreateRunInstancesRequest()
	request.Scheme = "https"

	request.RegionId = m.Region
	request.ImageId = m.ImageId
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
	if m.Charge.ChargeType == PrePaid {
		request.InstanceChargeType = m.Charge.ChargeType
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
	response, err := p.client.RunInstances(request)
	return response.InstanceIdSets.InstanceIdSet, err
}

func (p *AlibabaCloud) GetInstances(ids []string) (instances []cloud.Instance, err error) {
	batchIds := utils.StringSliceSplit(ids, 50)
	cloudInstance := make([]ecs.Instance, 0)
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

func (p *AlibabaCloud) BatchDelete(ids []string, regionId string) (err error) {
	request := ecs.CreateDeleteInstancesRequest()
	request.Scheme = "https"
	request.RegionId = regionId
	request.Force = requests.NewBoolean(true)
	batchIds := utils.StringSliceSplit(ids, 50)
	var response *ecs.DeleteInstancesResponse
	for _, onceIds := range batchIds {
		request.InstanceId = &onceIds
		response, err = p.client.DeleteInstances(request)
		logs.Logger.Infof("[BatchDelete] requestId: %s", response.RequestId)
	}
	return err
}

func (p *AlibabaCloud) StartInstance(id string) error {
	request := ecs.CreateStartInstanceRequest()
	request.Scheme = "https"

	request.InstanceId = id

	response, err := p.client.StartInstance(request)
	logs.Logger.Infof("[StartInstance] requestId: %s", response.RequestId)
	return err
}

func (p *AlibabaCloud) StopInstance(id string) error {
	request := ecs.CreateStopInstanceRequest()
	request.Scheme = "https"
	request.InstanceId = id

	response, err := p.client.StopInstance(request)
	logs.Logger.Infof("[StopInstance] requestId: %s", response.RequestId)
	return err
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
		switchIds := make([]string, 0, 64)
		if response.Body.VSwitchIds != nil {
			for _, switchId := range response.Body.VSwitchIds.VSwitchId {
				switchIds = append(switchIds, *switchId)
			}
		}
		res := cloud.GetVpcResponse{
			Vpc: cloud.VPC{
				VpcId:     *response.Body.VpcId,
				VpcName:   *response.Body.VpcName,
				CidrBlock: *response.Body.CidrBlock,
				Status:    *response.Body.Status,
				SwitchIds: switchIds,
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
					SwitchIds: tea.StringSliceValue(vpc.VSwitchIds.VSwitchId),
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
				VStatus:                 *response.Body.Status,
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
					VStatus:                 *vswitch.Status,
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
	request := &ecsClient.AuthorizeSecurityGroupRequest{
		RegionId:           tea.String(req.RegionId),
		SecurityGroupId:    tea.String(req.SecurityGroupId),
		IpProtocol:         tea.String(req.IpProtocol),
		PortRange:          tea.String(req.PortRange),
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
	request := &ecsClient.AuthorizeSecurityGroupEgressRequest{
		RegionId:         tea.String(req.RegionId),
		SecurityGroupId:  tea.String(req.SecurityGroupId),
		IpProtocol:       tea.String(req.IpProtocol),
		PortRange:        tea.String(req.PortRange),
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

func (p *AlibabaCloud) DescribeAvailableResource(req cloud.DescribeAvailableResourceRequest) (cloud.DescribeAvailableResourceResponse, error) {
	response, err := p.ecsClient.DescribeAvailableResource(&ecsClient.DescribeAvailableResourceRequest{
		RegionId:            tea.String(req.RegionId),
		ZoneId:              tea.String(req.ZoneId),
		DestinationResource: tea.String(Instancetype),
		NetworkCategory:     tea.String("vpc"),
	})
	if err != nil {
		logs.Logger.Errorf("DescribeAvailableResource AlibabaCloud failed.err: [%v] req[%v]", err, req)
		return cloud.DescribeAvailableResourceResponse{}, err
	}
	if response != nil && response.Body != nil && response.Body.AvailableZones != nil {
		zoneInsType := make(map[string][]cloud.InstanceType, 64)
		for _, zone := range response.Body.AvailableZones.AvailableZone {
			if zone.AvailableResources == nil {
				continue
			}
			resources := make([]cloud.InstanceType, 0, 100)
			for _, resource := range zone.AvailableResources.AvailableResource {
				if resource.SupportedResources == nil {
					continue
				}
				for _, ins := range resource.SupportedResources.SupportedResource {
					if ins != nil {
						resources = append(resources, cloud.InstanceType{
							Status:         *ins.Status,
							StatusCategory: *ins.StatusCategory,
							Value:          *ins.Value,
						})
					}
				}
			}
			zoneInsType[*zone.ZoneId] = resources
		}
		return cloud.DescribeAvailableResourceResponse{
			InstanceTypes: zoneInsType,
		}, nil
	}
	return cloud.DescribeAvailableResourceResponse{}, err
}

func (p *AlibabaCloud) DescribeInstanceTypes(req cloud.DescribeInstanceTypesRequest) (cloud.DescribeInstanceTypesResponse, error) {
	response, err := p.ecsClient.DescribeInstanceTypes(&ecsClient.DescribeInstanceTypesRequest{
		InstanceTypes: tea.StringSlice(req.TypeName),
	})
	if err != nil {
		logs.Logger.Errorf("DescribeInstanceTypes AlibabaCloud failed.err: [%v] req[%v]", err, req)
		return cloud.DescribeInstanceTypesResponse{}, err
	}
	if response != nil && response.Body != nil && response.Body.InstanceTypes != nil {
		InsTypeInfo := make([]cloud.InstanceInfo, 0, len(req.TypeName))
		for _, info := range response.Body.InstanceTypes.InstanceType {
			InsTypeInfo = append(InsTypeInfo, cloud.InstanceInfo{
				Core:        int(*info.CpuCoreCount),
				Memory:      int(*info.MemorySize),
				Family:      *info.InstanceTypeFamily,
				InsTypeName: *info.InstanceTypeId,
			})
		}
		return cloud.DescribeInstanceTypesResponse{
			Infos: InsTypeInfo,
		}, nil
	}
	return cloud.DescribeInstanceTypesResponse{}, err
}

func (p *AlibabaCloud) DescribeImages(req cloud.DescribeImagesRequest) (cloud.DescribeImagesResponse, error) {
	var page int32 = 1
	images := make([]cloud.Image, 0)
	for {
		response, err := p.ecsClient.DescribeImages(&ecsClient.DescribeImagesRequest{
			RegionId:   tea.String(req.RegionId),
			PageSize:   tea.Int32(50),
			PageNumber: tea.Int32(page),
		})
		if response != nil && response.Body != nil && response.Body.Images != nil {
			for _, img := range response.Body.Images.Image {
				images = append(images, cloud.Image{
					OsType:  *img.OSType,
					OsName:  *img.OSName,
					ImageId: *img.ImageId,
				})
			}

			if *response.Body.TotalCount > page*50 {
				page++
			} else {
				break
			}
		}
		if err != nil {
			logs.Logger.Errorf("DescribeImages failed,error: %v pageNumber:%d pageSize:%d region:%s", err, page, 50, req.RegionId)
		}
	}
	return cloud.DescribeImagesResponse{Images: images}, nil
}

func (*AlibabaCloud) ProviderType() string {
	return CloudName
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
			switch *rule.Direction {
			case DirectionIn:
				otherGroupId = *rule.SourceGroupId
				cidrIp = *rule.SourceCidrIp
				prefixListId = *rule.SourcePrefixListId
			case DirectionOut:
				otherGroupId = *rule.DestGroupId
				cidrIp = *rule.DestCidrIp
				prefixListId = *rule.DestPrefixListId
			}
			rules = append(rules, cloud.SecurityGroupRule{
				VpcId:           *response.Body.VpcId,
				SecurityGroupId: *response.Body.SecurityGroupId,
				PortRange:       *rule.PortRange,
				Protocol:        *rule.IpProtocol,
				Direction:       *rule.Direction,
				GroupId:         otherGroupId,
				CidrIp:          cidrIp,
				PrefixListId:    prefixListId,
				CreateAt:        *rule.CreateTime,
			})
		}

	}
	if err != nil {
		logs.Logger.Errorf("DescribeGroupRules failed,error: %v groupId:%s", err, req.SecurityGroupId)
	}
	return cloud.DescribeGroupRulesResponse{Rules: rules}, nil
}

var ChargeType = map[string]string{
	PostPaid:   constants.PostPaid,
	PayAsYouGo: constants.PayAsYouGo,
}
var PayStatus = map[string]int8{
	Paid:      constants.Paid,
	Unpaid:    constants.Unpaid,
	Cancelled: constants.Cancelled,
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

	orders := make([]cloud.Order, 0, orderNum*SubOrderNumPerMain)
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
			usageStartTime, _ := time.Parse("2006-01-02T15:04:05Z", subOrder.UsageStartTime)
			usageEndTime, _ := time.Parse("2006-01-02T15:04:05Z", subOrder.UsageEndTime)
			if subOrder.SubscriptionType == PayAsYouGo && usageEndTime.Sub(usageStartTime).Hours() > 24*365*20 {
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
				ChargeType:     ChargeType[subOrder.SubscriptionType],
				PayStatus:      PayStatus[subOrder.PaymentStatus],
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
