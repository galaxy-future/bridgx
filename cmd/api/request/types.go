package request

import (
	"github.com/galaxy-future/BridgX/internal/service"
)

type AddTagRequest struct {
	ClusterName string            `json:"cluster_name"`
	Tags        map[string]string `json:"tags"`
}

type SetExpectInstanceCountRequest struct {
	ClusterName string `json:"cluster_name"`
	ExpectCount int    `json:"expect_count"`
}

type ExpandClusterRequest struct {
	TaskName    string `json:"task_name"`
	ClusterName string `json:"cluster_name"`
	Count       int    `json:"count"`
}

type ShrinkClusterRequest struct {
	TaskName    string   `json:"task_name"`
	ClusterName string   `json:"cluster_name"`
	IPs         []string `json:"ips"`
	Count       int      `json:"count"`
}

type CreateVpcRequest struct {
	Provider  string `json:"provider"`
	RegionId  string `json:"region_id"`
	VpcName   string `json:"vpc_name"`
	CidrBlock string `json:"cidr_block"`
	Ak        string `json:"ak"`
}

func (c *CreateVpcRequest) Check() bool {
	return c.Provider != "" && c.RegionId != "" && c.VpcName != "" && c.Ak != ""
}

type CreateSwitchRequest struct {
	SwitchName string `json:"switch_name"`
	RegionId   string `json:"region_id"`
	VpcId      string `json:"vpc_id"`
	CidrBlock  string `json:"cidr_block"`
	ZoneId     string `json:"zone_id"`
}

func (c *CreateSwitchRequest) Check() bool {
	return c.SwitchName != "" && c.VpcId != "" && c.CidrBlock != "" && c.ZoneId != ""
}

type CreateSecurityGroupRequest struct {
	VpcId             string `json:"vpc_id"`
	RegionId          string `json:"region_id"`
	SecurityGroupName string `json:"security_group_name"`
	SecurityGroupType string `json:"security_group_type"`
}

func (c *CreateSecurityGroupRequest) Check() bool {
	return c.SecurityGroupName != "" && c.RegionId != "" && c.VpcId != ""
}

type AddSecurityGroupRuleRequest struct {
	VpcId           string              `json:"vpc_id"`
	RegionId        string              `json:"region_id"`
	SecurityGroupId string              `json:"security_group_id"`
	Rules           []service.GroupRule `json:"rules"`
}

func (c *AddSecurityGroupRuleRequest) Check() bool {
	return c.VpcId != "" && c.RegionId != "" && c.SecurityGroupId != "" &&
		len(c.Rules) > 0
}

type CreateSecurityGroupWithRuleRequest struct {
	VpcId             string              `json:"vpc_id"`
	RegionId          string              `json:"region_id"`
	SecurityGroupName string              `json:"security_group_name"`
	SecurityGroupType string              `json:"security_group_type"`
	Rules             []service.GroupRule `json:"rules"`
}

func (c *CreateSecurityGroupWithRuleRequest) Check() bool {
	return c.SecurityGroupName != "" && c.RegionId != "" && c.VpcId != ""
}

type CreateNetworkRequest struct {
	Provider          string `json:"provider"`
	RegionId          string `json:"region_id"`
	CidrBlock         string `json:"cidr_block"`
	VpcName           string `json:"vpc_name"`
	ZoneId            string `json:"zone_id"`
	SwitchCidrBlock   string `json:"switch_cidr_block"`
	SwitchName        string `json:"switch_name"`
	SecurityGroupName string `json:"security_group_name"`
	SecurityGroupType string `json:"security_group_type"`
	Ak                string `json:"ak"`
}

func (c *CreateNetworkRequest) Check() bool {
	return c.Provider != "" && c.RegionId != "" && c.VpcName != "" && c.Ak != "" && c.SwitchName != "" &&
		c.SwitchCidrBlock != "" && c.ZoneId != "" && c.SecurityGroupName != "" && c.SecurityGroupType != ""
}

type LoginRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type CreateCloudAccountRequest struct {
	AccountName   string `json:"account_name"`
	Provider      string `json:"provider"`
	AccountKey    string `json:"account_key"`
	AccountSecret string `json:"account_secret"`
}

type EditCloudAccountRequest struct {
	AccountId   string `json:"account_id"`
	AccountName string `json:"account_name"`
	Provider    string `json:"provider"`
}

type EditOrgRequest struct {
	OrgId   int64  `json:"org_id"`
	OrgName string `json:"org_name"`
}

type CreateUserRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type ModifyAdminPasswordRequest struct {
	OldPassword string `json:"password"`
	NewPassword string `json:"new_password"`
}

type ModifyUsernameRequest struct {
	UserId      string `json:"user_id"`
	NewUsername string `json:"new_username"`
}

type UserStatusRequest struct {
	UserNames []string `json:"usernames"`
	Action    string   `json:"action"`
}

type CreateOrgRequest struct {
	OrgName  string `json:"org_name"`
	UserName string `json:"username"`
	Password string `json:"password"`
}
