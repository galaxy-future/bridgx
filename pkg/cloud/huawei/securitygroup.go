package huawei

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/galaxy-future/BridgX/pkg/utils"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3/model"
	"github.com/spf13/cast"
)

// CreateSecurityGroup 将VpcId写入Description，方便查找
func (p *HuaweiCloud) CreateSecurityGroup(req cloud.CreateSecurityGroupRequest) (cloud.CreateSecurityGroupResponse, error) {
	request := &model.CreateSecurityGroupRequest{}
	securityGroupOpt := &model.CreateSecurityGroupOption{
		Name: req.SecurityGroupName,
	}
	request.Body = &model.CreateSecurityGroupRequestBody{
		SecurityGroup: securityGroupOpt,
	}
	response, err := p.secGrpClient.CreateSecurityGroup(request)
	if err != nil {
		return cloud.CreateSecurityGroupResponse{}, err
	}
	if response.HttpStatusCode != http.StatusCreated {
		return cloud.CreateSecurityGroupResponse{}, fmt.Errorf("httpcode %d", response.HttpStatusCode)
	}

	return cloud.CreateSecurityGroupResponse{SecurityGroupId: response.SecurityGroup.Id,
		RequestId: *response.RequestId}, nil
}

// AddIngressSecurityGroupRule 入参各云得统一
func (p *HuaweiCloud) AddIngressSecurityGroupRule(req cloud.AddSecurityGroupRuleRequest) error {
	return p.addSecGrpRule(req, cloud.SecGroupRuleIn)
}

func (p *HuaweiCloud) AddEgressSecurityGroupRule(req cloud.AddSecurityGroupRuleRequest) error {
	return p.addSecGrpRule(req, cloud.SecGroupRuleOut)
}

func (p *HuaweiCloud) DescribeSecurityGroups(req cloud.DescribeSecurityGroupsRequest) (cloud.DescribeSecurityGroupsResponse, error) {
	groups := make([]cloud.SecurityGroup, 0, _pageSize)
	request := &model.ListSecurityGroupsRequest{}
	limitRequest := int32(_pageSize)
	request.Limit = &limitRequest
	markerRequest := ""
	for {
		if markerRequest != "" {
			request.Marker = &markerRequest
		}
		response, err := p.secGrpClient.ListSecurityGroups(request)
		if err != nil {
			return cloud.DescribeSecurityGroupsResponse{}, err
		}
		if response.HttpStatusCode != http.StatusOK {
			return cloud.DescribeSecurityGroupsResponse{}, fmt.Errorf("httpcode %d", response.HttpStatusCode)
		}
		for _, group := range *response.SecurityGroups {
			groups = append(groups, cloud.SecurityGroup{
				SecurityGroupId:   group.Id,
				SecurityGroupType: "normal",
				SecurityGroupName: group.Name,
				CreateAt:          group.CreatedAt.String(),
				RegionId:          req.RegionId,
			})
		}
		secGrpNum := len(*response.SecurityGroups)
		if secGrpNum < _pageSize {
			break
		}
		markerRequest = (*response.SecurityGroups)[secGrpNum-1].Id
	}
	return cloud.DescribeSecurityGroupsResponse{Groups: groups}, nil
}

func (p *HuaweiCloud) DescribeGroupRules(req cloud.DescribeGroupRulesRequest) (cloud.DescribeGroupRulesResponse, error) {
	rules := make([]cloud.SecurityGroupRule, 0, _pageSize)
	request := &model.ShowSecurityGroupRequest{
		SecurityGroupId: req.SecurityGroupId,
	}
	response, err := p.secGrpClient.ShowSecurityGroup(request)
	if err != nil {
		return cloud.DescribeGroupRulesResponse{}, err
	}
	if response.HttpStatusCode != http.StatusOK {
		return cloud.DescribeGroupRulesResponse{}, fmt.Errorf("httpcode %d", response.HttpStatusCode)
	}

	for _, rule := range response.SecurityGroup.SecurityGroupRules {
		from, to := portRange2Int(rule.Multiport)
		protocol := _outProtocol[rule.Protocol]
		if protocol == "" {
			protocol = rule.Protocol
		}
		ipCidr := rule.RemoteIpPrefix
		if rule.RemoteGroupId == "" && rule.RemoteIpPrefix == "" && rule.RemoteAddressGroupId == "" {
			ipCidr = "0.0.0.0/0"
		}
		rules = append(rules, cloud.SecurityGroupRule{
			SecurityGroupId: response.SecurityGroup.Id,
			PortFrom:        from,
			PortTo:          to,
			Protocol:        protocol,
			Direction:       _secGrpRuleDirection[rule.Direction],
			GroupId:         rule.RemoteGroupId,
			CidrIp:          ipCidr,
			PrefixListId:    rule.RemoteAddressGroupId,
			CreateAt:        rule.CreatedAt.String(),
		})
	}

	return cloud.DescribeGroupRulesResponse{Rules: rules}, nil
}

func (p *HuaweiCloud) addSecGrpRule(req cloud.AddSecurityGroupRuleRequest, direction string) error {
	request := &model.CreateSecurityGroupRuleRequest{}
	secGrpRuleOpt := &model.CreateSecurityGroupRuleOption{
		SecurityGroupId: req.SecurityGroupId,
		Direction:       direction,
	}
	if req.IpProtocol != "" && _protocol[req.IpProtocol] != "" {
		secGrpRuleOpt.Protocol = utils.String(_protocol[req.IpProtocol])
		if req.IpProtocol == cloud.ProtocolIcmpV6 {
			secGrpRuleOpt.Ethertype = utils.String(cloud.IpV6)
		}
	}
	if (req.IpProtocol == cloud.ProtocolTcp || req.IpProtocol == cloud.ProtocolUdp) && req.PortFrom > 0 {
		secGrpRuleOpt.Multiport = utils.String(getPortRange(req.PortFrom, req.PortTo))
	}
	if req.CidrIp != "" {
		secGrpRuleOpt.RemoteIpPrefix = &req.CidrIp
	}
	if req.GroupId != "" {
		secGrpRuleOpt.RemoteGroupId = &req.GroupId
	}
	if req.PrefixListId != "" {
		secGrpRuleOpt.RemoteAddressGroupId = &req.PrefixListId
	}
	request.Body = &model.CreateSecurityGroupRuleRequestBody{
		SecurityGroupRule: secGrpRuleOpt,
	}

	response, err := p.secGrpClient.CreateSecurityGroupRule(request)
	if err != nil {
		return err
	}
	if response.HttpStatusCode != http.StatusCreated {
		return fmt.Errorf("httpcode %d", response.HttpStatusCode)
	}
	return nil
}

func getPortRange(from, to int) (portRange string) {
	if from < 1 {
		return
	}

	if from == to {
		portRange = cast.ToString(from)
	} else {
		portRange = fmt.Sprintf("%d-%d", from, to)
	}
	return
}

func portRange2Int(portRange string) (from, to int) {
	if portRange == "" {
		return 0, 0
	}

	idx := strings.Index(portRange, "-")
	if idx == -1 {
		from = cast.ToInt(portRange)
		to = from
	} else {
		from = cast.ToInt(portRange[:idx])
		to = cast.ToInt(portRange[idx+1:])
	}
	return
}
