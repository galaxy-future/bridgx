package tencent

import (
	"strings"
	"time"

	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/galaxy-future/BridgX/pkg/utils"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

// BatchCreate AsVpcGateway
func (p *TencentCloud) BatchCreate(m cloud.Params, num int) ([]string, error) {
	request := cvm.NewRunInstancesRequest()
	request.InstanceChargeType = common.StringPtr(_inEcsChargeType[m.Charge.ChargeType])
	if m.Charge.ChargeType == cloud.InstanceChargeTypePrePaid {
		request.InstanceChargePrepaid = &cvm.InstanceChargePrepaid{
			Period:    common.Int64Ptr(int64(m.Charge.Period)),
			RenewFlag: common.StringPtr("NOTIFY_AND_MANUAL_RENEW"),
		}
	}

	request.Placement = &cvm.Placement{
		Zone: common.StringPtr(m.Zone),
	}
	request.InstanceType = common.StringPtr(m.InstanceType)
	request.ImageId = common.StringPtr(m.ImageId)
	request.SystemDisk = &cvm.SystemDisk{
		DiskType: common.StringPtr(m.Disks.SystemDisk.Category),
		DiskSize: common.Int64Ptr(int64(m.Disks.SystemDisk.Size)),
	}
	for _, disk := range m.Disks.DataDisk {
		request.DataDisks = append(request.DataDisks, &cvm.DataDisk{
			DiskType:           common.StringPtr(disk.Category),
			DiskSize:           common.Int64Ptr(int64(disk.Size)),
			DeleteWithInstance: common.BoolPtr(true),
		})
	}
	request.VirtualPrivateCloud = &cvm.VirtualPrivateCloud{
		VpcId:        common.StringPtr(m.Network.VpcId),
		SubnetId:     common.StringPtr(m.Network.SubnetId),
		AsVpcGateway: common.BoolPtr(false),
	}
	request.SecurityGroupIds = common.StringPtrs([]string{m.Network.SecurityGroup})
	if m.Network.InternetMaxBandwidthOut > 0 {
		request.InternetAccessible = &cvm.InternetAccessible{
			InternetChargeType:      common.StringPtr(_bandwidthChargeMode[m.Network.InternetChargeType]),
			InternetMaxBandwidthOut: common.Int64Ptr(int64(m.Network.InternetMaxBandwidthOut)),
		}
	}
	request.InstanceCount = common.Int64Ptr(int64(num))
	request.LoginSettings = &cvm.LoginSettings{
		Password: common.StringPtr(m.Password),
	}
	if len(m.Tags) > 0 {
		request.TagSpecification = []*cvm.TagSpecification{
			{
				ResourceType: common.StringPtr("instance"),
			},
		}
		for _, tag := range m.Tags {
			request.TagSpecification[0].Tags = append(request.TagSpecification[0].Tags, &cvm.Tag{
				Key:   common.StringPtr(tag.Key),
				Value: common.StringPtr(tag.Value),
			})
		}
	}
	request.DryRun = common.BoolPtr(m.DryRun)

	response, err := p.cvmClient.RunInstances(request)
	if err != nil {
		return nil, err
	}
	return common.StringValues(response.Response.InstanceIdSet), nil
}

func (p *TencentCloud) GetInstances(ids []string) (instances []cloud.Instance, err error) {
	idNum := len(ids)
	if idNum < 1 {
		return []cloud.Instance{}, nil
	}
	batchIds := utils.StringSliceSplit(ids, _maxNumEcsPerOperation)
	cvmInstances := make([]*cvm.Instance, 0, idNum)
	for _, onceIds := range batchIds {
		request := cvm.NewDescribeInstancesRequest()
		request.InstanceIds = common.StringPtrs(onceIds)
		response, err := p.cvmClient.DescribeInstances(request)
		if err != nil {
			return nil, err
		}
		cvmInstances = append(cvmInstances, response.Response.InstanceSet...)
	}
	return cvmIns2CloudIns(cvmInstances), nil
}

func (p *TencentCloud) GetInstancesByTags(regionId string, tags []cloud.Tag) (instances []cloud.Instance, err error) {
	request := cvm.NewDescribeInstancesRequest()
	request.Filters = make([]*cvm.Filter, 0, len(tags))
	for _, tag := range tags {
		request.Filters = append(request.Filters, &cvm.Filter{
			Name:   common.StringPtr("tag:" + tag.Key),
			Values: common.StringPtrs([]string{tag.Value}),
		})
	}
	request.Limit = common.Int64Ptr(_pageSize)
	cvmInstances := make([]*cvm.Instance, 0, _pageSize)
	var offset int64 = 0
	for {
		request.Offset = common.Int64Ptr(offset)
		response, err := p.cvmClient.DescribeInstances(request)
		if err != nil {
			return nil, err
		}
		cvmInstances = append(cvmInstances, response.Response.InstanceSet...)

		if offset+_pageSize > *response.Response.TotalCount {
			break
		}
		offset += _pageSize
	}
	return cvmIns2CloudIns(cvmInstances), nil
}

func (p *TencentCloud) GetInstancesByCluster(regionId, clusterName string) (instances []cloud.Instance, err error) {
	return p.GetInstancesByTags(regionId, []cloud.Tag{{
		Key:   cloud.ClusterName,
		Value: clusterName,
	}})
}

// BatchDelete 非后付费机器，子机销毁时，无法自动删除数据盘
func (p *TencentCloud) BatchDelete(ids []string, regionId string) error {
	batchIds := utils.StringSliceSplit(ids, _maxNumEcsPerOperation)
	for _, onceIds := range batchIds {
		request := cvm.NewTerminateInstancesRequest()
		request.InstanceIds = common.StringPtrs(onceIds)
		_, err := p.cvmClient.TerminateInstances(request)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *TencentCloud) StartInstances(ids []string) error {
	batchIds := utils.StringSliceSplit(ids, _maxNumEcsPerOperation)
	request := cvm.NewStartInstancesRequest()
	for _, onceIds := range batchIds {
		request.InstanceIds = common.StringPtrs(onceIds)
		_, err := p.cvmClient.StartInstances(request)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *TencentCloud) StopInstances(ids []string) error {
	batchIds := utils.StringSliceSplit(ids, _maxNumEcsPerOperation)
	request := cvm.NewStopInstancesRequest()
	for _, onceIds := range batchIds {
		request.InstanceIds = common.StringPtrs(onceIds)
		request.StopType = common.StringPtr("SOFT_FIRST")
		request.StoppedMode = common.StringPtr("STOP_CHARGING")
		_, err := p.cvmClient.StopInstances(request)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *TencentCloud) DescribeAvailableResource(req cloud.DescribeAvailableResourceRequest) (cloud.DescribeAvailableResourceResponse, error) {
	request := cvm.NewDescribeZoneInstanceConfigInfosRequest()
	if req.ZoneId != "" {
		request.Filters = []*cvm.Filter{
			{
				Name:   common.StringPtr("zone"),
				Values: common.StringPtrs([]string{req.ZoneId}),
			},
		}
	}
	response, err := p.cvmClient.DescribeZoneInstanceConfigInfos(request)
	if err != nil {
		return cloud.DescribeAvailableResourceResponse{}, err
	}

	zoneInsType := make(map[string][]cloud.InstanceType, 8)
	zoneCvmInsType := make(map[string][]*cvm.InstanceTypeQuotaItem, 8)
	for _, insType := range response.Response.InstanceTypeQuotaSet {
		_, ok := zoneCvmInsType[*insType.Zone]
		if !ok {
			zoneCvmInsType[*insType.Zone] = make([]*cvm.InstanceTypeQuotaItem, 0, _pageSize)
		}

		zoneCvmInsType[*insType.Zone] = append(zoneCvmInsType[*insType.Zone], insType)
	}
	for k, v := range zoneCvmInsType {
		zoneInsType[k] = cvmInsType2CloudInsType(v)
	}
	return cloud.DescribeAvailableResourceResponse{InstanceTypes: zoneInsType}, nil
}

func (p *TencentCloud) DescribeInstanceTypes(req cloud.DescribeInstanceTypesRequest) (cloud.DescribeInstanceTypesResponse, error) {
	batchIds := utils.StringSliceSplit(req.TypeName, _maxNumEcsPerOperation)
	request := cvm.NewDescribeZoneInstanceConfigInfosRequest()
	cvmInsType := make([]*cvm.InstanceTypeQuotaItem, 0, len(req.TypeName))
	for _, onceIds := range batchIds {
		request.Filters = []*cvm.Filter{
			{
				Name:   common.StringPtr("instance-type"),
				Values: common.StringPtrs(onceIds),
			},
		}
		response, err := p.cvmClient.DescribeZoneInstanceConfigInfos(request)
		if err != nil {
			return cloud.DescribeInstanceTypesResponse{}, err
		}

		for _, insType := range response.Response.InstanceTypeQuotaSet {
			cvmInsType = append(cvmInsType, insType)
		}
	}
	insType := cvmInsType2CloudInsType(cvmInsType)
	return cloud.DescribeInstanceTypesResponse{Infos: insType}, nil
}

// DescribeImages miss OsType
func (p *TencentCloud) DescribeImages(req cloud.DescribeImagesRequest) (cloud.DescribeImagesResponse, error) {
	request := cvm.NewDescribeImagesRequest()
	request.Filters = []*cvm.Filter{
		{
			Name:   common.StringPtr("image-type"),
			Values: common.StringPtrs([]string{_imageType[req.ImageType]}),
		},
	}
	request.Limit = common.Uint64Ptr(uint64(_pageSize))
	if req.ImageType == cloud.ImageGlobal && req.InsType != "" {
		request.InstanceType = common.StringPtr(req.InsType)
	}

	images := make([]cloud.Image, 0, _pageSize)
	var offset uint64 = 0
	for {
		request.Offset = common.Uint64Ptr(offset)
		response, err := p.cvmClient.DescribeImages(request)
		if err != nil {
			return cloud.DescribeImagesResponse{}, err
		}

		for _, img := range response.Response.ImageSet {
			images = append(images, cloud.Image{
				Platform: *img.Platform,
				OsName:   *img.OsName,
				Size:     int(utils.Int64Value(img.ImageSize)),
				ImageId:  *img.ImageId,
			})
		}
		if offset+_pageSize > uint64(*response.Response.TotalCount) {
			break
		}
		offset += _pageSize
	}
	return cloud.DescribeImagesResponse{Images: images}, nil
}

func cvmIns2CloudIns(cvmInstances []*cvm.Instance) []cloud.Instance {
	instances := make([]cloud.Instance, 0, len(cvmInstances))
	for _, info := range cvmInstances {
		ipInner := common.StringValues(info.PrivateIpAddresses)
		ipOut := ""
		if len(info.PublicIpAddresses) > 0 {
			ipOut = *info.PublicIpAddresses[0]
		}
		securityGroup := common.StringValues(info.SecurityGroupIds)
		bandwidthChargeType := ""
		bandwidthOut := 0
		if info.InternetAccessible.InternetChargeType != nil {
			bandwidthChargeType = _bandwidthChargeType[*info.InternetAccessible.InternetChargeType]
			bandwidthOut = int(*info.InternetAccessible.InternetMaxBandwidthOut)
		}
		var expireAt *time.Time
		if info.ExpiredTime != nil {
			expireTime, _ := time.Parse("2006-01-02T15:04:05Z", *info.ExpiredTime)
			expireAt = &expireTime
		}

		instances = append(instances, cloud.Instance{
			Id:       *info.InstanceId,
			CostWay:  _ecsChargeType[*info.InstanceChargeType],
			Provider: cloud.TencentCloud,
			IpInner:  strings.Join(ipInner, ","),
			IpOuter:  ipOut,
			ImageId:  *info.ImageId,
			Network: &cloud.Network{
				VpcId:                   *info.VirtualPrivateCloud.VpcId,
				SubnetId:                *info.VirtualPrivateCloud.SubnetId,
				SecurityGroup:           strings.Join(securityGroup, ","),
				InternetChargeType:      bandwidthChargeType,
				InternetMaxBandwidthOut: bandwidthOut,
			},
			Status:   _ecsStatus[utils.StringValue(info.InstanceState)],
			ExpireAt: expireAt,
		})
	}
	return instances
}

func cvmInsType2CloudInsType(cvmInsType []*cvm.InstanceTypeQuotaItem) []cloud.InstanceType {
	insType := make([]cloud.InstanceType, 0, len(cvmInsType))
	for _, info := range cvmInsType {
		stat := _insTypeStat[utils.StringValue(info.Status)]
		chargeType := _insTypeChargeType[utils.StringValue(info.InstanceChargeType)]
		if stat != cloud.InsTypeAvailable || chargeType == "" {
			continue
		}
		isGpu := false
		if utils.Int64Value(info.Gpu) > 0 {
			isGpu = true
		}
		insType = append(insType, cloud.InstanceType{
			ChargeType:  chargeType,
			IsGpu:       isGpu,
			Core:        int(utils.Int64Value(info.Cpu)),
			Memory:      int(utils.Int64Value(info.Memory)),
			Family:      utils.StringValue(info.InstanceFamily),
			InsTypeName: utils.StringValue(info.InstanceType),
			Status:      stat,
		})
	}
	return insType
}
