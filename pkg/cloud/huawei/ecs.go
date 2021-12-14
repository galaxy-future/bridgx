package huawei

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/galaxy-future/BridgX/pkg/utils"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/model"
	"github.com/spf13/cast"
)

func (p *HuaweiCloud) BatchCreate(m cloud.Params, num int) ([]string, error) {
	request := &model.CreateServersRequest{}
	listNicsServer := []model.PrePaidServerNic{
		{
			SubnetId: m.Network.SubnetId,
		},
	}

	sizeRootVolumePrePaidServerRootVolume := int32(m.Disks.SystemDisk.Size)
	rootVolumeServer := &model.PrePaidServerRootVolume{
		Volumetype: _rootDiskCategory[m.Disks.SystemDisk.Category],
		Size:       &sizeRootVolumePrePaidServerRootVolume,
	}
	listDataVolumesServer := make([]model.PrePaidServerDataVolume, 0, len(m.Disks.DataDisk))
	for _, disk := range m.Disks.DataDisk {
		listDataVolumesServer = append(listDataVolumesServer, model.PrePaidServerDataVolume{
			Volumetype: _dataDiskCategory[disk.Category],
			Size:       int32(disk.Size),
		})
	}

	adminPassServerPrePaidServer := m.Password
	countServerPrePaidServer := int32(num)

	chargeType := _inEcsChargeType[m.Charge.ChargeType]
	periodType := _ecsPeriodType[m.Charge.PeriodUnit]
	periodNum := int32(m.Charge.Period)
	extendParam := &model.PrePaidServerExtendParam{
		ChargingMode: &chargeType,
		PeriodType:   &periodType,
		PeriodNum:    &periodNum,
	}
	listServerTagsServer := make([]model.PrePaidServerTag, 0, len(m.Tags))
	for _, tag := range m.Tags {
		listServerTagsServer = append(listServerTagsServer, model.PrePaidServerTag{
			Key:   tag.Key,
			Value: tag.Value,
		})
	}

	serverbody := &model.PrePaidServer{
		ImageRef:    m.ImageId,
		FlavorRef:   m.InstanceType,
		Name:        fmt.Sprintf("ins%v", time.Now().UnixNano()),
		AdminPass:   &adminPassServerPrePaidServer,
		Vpcid:       m.Network.VpcId,
		Nics:        listNicsServer,
		Count:       &countServerPrePaidServer,
		RootVolume:  rootVolumeServer,
		DataVolumes: &listDataVolumesServer,
		ServerTags:  &listServerTagsServer,
		Extendparam: extendParam,
	}
	if m.Network.InternetMaxBandwidthOut > 0 {
		sizeBandwith := int32(m.Network.InternetMaxBandwidthOut)
		chargemodeBandwidth := _bandwidthChargeMode[m.Network.InternetChargeType]
		bandwidthEip := &model.PrePaidServerEipBandwidth{
			Size:       &sizeBandwith,
			Sharetype:  model.GetPrePaidServerEipBandwidthSharetypeEnum().PER,
			Chargemode: &chargemodeBandwidth,
		}
		eipPublicip := &model.PrePaidServerEip{
			Iptype:    m.Network.InternetIpType,
			Bandwidth: bandwidthEip,
		}
		publicipServer := &model.PrePaidServerPublicip{
			Eip: eipPublicip,
		}
		serverbody.Publicip = publicipServer
	}
	if m.Network.SecurityGroup != "" {
		idSecurityGroup := m.Network.SecurityGroup
		var listSecurityGroupsServer = []model.PrePaidServerSecurityGroup{
			{
				Id: &idSecurityGroup,
			},
		}
		serverbody.SecurityGroups = &listSecurityGroupsServer
	}

	request.Body = &model.CreateServersRequestBody{
		Server: serverbody,
	}
	response, err := p.ecsClient.CreateServers(request)
	if err != nil {
		return []string{}, err
	}
	if response.HttpStatusCode != http.StatusOK {
		return []string{}, fmt.Errorf("httpcode %d, %v", response.HttpStatusCode, *response.JobId)
	}
	return *(response.ServerIds), nil
}

func (p *HuaweiCloud) GetInstances(ids []string) (instances []cloud.Instance, err error) {
	idNum := len(ids)
	if idNum < 1 {
		return []cloud.Instance{}, nil
	}
	request := &model.ShowServerRequest{}
	ecsInfos := make([]model.ServerDetail, 0, idNum)
	for _, id := range ids {
		if id == "" {
			continue
		}
		request.ServerId = id
		response, err := p.ecsClient.ShowServer(request)
		if err != nil {
			logs.Logger.Errorf("ShowServer failed, %s, %s", id, err.Error())
			continue
		}
		if response.HttpStatusCode != http.StatusOK {
			logs.Logger.Errorf("id %s, httpcode %d", id, response.HttpStatusCode)
			continue
		}
		ecsInfos = append(ecsInfos, *(response.Server))
	}
	return ecsInfo2CloudIns(ecsInfos), nil
}

func (p *HuaweiCloud) GetInstancesByTags(regionId string, tags []cloud.Tag) (instances []cloud.Instance, err error) {
	ecsInfos := make([]model.ServerDetail, 0, _pageSize)
	request := &model.ListServersDetailsRequest{}
	listTag := make([]string, 0, len(tags))
	for _, tag := range tags {
		listTag = append(listTag, tag.Key+"="+tag.Value)
	}
	tagsRequest := strings.Join(listTag, ",")
	request.Tags = &tagsRequest
	limitRequest := int32(_pageSize)
	request.Limit = &limitRequest
	pageNum := 1
	for {
		offsetRequest := int32(pageNum)
		request.Offset = &offsetRequest
		response, err := p.ecsClient.ListServersDetails(request)
		if err != nil {
			return nil, err
		}
		if response.HttpStatusCode != http.StatusOK {
			return nil, fmt.Errorf("httpcode %d", response.HttpStatusCode)
		}
		ecsInfos = append(ecsInfos, *(response.Servers)...)
		if int32(pageNum*_pageSize) >= *response.Count {
			break
		}
		pageNum++
	}
	return ecsInfo2CloudIns(ecsInfos), nil
}

func (p *HuaweiCloud) GetInstancesByCluster(regionId, clusterName string) (instances []cloud.Instance, err error) {
	return p.GetInstancesByTags(regionId, []cloud.Tag{{
		Key:   cloud.ClusterName,
		Value: clusterName,
	}})
}

// BatchDelete 华为云限制一次最多操作_maxNumEcsPerOperation台
func (p *HuaweiCloud) BatchDelete(ids []string, regionId string) error {
	batchIds := utils.StringSliceSplit(ids, _maxNumEcsPerOperation)
	request := &model.DeleteServersRequest{}
	for _, onceIds := range batchIds {
		listServerIds := make([]model.ServerId, 0, len(onceIds))
		for _, id := range onceIds {
			listServerIds = append(listServerIds, model.ServerId{
				Id: id,
			})
		}
		deleteVolumeDeleteServersRequestBody := true
		deletePublicipDeleteServersRequestBody := true
		request.Body = &model.DeleteServersRequestBody{
			Servers:        listServerIds,
			DeleteVolume:   &deleteVolumeDeleteServersRequestBody,
			DeletePublicip: &deletePublicipDeleteServersRequestBody,
		}
		response, err := p.ecsClient.DeleteServers(request)
		if err != nil {
			return err
		}
		if response.HttpStatusCode != http.StatusOK {
			return fmt.Errorf("httpcode %d, %v", response.HttpStatusCode, *response.JobId)
		}
	}
	return nil
}

func (p *HuaweiCloud) StartInstances(ids []string) error {
	batchIds := utils.StringSliceSplit(ids, _maxNumEcsPerOperation)
	request := &model.BatchStartServersRequest{}
	for _, onceIds := range batchIds {
		listServersOsStart := make([]model.ServerId, 0, len(onceIds))
		for _, id := range onceIds {
			listServersOsStart = append(listServersOsStart, model.ServerId{
				Id: id,
			})
		}
		osStartOpt := &model.BatchStartServersOption{
			Servers: listServersOsStart,
		}
		request.Body = &model.BatchStartServersRequestBody{
			OsStart: osStartOpt,
		}
		response, err := p.ecsClient.BatchStartServers(request)
		if err != nil {
			return err
		}
		if response.HttpStatusCode != http.StatusOK {
			return fmt.Errorf("httpcode %d, %v", response.HttpStatusCode, *response.JobId)
		}
	}
	return nil
}

func (p *HuaweiCloud) StopInstances(ids []string) error {
	batchIds := utils.StringSliceSplit(ids, _maxNumEcsPerOperation)
	request := &model.BatchStopServersRequest{}
	for _, onceIds := range batchIds {
		listServersOsStop := make([]model.ServerId, 0, len(onceIds))
		for _, id := range onceIds {
			listServersOsStop = append(listServersOsStop, model.ServerId{
				Id: id,
			})
		}
		osStopOpt := &model.BatchStopServersOption{
			Servers: listServersOsStop,
		}
		request.Body = &model.BatchStopServersRequestBody{
			OsStop: osStopOpt,
		}
		response, err := p.ecsClient.BatchStopServers(request)
		if err != nil {
			return err
		}
		if response.HttpStatusCode != http.StatusOK {
			return fmt.Errorf("httpcode %d, %v", response.HttpStatusCode, *response.JobId)
		}
	}
	return nil
}

// GetZones 华为云无ZoneId字段用ZoneName填充
func (p *HuaweiCloud) GetZones(req cloud.GetZonesRequest) (cloud.GetZonesResponse, error) {
	request := &model.NovaListAvailabilityZonesRequest{}
	response, err := p.ecsClient.NovaListAvailabilityZones(request)
	if err != nil {
		return cloud.GetZonesResponse{}, err
	}
	if response.HttpStatusCode != http.StatusOK {
		return cloud.GetZonesResponse{}, fmt.Errorf("httpcode %d", response.HttpStatusCode)
	}

	zones := make([]cloud.Zone, 0, len(*response.AvailabilityZoneInfo))
	for _, zone := range *response.AvailabilityZoneInfo {
		if !zone.ZoneState.Available {
			continue
		}
		zones = append(zones, cloud.Zone{
			ZoneId:    zone.ZoneName,
			LocalName: zone.ZoneName,
		})
	}
	return cloud.GetZonesResponse{Zones: zones}, nil
}

func (p *HuaweiCloud) DescribeAvailableResource(req cloud.DescribeAvailableResourceRequest) (cloud.DescribeAvailableResourceResponse, error) {
	zoneIds := make([]string, 0, 8)
	if req.ZoneId == "" {
		zones, err := p.GetZones(cloud.GetZonesRequest{})
		if err != nil {
			return cloud.DescribeAvailableResourceResponse{}, err
		}
		for _, zone := range zones.Zones {
			zoneIds = append(zoneIds, zone.ZoneId)
		}

	} else {
		zoneIds = append(zoneIds, req.ZoneId)
	}

	zoneInsType := make(map[string][]cloud.InstanceType, len(zoneIds))
	request := &model.ListFlavorsRequest{}
	for _, zoneId := range zoneIds {
		request.AvailabilityZone = &zoneId
		response, err := p.ecsClient.ListFlavors(request)
		if err != nil {
			return cloud.DescribeAvailableResourceResponse{}, err
		}
		if response.HttpStatusCode != http.StatusOK {
			return cloud.DescribeAvailableResourceResponse{}, fmt.Errorf("httpcode %d", response.HttpStatusCode)
		}

		insType := make([]cloud.InstanceType, 0, len(*response.Flavors))
		for _, flavor := range *response.Flavors {
			insType = append(insType, cloud.InstanceType{
				InstanceInfo: cloud.InstanceInfo{
					Core:        cast.ToInt(flavor.Vcpus),
					Memory:      cast.ToInt(flavor.Ram / 1024),
					Family:      *flavor.OsExtraSpecs.Ecsperformancetype,
					InsTypeName: flavor.Id,
				},
				Status: getFlavorStatus(flavor.OsExtraSpecs, zoneId),
			})
		}
		zoneInsType[zoneId] = insType
	}

	return cloud.DescribeAvailableResourceResponse{InstanceTypes: zoneInsType}, nil
}

//DescribeInstanceTypes NovaShowFlavor 华为云还没实现
func (p *HuaweiCloud) DescribeInstanceTypes(req cloud.DescribeInstanceTypesRequest) (cloud.DescribeInstanceTypesResponse, error) {
	return cloud.DescribeInstanceTypesResponse{}, nil
}

//缺少子网id,eip带宽相关信息. ListServerInterfaces 可以拿到子网id,ListPublicips 可以获取eip信息
func ecsInfo2CloudIns(ecsInfos []model.ServerDetail) []cloud.Instance {
	instances := make([]cloud.Instance, 0, len(ecsInfos))
	for _, info := range ecsInfos {
		var ipInner []string
		ipOut := ""
		for _, v := range info.Addresses {
			for _, row := range v {
				if *(row.OSEXTIPStype) == model.GetServerAddressOSEXTIPStypeEnum().FIXED {
					ipInner = append(ipInner, row.Addr)
				} else if *(row.OSEXTIPStype) == model.GetServerAddressOSEXTIPStypeEnum().FLOATING {
					ipOut = row.Addr
				}
			}
		}
		var securityGroup []string
		for _, row := range info.SecurityGroups {
			securityGroup = append(securityGroup, row.Id)
		}

		instances = append(instances, cloud.Instance{
			Id:       info.Id,
			CostWay:  _ecsChargeType[info.Metadata["charging_mode"]],
			Provider: cloud.HuaweiCloud,
			IpInner:  strings.Join(ipInner, ","),
			IpOuter:  ipOut,
			ImageId:  info.Image.Id,
			Network: &cloud.Network{
				VpcId:         info.Metadata["vpc_id"],
				SecurityGroup: strings.Join(securityGroup, ","),
			},
			Status: _ecsStatus[info.Status],
		})
	}
	return instances
}

func getFlavorStatus(flavor *model.FlavorExtraSpec, zoneId string) string {
	status := *flavor.Condoperationaz
	if status == "" {
		return _insTypeStat[*flavor.Condoperationstatus]
	}

	staStr := status
	for {
		begin := strings.Index(staStr, zoneId)
		if begin == -1 {
			return _insTypeStat[*flavor.Condoperationstatus]
		}
		staStr = staStr[begin:]
		zoneIdx := strings.Index(staStr, "(")
		if staStr[:zoneIdx] != zoneId {
			staStr = staStr[zoneIdx:]
			continue
		}

		end := strings.Index(staStr, ")")
		return _insTypeStat[staStr[zoneIdx+1:end]]
	}
}
