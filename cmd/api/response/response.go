package response

type ClusterCountResponse struct {
	ClusterNum int64 `json:"cluster_num"`
}

type InstanceCountResponse struct {
	InstanceNum int64 `json:"instance_num"`
}

type TaskCountResponse struct {
	TaskNum int64 `json:"task_num"`
}

type TaskListResponse struct {
	TaskList []TaskThumb `json:"task_list"`
	Pager    Pager       `json:"pager"`
}

type TaskThumb struct {
	TaskId      string `json:"task_id"`
	TaskName    string `json:"task_name"`
	TaskAction  string `json:"task_action"`
	Status      string `json:"status"`
	ClusterName string `json:"cluster_name"`
	CreateAt    string `json:"create_at"`
	ExecuteTime int    `json:"execute_time"`
	FinishAt    string `json:"finish_at"`
}

type ClusterThumb struct {
	ClusterId     string `json:"cluster_id"`
	ClusterName   string `json:"cluster_name"`
	InstanceCount int64  `json:"instance_count"`
	InstanceType  string `json:"instance_type"`
	ChargeType    string `json:"charge_type"`
	Provider      string `json:"provider"`
	Account       string `json:"account"`
	Usage         string `json:"usage"`
	CreateAt      string `json:"create_at"`
	CreateBy      string `json:"create_by"`
	UpdateBy      string `json:"update_by"`
	UpdateAt      string `json:"update_at"`
}

type TaskDetailResponse struct {
	TaskId              string `json:"task_id"`
	TaskName            string `json:"task_name"`
	ClusterName         string `json:"cluster_name"`
	TaskStatus          string `json:"task_status"`
	TaskResult          string `json:"task_result"`
	TaskAction          string `json:"task_action"`
	FailReason          string `json:"fail_reason"`
	RunNum              int    `json:"run_num"`
	SuspendNum          int    `json:"suspend_num"`
	SuccessNum          int    `json:"success_num"`
	FailNum             int    `json:"fail_num"`
	TotalNum            int    `json:"total_num"`
	SuccessRate         string `json:"success_rate"`
	ExecuteTime         int    `json:"execute_time"`
	BeforeInstanceCount int    `json:"before_instance_count"`
	AfterInstanceCount  int    `json:"after_instance_count"`
	ExpectInstanceCount int    `json:"expect_instance_count"`
	CreateAt            string `json:"create_at"`
}

type TaskDetailListResponse struct {
	TaskList []*TaskDetailResponse `json:"task_list"`
	Pager    Pager                 `json:"pager"`
}

type InstanceResponse struct {
	InstanceDetail
}

type InstanceDetail struct {
	InstanceId    string         `json:"instance_id"`
	Provider      string         `json:"provider"`
	RegionId      string         `json:"region_id"`
	ImageId       string         `json:"image_id"`
	InstanceType  string         `json:"instance_type"`
	IpInner       string         `json:"ip_inner"`
	IpOuter       string         `json:"ip_outer"`
	CreateAt      string         `json:"create_at"`
	StorageConfig *StorageConfig `json:"storage_config"`
	NetworkConfig *NetworkConfig `json:"network_config"`
}

type StorageConfig struct {
	SystemDiskType string     `json:"system_disk_type"`
	SystemDiskSize int        `json:"system_disk_size"`
	DataDisks      []DataDisk `json:"data_disks"`
	DataDiskNum    int        `json:"data_disk_num"`
}

type DataDisk struct {
	DataDiskType string `json:"data_disk_type"`
	DataDiskSize int    `json:"data_disk_size"`
}

type NetworkConfig struct {
	VpcName           string `json:"vpc_name"`
	SubnetIdName      string `json:"subnet_id_name"`
	SecurityGroupName string `json:"security_group_name"`
}

type InstanceListResponse struct {
	InstanceList []InstanceThumb `json:"instance_list"`
	Pager        Pager           `json:"pager"`
}

type InstanceThumb struct {
	InstanceId    string `json:"instance_id"`
	IpInner       string `json:"ip_inner"`
	IpOuter       string `json:"ip_outer"`
	Provider      string `json:"provider"`
	CreateAt      string `json:"create_at"`
	Status        string `json:"status"`
	StartupTime   int    `json:"startup_time"`
	ClusterName   string `json:"cluster_name"`
	InstanceType  string `json:"instance_type"`
	LoginName     string `json:"login_name"`
	LoginPassword string `json:"login_password"`
	ChargeType    string `json:"charge_type"`
}

type InstanceUsage struct {
	Id           string `json:"id"`
	ClusterName  string `json:"cluster_name"`
	InstanceId   string `json:"instance_id"`
	StartupAt    string `json:"startup_at"`
	ShutdownAt   string `json:"shutdown_at"`
	StartupTime  int    `json:"startup_time"`
	InstanceType string `json:"instance_type"`
}

type InstanceUsageResponse struct {
	InstanceList []InstanceUsage `json:"instance_list"`
	Pager        Pager           `json:"pager"`
}

type ListClustersResponse struct {
	ClusterList []ClusterThumb `json:"cluster_list"`
	Pager       Pager          `json:"pager"`
}

type ListClustersWithTagResponse struct {
	ClusterList []ClusterThumbWithTag `json:"cluster_list"`
	Pager       Pager                 `json:"pager"`
}

type ClusterTagsResponse struct {
	ClusterTags map[string]map[string]string `json:"cluster_tags"`
	Pager       Pager                        `json:"pager"`
}

type Pager struct {
	PageNumber int `json:"page_number"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
}

type ListCloudAccountResponse struct {
	CloudAccountList []CloudAccount `json:"account_list"`
	Pager            Pager          `json:"pager"`
}

type CloudAccount struct {
	Id          string `json:"id"`
	AccountName string `json:"account_name"`
	AccountKey  string `json:"account"`
	Provider    string `json:"provider"`
	CreateAt    string `json:"create_at"`
	CreateBy    string `json:"create_by"`
}

type EncryptCloudAccountInfo struct {
	AccountName          string `json:"account_name"`
	AccountKey           string `json:"account_key"`
	AccountSecretEncrypt string `json:"account_secret_encrypt"`
	Provider             string `json:"provider"`
}

type TaskInstancesResponse struct {
	InstanceList []InstanceThumb `json:"instance_list"`
	Pager        Pager           `json:"pager"`
}

type UserInfo struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
	UserType string `json:"user_type"`
	OrgId    int64  `json:"org_id"`
}

type UserThumb struct {
	UserId     string `json:"user_id"`
	UserName   string `json:"user_name"`
	CreateAt   string `json:"create_at"`
	CreateBy   string `json:"create_by"`
	UserStatus string `json:"user_status"`
}

type ListUsersResponse struct {
	UserList []UserThumb `json:"user_list"`
	Pager    Pager       `json:"pager"`
}

type OrgThumb struct {
	OrgId   string `json:"org_id"`
	OrgName string `json:"org_name"`
	UserNum string `json:"user_num"`
}

type ListOrgsResponse struct {
	OrgList []OrgThumb `json:"org_list"`
}

type InstanceStatResponse struct {
	InstanceTypeDesc string `json:"instance_type_desc"`
	InstanceCount    int64  `json:"instance_count"`
}

type ClusterThumbWithTag struct {
	ClusterId   string            `json:"cluster_id"`
	ClusterName string            `json:"cluster_name"`
	Provider    string            `json:"provider"`
	Tags        map[string]string `json:"tags"`
	CreateAt    string            `json:"create_at"`
	CreateBy    string            `json:"create_by"`
}
