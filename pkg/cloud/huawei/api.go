package huawei

import (
	"fmt"
	"net/http"
	"time"

	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/galaxy-future/BridgX/pkg/utils"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	bss "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/bss/v2"
	bssModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/bss/v2/model"
	bssRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/bss/v2/region"
	ecs "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2"
	ecsRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/region"
	iam "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3"
	iamModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	iamRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/region"
	ims "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ims/v2"
	imsModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ims/v2/model"
	imsRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ims/v2/region"
	vpc "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v2"
	vpcRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v2/region"
	secGrp "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3"
	secGrpRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3/region"
)

type HuaweiCloud struct {
	ecsClient    *ecs.EcsClient
	imsClient    *ims.ImsClient
	secGrpClient *secGrp.VpcClient
	vpcClient    *vpc.VpcClient
	iamClient    *iam.IamClient
	bssClient    *bss.BssClient
}

func New(ak, sk, regionId string) (h *HuaweiCloud, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%s,%s,%v", ak, regionId, e)
		}
	}()
	auth := basic.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		Build()

	ecsClt := ecs.NewEcsClient(
		ecs.EcsClientBuilder().
			WithRegion(ecsRegion.ValueOf(regionId)).
			WithCredential(auth).
			Build())
	imsClt := ims.NewImsClient(
		ims.ImsClientBuilder().
			WithRegion(imsRegion.ValueOf(regionId)).
			WithCredential(auth).
			Build())
	secGrpClt := secGrp.NewVpcClient(
		secGrp.VpcClientBuilder().
			WithRegion(secGrpRegion.ValueOf(regionId)).
			WithCredential(auth).
			Build())
	vpcClt := vpc.NewVpcClient(
		vpc.VpcClientBuilder().
			WithRegion(vpcRegion.ValueOf(regionId)).
			WithCredential(auth).
			Build())

	gAuth := global.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		Build()

	iamClt := iam.NewIamClient(
		iam.IamClientBuilder().
			WithRegion(iamRegion.ValueOf(regionId)).
			WithCredential(gAuth).
			Build())
	//bss region must be cn-north-1
	bssClt := bss.NewBssClient(
		bss.BssClientBuilder().
			WithRegion(bssRegion.ValueOf("cn-north-1")).
			WithCredential(gAuth).
			Build())
	return &HuaweiCloud{ecsClient: ecsClt, imsClient: imsClt, secGrpClient: secGrpClt, vpcClient: vpcClt,
		iamClient: iamClt, bssClient: bssClt}, nil
}

func (HuaweiCloud) ProviderType() string {
	return cloud.HuaweiCloud
}

// GetRegions 暂时返回中文名字
func (p *HuaweiCloud) GetRegions() (cloud.GetRegionsResponse, error) {
	request := &iamModel.KeystoneListRegionsRequest{}
	response, err := p.iamClient.KeystoneListRegions(request)
	if err != nil {
		return cloud.GetRegionsResponse{}, err
	}
	if response.HttpStatusCode != http.StatusOK {
		return cloud.GetRegionsResponse{}, fmt.Errorf("httpcode %d", response.HttpStatusCode)
	}

	regions := make([]cloud.Region, 0, len(*response.Regions))
	for _, region := range *response.Regions {
		regions = append(regions, cloud.Region{
			RegionId:  region.Id,
			LocalName: region.Locales.ZhCn,
		})
	}
	return cloud.GetRegionsResponse{Regions: regions}, nil
}

func (p *HuaweiCloud) DescribeImages(req cloud.DescribeImagesRequest) (cloud.DescribeImagesResponse, error) {
	pageSize := 500
	images := make([]cloud.Image, 0, pageSize)
	request := &imsModel.ListImagesRequest{}
	imageType := _imageType[req.ImageType]
	request.Imagetype = &imageType
	if req.ImageType == cloud.ImageGlobal {
		protected := true
		request.Protected = &protected
		if req.InsType != "" {
			request.FlavorId = &req.InsType
		}
	}
	statusRequest := imsModel.GetListImagesRequestStatusEnum().ACTIVE
	request.Status = &statusRequest
	limitRequest := int32(pageSize)
	request.Limit = &limitRequest
	markerRequest := ""
	for {
		if markerRequest != "" {
			request.Marker = &markerRequest
		}
		response, err := p.imsClient.ListImages(request)
		if err != nil {
			return cloud.DescribeImagesResponse{}, err
		}
		if response.HttpStatusCode != http.StatusOK {
			return cloud.DescribeImagesResponse{}, fmt.Errorf("httpcode %d", response.HttpStatusCode)
		}

		for _, img := range *response.Images {
			osType, _ := img.OsType.MarshalJSON()
			tmp, _ := img.Platform.MarshalJSON()
			platform := string(tmp)
			images = append(images, cloud.Image{
				Platform:  platform[1 : len(platform)-2],
				ImageId:   img.Id,
				OsType:    _osType[string(osType)],
				Size:      int(img.MinDisk),
				OsName:    img.Name,
				ImageName: img.Name,
			})
		}

		imgNum := len(*response.Images)
		if imgNum < pageSize {
			break
		}
		markerRequest = (*response.Images)[imgNum-1].Id
	}
	return cloud.DescribeImagesResponse{Images: images}, nil
}

func (p *HuaweiCloud) payOrders(orderId string) error {
	request := &bssModel.PayOrdersRequest{}
	request.Body = &bssModel.PayCustomerOrderReq{
		OrderId: orderId,
	}
	response, err := p.bssClient.PayOrders(request)
	if err != nil {
		return err
	}
	if response.HttpStatusCode != http.StatusNoContent {
		return fmt.Errorf("httpcode %d", response.HttpStatusCode)
	}
	return nil
}

//up to 50 at once
func (p *HuaweiCloud) listPrePaidResources(ids []string) (map[string]prePaidResources, error) {
	pageSize := 50
	batchIds := utils.StringSliceSplit(ids, int64(pageSize))
	resource := make(map[string]prePaidResources, len(ids))
	request := &bssModel.ListPayPerUseCustomerResourcesRequest{}
	limitQueryResourcesReq := int32(pageSize)
	onlyMainResourceQueryResourcesReq := int32(1)
	for _, onceIds := range batchIds {
		request.Body = &bssModel.QueryResourcesReq{
			Limit:            &limitQueryResourcesReq,
			OnlyMainResource: &onlyMainResourceQueryResourcesReq,
			ResourceIds:      &onceIds,
		}
		response, err := p.bssClient.ListPayPerUseCustomerResources(request)
		if err != nil {
			return nil, err
		}
		if response.HttpStatusCode != http.StatusOK {
			return nil, fmt.Errorf("httpcode %d", response.HttpStatusCode)
		}

		for _, res := range *response.Data {
			effTime, _ := time.Parse("2006-01-02T15:04:05Z", *res.EffectiveTime)
			expTime, _ := time.Parse("2006-01-02T15:04:05Z", *res.ExpireTime)
			resource[*res.ResourceId] = prePaidResources{
				Id:            *res.ResourceId,
				Name:          *res.ResourceName,
				RegionId:      *res.RegionCode,
				EffectiveTime: effTime,
				ExpireTime:    expTime,
				ExpirePolicy:  int(*res.ExpirePolicy),
				Status:        int(*res.Status),
			}
		}
	}
	return resource, nil
}

func (p *HuaweiCloud) GetOrders(req cloud.GetOrdersRequest) (cloud.GetOrdersResponse, error) {
	orders := make([]cloud.Order, 0, 0)
	return cloud.GetOrdersResponse{Orders: orders}, nil
}
