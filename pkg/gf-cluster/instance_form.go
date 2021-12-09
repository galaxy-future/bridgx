package gf_cluster

type InstanceForm struct {
	Id                   int64  `json:"id"`
	ExecuteStatus        string `json:"execute_status"`
	InstanceGroup        string `json:"instance_group"`
	Cpu                  string `json:"cpu"`
	Memory               string `json:"memory"`
	Disk                 string `json:"disk"`
	UpdatedInstanceCount int    `json:"updated_instance_count"`
	HostTime             int64  `json:"host_time"`
	OptType              string `json:"opt_type"`
	CreatedUserId        int64  `json:"created_user_id"`
	CreatedUserName      string `json:"created_user_name"`
	CreatedTime          int64  `json:"created_time"`
	ClusterName          string `json:"cluster_name"`
}

type InstanceFormListResponse struct {
	*ResponseBase
	InstanceForms []*InstanceForm `json:"instance_forms"`
	Pager         Pager           `json:"pager"`
}

func NewInstanceFormListResponse(instanceForms []*InstanceForm, pager Pager) *InstanceFormListResponse {
	return &InstanceFormListResponse{
		ResponseBase:  NewSuccessResponse(),
		InstanceForms: instanceForms,
		Pager:         pager,
	}
}
