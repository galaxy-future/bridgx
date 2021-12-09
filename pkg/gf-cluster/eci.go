package gf_cluster

//cluster

type InstanceGroup struct {
	Id            int64  `json:"id"`
	KubernetesId  int64  `json:"kubernetes_id"`
	Name          string `json:"name"`
	Image         string `json:"image"`
	Cpu           string `json:"cpu"`
	Memory        string `json:"memory"`
	Disk          string `json:"disk"`
	InstanceCount int    `json:"instance_count"`
	CreatedUser   string `json:"created_user"`
	CreatedUserId int64  `json:"created_user_id"`
}

type InstanceGroupCreateRequest struct {
	KubernetesId  int64  `json:"kubernetes_id"`
	Name          string `json:"name"`
	Image         string `json:"image"`
	Cpu           string `json:"cpu"`
	Memory        string `json:"memory"`
	Disk          string `json:"disk"`
	InstanceCount int    `json:"instance_count"`
	CreatedUser   string `json:"created_user"`
	CreatedUserId int64  `json:"created_user_id"`
}

type InstanceGroupBatchCreateRequest struct {
	InstanceGroups []*InstanceGroupCreateRequest
}

type InstanceGroupBatchDeleteRequest struct {
	Ids []int64 `json:"ids"`
}

type InstanceGroupUpdateRequest struct {
	Id            int64  `json:"id"`
	KubernetesId  int64  `json:"kubernetes_id"`
	Name          string `json:"name"`
	Image         string `json:"image"`
	Cpu           string `json:"cpu"`
	Memory        string `json:"memory"`
	Disk          string `json:"disk"`
	InstanceCount int    `json:"instance_count"`
	CreatedUser   string `json:"created_user"`
	CreatedUserId int64  `json:"created_user_id"`
}

type InstanceGroupGetRequest struct {
	InstanceGroupID string `json:"instance_group_id"`
}

type InstanceGroupGetResponse struct {
	*ResponseBase
	InstanceGroup *InstanceGroup `json:"instance_group"`
}

func NewGetInstanceGroupResponse(cluster *InstanceGroup) InstanceGroupGetResponse {
	return InstanceGroupGetResponse{
		ResponseBase:  NewSuccessResponse(),
		InstanceGroup: cluster,
	}
}

type InstanceGroupListResponse struct {
	*ResponseBase
	InstanceGroups []*InstanceGroup `json:"instance_groups"`
	Pager          Pager            `json:"pager"`
}

func NewListInstanceGroupResponse(instanceGroups []*InstanceGroup, pager Pager) InstanceGroupListResponse {
	return InstanceGroupListResponse{
		ResponseBase:   NewSuccessResponse(),
		InstanceGroups: instanceGroups,
		Pager:          pager,
	}
}

// instances
type InstanceGroupExpandRequest struct {
	InstanceGroupId int64 `json:"Instance_group_id"`
	Count           int   `json:"count"`
}

type InstanceListRequest struct {
	InstanceGroupId string `json:"Instance_group_id"`
	Name            string `json:"name"`
}

type InstanceListResponse struct {
	*ResponseBase
	Instances []*Instance `json:"instances"`
}

func NewInstanceListResponse(instances []*Instance) *InstanceListResponse {
	return &InstanceListResponse{
		ResponseBase: NewSuccessResponse(),
		Instances:    instances,
	}
}

type InstanceGroupExpandOrShrinkRequest struct {
	InstanceGroupId int64 `json:"Instance_group_id"`
	Count           int   `json:"count"`
}

type InstanceRestartRequest struct {
	InstanceGroupId int64  `json:"Instance_group_id"`
	InstanceName    string `json:"instance_name"`
}

type InstanceDeleteRequest struct {
	InstanceGroupId int64  `json:"Instance_group_id"`
	InstanceName    string `json:"instance_name"`
}

type InstanceGroupShrinkRequest struct {
	InstanceGroupId int64 `json:"Instance_group_id"`
	Count           int   `json:"count"`
}

type InstanceGroupShrinkResponse struct {
	InstanceGroupId  int64    `json:"Instance_group_id"`
	InstanceNameList []string `json:"instance_name_list"`
}

type InstanceRemoveRequest struct {
	ClusterName string `json:"cluster_name"`
}

type Instance struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Ip      string `json:"ip"`
	HostIp  string `json:"host_ip"`
}