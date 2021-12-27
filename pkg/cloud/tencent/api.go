package tencent

import (
	"github.com/galaxy-future/BridgX/pkg/cloud"
	api "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/api/v20201106"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

type TencentCloud struct {
	vpcClient *vpc.Client
	cvmClient *cvm.Client
	apiClient *api.Client
}

func New(ak, sk, region string) (h *TencentCloud, err error) {
	credential := common.NewCredential(ak, sk)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = _vpcEndpoint
	vpcClient, err := vpc.NewClient(credential, region, cpf)
	if err != nil {
		return nil, err
	}

	cpf = profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = _cvmEndpoint
	cvmClient, err := cvm.NewClient(credential, region, cpf)
	if err != nil {
		return nil, err
	}

	cpf = profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = _apiEndpoint
	apiClient, err := api.NewClient(credential, region, cpf)
	if err != nil {
		return nil, err
	}

	return &TencentCloud{vpcClient: vpcClient, cvmClient: cvmClient, apiClient: apiClient}, nil
}

func (TencentCloud) ProviderType() string {
	return cloud.TencentCloud
}

// GetRegions 暂时返回中文名字
func (p *TencentCloud) GetRegions() (cloud.GetRegionsResponse, error) {
	request := api.NewDescribeRegionsRequest()
	request.Product = common.StringPtr("cvm")
	response, err := p.apiClient.DescribeRegions(request)
	if err != nil {
		return cloud.GetRegionsResponse{}, err
	}

	regions := make([]cloud.Region, 0, len(response.Response.RegionSet))
	for _, region := range response.Response.RegionSet {
		if *region.RegionState != "AVAILABLE" {
			continue
		}
		regions = append(regions, cloud.Region{
			RegionId:  *region.Region,
			LocalName: *region.RegionName,
		})
	}
	return cloud.GetRegionsResponse{Regions: regions}, nil
}

// GetZones zoneId zone
func (p *TencentCloud) GetZones(req cloud.GetZonesRequest) (cloud.GetZonesResponse, error) {
	request := api.NewDescribeZonesRequest()
	request.Product = common.StringPtr("cvm")
	response, err := p.apiClient.DescribeZones(request)
	if err != nil {
		return cloud.GetZonesResponse{}, err
	}

	zones := make([]cloud.Zone, 0, len(response.Response.ZoneSet))
	for _, zone := range response.Response.ZoneSet {
		if *zone.ZoneState != "AVAILABLE" {
			continue
		}
		zones = append(zones, cloud.Zone{
			ZoneId:    *zone.Zone,
			LocalName: *zone.ZoneName,
		})
	}
	return cloud.GetZonesResponse{Zones: zones}, nil
}

func (p *TencentCloud) GetOrders(req cloud.GetOrdersRequest) (cloud.GetOrdersResponse, error) {
	return cloud.GetOrdersResponse{}, nil
}
