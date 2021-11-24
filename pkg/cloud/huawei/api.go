package huawei

import (
	"fmt"

	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	ecs "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/model"
	region "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/region"
)

type HuaweiCloud struct {
	ecsClient *ecs.EcsClient
}

func New(AK, SK, regionId string) *HuaweiCloud {
	auth := basic.NewCredentialsBuilder().
		WithAk(AK).
		WithSk(SK).
		Build()

	client := ecs.NewEcsClient(
		ecs.EcsClientBuilder().
			WithRegion(region.ValueOf(regionId)).
			WithCredential(auth).
			Build())
	return &HuaweiCloud{ecsClient: client}
}

func (p *HuaweiCloud) GetInstances(ids []string) (instances []cloud.Instance, err error) {
	request := &model.ListServersDetailsRequest{}
	response, err := p.ecsClient.ListServersDetails(request)
	if err == nil {
		fmt.Printf("%+v\n", response)
	} else {
		fmt.Println(err)
	}
	return nil, err
}
