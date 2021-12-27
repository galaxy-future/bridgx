package tencent

import (
	"github.com/galaxy-future/BridgX/pkg/cloud"
)

func (p *TencentCloud) CreateSecurityGroup(req cloud.CreateSecurityGroupRequest) (cloud.CreateSecurityGroupResponse, error) {
	return cloud.CreateSecurityGroupResponse{}, nil
}

// AddIngressSecurityGroupRule 入参各云得统一
func (p *TencentCloud) AddIngressSecurityGroupRule(req cloud.AddSecurityGroupRuleRequest) error {
	return p.addSecGrpRule(req, cloud.SecGroupRuleIn)
}

func (p *TencentCloud) AddEgressSecurityGroupRule(req cloud.AddSecurityGroupRuleRequest) error {
	return p.addSecGrpRule(req, cloud.SecGroupRuleOut)
}

func (p *TencentCloud) DescribeSecurityGroups(req cloud.DescribeSecurityGroupsRequest) (cloud.DescribeSecurityGroupsResponse, error) {
	return cloud.DescribeSecurityGroupsResponse{}, nil
}

func (p *TencentCloud) DescribeGroupRules(req cloud.DescribeGroupRulesRequest) (cloud.DescribeGroupRulesResponse, error) {
	return cloud.DescribeGroupRulesResponse{}, nil
}

func (p *TencentCloud) addSecGrpRule(req cloud.AddSecurityGroupRuleRequest, direction string) error {
	return nil
}
