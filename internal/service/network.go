package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
	"github.com/galaxy-future/BridgX/internal/errs"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/internal/types"
	"github.com/galaxy-future/BridgX/pkg/cloud"
)

type targetType int

const (
	TargetTypeVpc targetType = iota
	TargetTypeSwitch
	TargetTypeSecurityGroup
	TargetTypeNetwork
	TargetTypeAccount

	DefaultRegion = "cn-qingdao"
)

var H *SimpleTaskHandler

type SimpleTask struct {
	VpcId        string
	VpcName      string
	RegionId     string
	Provider     cloud.Provider
	ProviderName string
	SwitchId     string
	AccountKey   string
	TargetType   targetType
	Retry        int
}

type SimpleTaskHandler struct {
	Tasks    chan *SimpleTask
	capacity int
	running  int32
	failed   []*SimpleTask
	lock     sync.Mutex
}

func Init(workerCount int) {
	H = &SimpleTaskHandler{make(chan *SimpleTask, workerCount), workerCount, 0, make([]*SimpleTask, 0, 1000), sync.Mutex{}}
	H.run()
	RefreshCache(context.Background())
}

func (s *SimpleTaskHandler) SubmitTask(t *SimpleTask) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Logger.Errorf("SubmitTask recover : %v", r)
			}
		}()
		select {
		case s.Tasks <- t:
			fmt.Printf("有任务啦 %v", t)
			s.run()
		case <-time.After(5 * time.Minute):
			s.lock.Lock()
			s.failed = append(s.failed, t)
			s.lock.Unlock()
		}
	}()
}

func (s *SimpleTaskHandler) run() {
	if atomic.LoadInt32(&s.running) >= int32(s.capacity) {
		return
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Logger.Errorf("SimpleTaskHandler recover : %v", r)
			}
		}()
		atomic.AddInt32(&s.running, 1)
		for {
			var t *SimpleTask
			select {
			case <-time.After(1 * time.Hour):
				s.lock.Lock()
				if len(s.failed) > 0 {
					t = s.failed[0]
					s.failed = s.failed[1:]
				}
				s.lock.Unlock()
			case t = <-s.Tasks:
			}
			if t != nil && (t.VpcId != "" || t.AccountKey != "" && t.TargetType == TargetTypeAccount) {
				s.taskHandle(t)
			}
		}
	}()
}

func (s *SimpleTaskHandler) taskHandle(t *SimpleTask) {
	var err error
	switch t.TargetType {
	case TargetTypeVpc:
		err = refreshVpc(t)
	case TargetTypeSwitch:
		err = refreshSwitch(t)
	//case TargetTypeSecurityGroup:
	//	err = refreshVpc(t)
	case TargetTypeNetwork:
		err = refreshVpc(t)
		if err != nil {
			break
		}
		err = refreshSwitch(t)
	case TargetTypeAccount:
		err = refreshAccount(t)
	}
	if err == nil {
		return
	}
	logs.Logger.Errorf("taskHandle failed,task: [%v] err: [%v]", t, err)
	if t.Retry > 0 {
		t.Retry--
		s.SubmitTask(t)
	}
}

func refreshAccount(t *SimpleTask) error {
	if t.AccountKey == "" {
		return nil
	}
	ctx := context.Background()
	accounts, err := GetOrgKeysByAk(ctx, t.AccountKey)
	regions, err := GetRegions(ctx, GetRegionsRequest{
		Provider: t.ProviderName,
		Account:  accounts,
	})
	if err != nil {
		return err
	}
	vpcs := updateOrCreateVpcs(ctx, regions, t)
	updateOrCreateSwitch(ctx, vpcs, t)
	groups := updateOrCreateSecurityGroups(ctx, vpcs, t)
	updateOrCreateSecurityGroupRules(ctx, groups, t)
	return nil
}

func updateOrCreateVpcs(ctx context.Context, regions []cloud.Region, t *SimpleTask) []cloud.VPC {
	vpcs := make([]cloud.VPC, 0, 64)
	describeVpcReq := cloud.DescribeVpcsRequest{}
	for _, region := range regions {
		describeVpcReq.RegionId = region.RegionId
		provider, err := getProvider(t.ProviderName, t.AccountKey, region.RegionId)
		if err != nil {
			logs.Logger.Errorf("getProvider failed.err: %s", err.Error())
			continue
		}
		vpcsRes, err := provider.DescribeVpcs(describeVpcReq)
		if err != nil {
			continue
		}
		vpcs = append(vpcs, vpcsRes.Vpcs...)
	}
	vpcModels := cloud2ModelVpc(vpcs, t.AccountKey, t.ProviderName)
	err := model.UpdateOrCreateVpcs(ctx, vpcModels)
	if err != nil {
		logs.Logger.Errorf("updateOrCreateVpcs failed.err : [%s]", err.Error())
	}
	return vpcs
}

func updateOrCreateSwitch(ctx context.Context, vpcs []cloud.VPC, t *SimpleTask) {
	switches := make([]cloud.Switch, 0, 64)
	describeSwitchesReq := cloud.DescribeSwitchesRequest{}
	for _, vpc := range vpcs {
		describeSwitchesReq.VpcId = vpc.VpcId
		provider, err := getProvider(t.ProviderName, t.AccountKey, vpc.RegionId)
		if err != nil {
			logs.Logger.Errorf("getProvider failed.err: %s", err.Error())
			continue
		}
		switchesRes, err := provider.DescribeSwitches(describeSwitchesReq)
		if err != nil {
			continue
		}
		switches = append(switches, switchesRes.Switches...)
	}
	switchesModels := cloud2ModelSwitches(switches)
	err := model.UpdateOrCreateSwitches(ctx, switchesModels)
	if err != nil {
		logs.Logger.Errorf("updateOrCreateSwitch failed.err : [%s]", err.Error())
	}
}

func updateOrCreateSecurityGroups(ctx context.Context, vpcs []cloud.VPC, t *SimpleTask) []cloud.SecurityGroup {
	groups := make([]cloud.SecurityGroup, 0, 64)
	groupReq := cloud.GetSecurityGroupRequest{}
	for _, vpc := range vpcs {
		groupReq.VpcId = vpc.VpcId
		groupReq.RegionId = vpc.RegionId
		provider, err := getProvider(t.ProviderName, t.AccountKey, vpc.RegionId)
		if err != nil {
			logs.Logger.Errorf("getProvider failed.err: %s", err.Error())
			continue
		}
		groupRes, err := provider.GetSecurityGroup(groupReq)
		if err != nil {
			continue
		}
		groups = append(groups, groupRes.Groups...)
	}
	groupsModels := cloud2ModelGroups(groups)
	err := model.UpdateOrCreateGroups(ctx, groupsModels)
	if err != nil {
		logs.Logger.Errorf("updateOrCreateSecurityGroups failed.err : [%s]", err.Error())
	}
	return groups
}

func updateOrCreateSecurityGroupRules(ctx context.Context, groups []cloud.SecurityGroup, t *SimpleTask) {
	rulesReq := cloud.DescribeGroupRulesRequest{}
	for _, group := range groups {
		rulesReq.RegionId = group.RegionId
		rulesReq.SecurityGroupId = group.SecurityGroupId
		provider, err := getProvider(t.ProviderName, t.AccountKey, group.RegionId)
		if err != nil {
			logs.Logger.Errorf("getProvider failed.err: %s", err.Error())
			continue
		}
		rulesRes, err := provider.DescribeGroupRules(rulesReq)
		if err != nil {
			continue
		}
		rulesModels := cloud2ModelRules(rulesRes.Rules)
		err = model.ReplaceRules(ctx, group.VpcId, group.SecurityGroupId, rulesModels)
		if err != nil {
			logs.Logger.Errorf("updateOrCreateSecurityGroupRules failed.err : [%s]", err.Error())
		}
	}
}

func cloud2ModelVpc(vpcs []cloud.VPC, ak, provider string) []model.Vpc {
	res := make([]model.Vpc, 0, len(vpcs))
	for _, vpc := range vpcs {
		createAt, _ := time.Parse("2006-01-02T15:04:05Z", vpc.CreateAt)
		now := time.Now()
		res = append(res, model.Vpc{
			Base: model.Base{
				CreateAt: &createAt,
				UpdateAt: &now,
			},
			Ak:        ak,
			RegionId:  vpc.RegionId,
			VpcId:     vpc.VpcId,
			Name:      vpc.VpcName,
			CidrBlock: vpc.CidrBlock,
			SwitchIds: strings.Join(vpc.SwitchIds, ","),
			Provider:  provider,
			VStatus:   vpc.Status,
		})
	}
	return res
}

func cloud2ModelSwitches(switches []cloud.Switch) []model.Switch {
	res := make([]model.Switch, 0, len(switches))
	for _, sw := range switches {
		createAt, _ := time.Parse("2006-01-02T15:04:05Z", sw.CreateAt)
		now := time.Now()
		res = append(res, model.Switch{
			Base: model.Base{
				CreateAt: &createAt,
				UpdateAt: &now,
			},
			VpcId:                   sw.VpcId,
			SwitchId:                sw.SwitchId,
			ZoneId:                  sw.ZoneId,
			Name:                    sw.Name,
			CidrBlock:               sw.CidrBlock,
			IsDefault:               sw.IsDefault,
			AvailableIpAddressCount: sw.AvailableIpAddressCount,
			VStatus:                 sw.VStatus,
		})
	}
	return res
}

func cloud2ModelGroups(groups []cloud.SecurityGroup) []model.SecurityGroup {
	res := make([]model.SecurityGroup, 0, len(groups))
	for _, group := range groups {
		createAt, _ := time.Parse("2006-01-02T15:04:05Z", group.CreateAt)
		now := time.Now()
		res = append(res, model.SecurityGroup{
			Base: model.Base{
				CreateAt: &createAt,
				UpdateAt: &now,
			},
			VpcId:             group.VpcId,
			SecurityGroupId:   group.SecurityGroupId,
			Name:              group.SecurityGroupName,
			SecurityGroupType: group.SecurityGroupType,
		})
	}
	return res
}

func cloud2ModelRules(rules []cloud.SecurityGroupRule) []model.SecurityGroupRule {
	res := make([]model.SecurityGroupRule, 0, len(rules))
	for _, rule := range rules {
		createAt, _ := time.Parse("2006-01-02T15:04:05Z", rule.CreateAt)
		now := time.Now()
		res = append(res, model.SecurityGroupRule{
			Base: model.Base{
				CreateAt: &createAt,
				UpdateAt: &now,
			},
			VpcId:           rule.VpcId,
			SecurityGroupId: rule.SecurityGroupId,
			PortRange:       rule.PortRange,
			Protocol:        rule.Protocol,
			Direction:       rule.Direction,
			GroupId:         rule.GroupId,
			CidrIp:          rule.CidrIp,
			PrefixListId:    rule.PrefixListId,
		})
	}
	return res
}

func refreshVpc(t *SimpleTask) error {
	res, err := t.Provider.GetVPC(cloud.GetVpcRequest{
		VpcId:    t.VpcId,
		RegionId: t.RegionId,
		VpcName:  t.VpcName,
	})
	if err != nil {
		logs.Logger.Errorf("refreshVpc failed.task: [%v]", t)
		return nil
	}
	vpc := res.Vpc
	return model.UpdateVpc(context.Background(), vpc.VpcId, vpc.CidrBlock, vpc.Status, vpc.SwitchIds)
}

func refreshSwitch(t *SimpleTask) error {
	res, err := t.Provider.GetSwitch(cloud.GetSwitchRequest{
		SwitchId: t.SwitchId,
	})
	if err != nil {
		return nil
	}
	vswitch := res.Switch
	return model.UpdateSwitch(context.Background(),
		vswitch.AvailableIpAddressCount, vswitch.IsDefault,
		vswitch.VpcId, vswitch.SwitchId, vswitch.Name,
		vswitch.VStatus, vswitch.CidrBlock)
}

const (
	DirectionIn  = "ingress"
	DirectionOut = "egress"
)

type CreateNetworkRequest struct {
	Provider          string
	RegionId          string
	CidrBlock         string
	VpcName           string
	ZoneId            string
	SwitchCidrBlock   string
	SwitchName        string
	SecurityGroupName string
	SecurityGroupType string
	Ak                string
}

type CreateNetworkResponse struct {
	VpcId           string
	SwitchId        string
	SecurityGroupId string
}

type CreateVPCRequest struct {
	Provider  string
	RegionId  string
	VpcName   string
	CidrBlock string
	Ak        string
}

func CreateNetwork(ctx context.Context, req *CreateNetworkRequest) (vpcRes CreateNetworkResponse, err error) {
	// createVpc
	vpcId, err := CreateVPC(ctx, CreateVPCRequest{
		Provider:  req.Provider,
		RegionId:  req.RegionId,
		VpcName:   req.VpcName,
		CidrBlock: req.CidrBlock,
		Ak:        req.Ak,
	})
	if err != nil {
		return CreateNetworkResponse{}, err
	}
	err = waitForVpcStatus(ctx, req, vpcId)
	if err != nil {
		return CreateNetworkResponse{}, err
	}
	switchId, err := CreateSwitch(ctx, CreateSwitchRequest{
		SwitchName: req.SwitchName,
		ZoneId:     req.ZoneId,
		VpcId:      vpcId,
		CidrBlock:  req.SwitchCidrBlock,
	})
	if err != nil {
		return CreateNetworkResponse{}, err
	}

	groupId, err := CreateSecurityGroup(ctx, CreateSecurityGroupRequest{
		VpcId:             vpcId,
		SecurityGroupName: req.SecurityGroupName,
		SecurityGroupType: req.SecurityGroupType,
	})
	if err != nil {
		return CreateNetworkResponse{}, err
	}
	return CreateNetworkResponse{
		VpcId:           vpcId,
		SwitchId:        switchId,
		SecurityGroupId: groupId,
	}, nil
}

func waitForVpcStatus(ctx context.Context, req *CreateNetworkRequest, vpcId string) error {
	getVpc := func(attempt uint) error {
		vpc, err := GetVPCFromCloud(ctx, GetVPCFromCloudRequest{
			Provider:   req.Provider,
			RegionId:   req.RegionId,
			VpcId:      vpcId,
			PageNumber: 1,
			PageSize:   10,
			Ak:         req.Ak,
		})
		if err != nil {
			return err
		}
		if vpc.Status == cloud.VPCStatusAvailable {
			return nil
		}
		if vpc.Status == cloud.VPCStatusPending {
			return errs.ErrVpcPending
		}
		return nil
	}

	return retry.Retry(getVpc, strategy.Limit(10), strategy.Backoff(backoff.BinaryExponential(10*time.Millisecond)))
}

func CreateVPC(ctx context.Context, req CreateVPCRequest) (vpcId string, err error) {
	/* name 如果限制了再打开这部分
	vpcIDStruct, err := model.FindVpcId(ctx, model.FindVpcConditions{
		AK:       req.Account.AK,
		VpcName:  req.VpcName,
		RegionId: req.RegionId,
	})
	if err != nil {
		return "", errs.ErrDBQueryFailed
	}
	if vpcIDStruct.VpcId != "" {
		return "", errs.ErrVpcNameExist
	}
	*/

	p, err := getProvider(req.Provider, req.Ak, req.RegionId)
	if err != nil {
		return "", err
	}

	res, err := p.CreateVPC(cloud.CreateVpcRequest{
		RegionId:  req.RegionId,
		VpcName:   req.VpcName,
		CidrBlock: req.CidrBlock,
	})
	if err != nil {
		return "", errs.ErrCreateVpcFailed
	}

	now := time.Now()
	err = model.CreateVpc(ctx, model.Vpc{
		Base: model.Base{
			CreateAt: &now,
			UpdateAt: &now,
		},
		Ak:        req.Ak,
		RegionId:  req.RegionId,
		VpcId:     res.VpcId,
		Name:      req.VpcName,
		CidrBlock: req.CidrBlock,
		Provider:  req.Provider,
	})
	if err != nil {
		logs.Logger.Errorf("save Vpc failed: %v, error: %v", res, err.Error())
		return "", nil
	}
	H.SubmitTask(&SimpleTask{
		VpcId:      res.VpcId,
		RegionId:   req.RegionId,
		Provider:   p,
		TargetType: TargetTypeVpc,
		Retry:      3,
	})
	return res.VpcId, nil
}

type GetVPCRequest struct {
	Provider   string
	RegionId   string
	VpcName    string
	PageNumber int
	PageSize   int

	Account *types.OrgKeys
}
type VPCResponse struct {
	Vpcs  []Vpc
	Pager types.Pager
}

type Vpc struct {
	VpcId     string
	VpcName   string
	CidrBlock string
	SwitchIds string
	Provider  string
	Status    string
	CreateAt  string
}

func model2VpcResponse(vpcs []model.Vpc, pageNumber, pageSize, total int) VPCResponse {
	vs := make([]Vpc, 0, len(vpcs))
	for _, v := range vpcs {
		vs = append(vs, Vpc{
			VpcId:     v.VpcId,
			VpcName:   v.Name,
			CidrBlock: v.CidrBlock,
			SwitchIds: v.SwitchIds,
			Provider:  v.Provider,
			Status:    v.VStatus,
			CreateAt:  v.CreateAt.String()})

	}
	return VPCResponse{
		Vpcs: vs,
		Pager: types.Pager{
			PageNumber: pageNumber,
			PageSize:   pageSize,
			Total:      total,
		},
	}
}

func GetVPC(ctx context.Context, req GetVPCRequest) (resp VPCResponse, err error) {
	// TODO: cache
	aks := make([]string, 0, len(req.Account.Info))
	for _, info := range req.Account.Info {
		aks = append(aks, info.AK)
		if req.Provider != "" {
			if info.Provider != req.Provider {
				continue
			}
		}
	}
	vs, total, err := model.FindVpcsWithPage(ctx, model.FindVpcConditions{
		Aks:        aks,
		VpcName:    req.VpcName,
		RegionId:   req.RegionId,
		PageNumber: req.PageNumber,
		PageSize:   req.PageSize,
	})
	return model2VpcResponse(vs, req.PageNumber, req.PageSize, int(total)), err
}

type GetVPCFromCloudRequest struct {
	Provider   string
	RegionId   string
	VpcName    string
	PageNumber int32
	PageSize   int32
	VpcId      string
	Ak         string
}

func GetVPCFromCloud(ctx context.Context, req GetVPCFromCloudRequest) (vpc cloud.VPC, err error) {
	// TODO: cache
	p, err := getProvider(req.Provider, req.Ak, req.RegionId)
	if err != nil {
		return cloud.VPC{}, err
	}
	res, err := p.GetVPC(cloud.GetVpcRequest{
		VpcId:    req.VpcId,
		RegionId: req.RegionId,
		VpcName:  req.VpcName,
	})
	if err != nil {
		return cloud.VPC{}, err
	}
	return res.Vpc, nil
}

type CreateSwitchRequest struct {
	SwitchName string
	ZoneId     string
	VpcId      string
	CidrBlock  string
}

func CreateSwitch(ctx context.Context, req CreateSwitchRequest) (switchId string, err error) {
	vpc, err := model.FindVpcById(ctx, model.FindVpcConditions{
		VpcId: req.VpcId,
	})
	if err != nil {
		logs.Logger.Errorf("FindVpcById failed.err: [%v] req[%v]", err, req)
		return "", errs.ErrDBQueryFailed
	}
	if vpc.VpcId == "" {
		return "", errs.ErrVpcNotExist
	}
	vpcId := vpc.VpcId
	/* name 如果限制了再打开这部分
	switchIdstruct, err := model.FindSwitchId(ctx, model.FindSwitchesConditions{VpcId: vpcId, SwitchName: req.SwitchName})
	if err != nil {
		return "", errs.ErrDBQueryFailed
	}
	if switchIdstruct.SwitchId != "" {
		return "", errs.ErrSwitchNameExist
	}

	*/

	p, err := getProvider(vpc.Provider, vpc.Ak, vpc.RegionId)
	if err != nil {
		return "", err
	}
	// TODO: lock
	// TODO: defer unlock

	res, err := p.CreateSwitch(cloud.CreateSwitchRequest{
		RegionId:    vpc.RegionId,
		ZoneId:      req.ZoneId,
		CidrBlock:   req.CidrBlock,
		VSwitchName: req.SwitchName,
		VpcId:       vpcId,
	})
	if err != nil {
		return "", errs.ErrCreateSwitchFailed
	}

	now := time.Now()
	err = model.CreateSwitch(ctx, model.Switch{
		Base: model.Base{
			CreateAt: &now,
			UpdateAt: &now,
		},
		VpcId:     vpcId,
		SwitchId:  res.SwitchId,
		ZoneId:    req.ZoneId,
		Name:      req.SwitchName,
		CidrBlock: req.CidrBlock,
		IsDel:     0,
	})
	if err != nil {
		logs.Logger.Errorf("save Switch failed: %v, error: %v", res, err.Error())
		return "", nil
	}
	H.SubmitTask(&SimpleTask{
		VpcId:      req.VpcId,
		RegionId:   vpc.RegionId,
		Provider:   p,
		TargetType: TargetTypeSwitch,
		Retry:      3,
	})
	return res.SwitchId, nil
}

type GetSwitchRequest struct {
	SwitchName string
	VpcId      string
	PageNumber int
	PageSize   int
}
type Switch struct {
	VpcId                   string
	SwitchId                string
	ZoneId                  string
	SwitchName              string
	CidrBlock               string
	VStatus                 string
	CreateAt                string
	IsDefault               string
	AvailableIpAddressCount int
}

type SwitchResponse struct {
	Switches []Switch
	Pager    types.Pager
}

func model2SwitchResponse(switches []model.Switch, pageNumber, pageSize, total int) SwitchResponse {
	vs := make([]Switch, 0, len(switches))
	for _, v := range switches {
		isDefault := "Y"
		if v.IsDefault == 0 {
			isDefault = "N"
		}
		vs = append(vs, Switch{
			VpcId:                   v.VpcId,
			SwitchId:                v.SwitchId,
			ZoneId:                  v.ZoneId,
			SwitchName:              v.Name,
			CidrBlock:               v.CidrBlock,
			VStatus:                 v.VStatus,
			CreateAt:                v.CreateAt.String(),
			IsDefault:               isDefault,
			AvailableIpAddressCount: v.AvailableIpAddressCount,
		})
	}
	return SwitchResponse{
		Switches: vs,
		Pager: types.Pager{
			PageNumber: pageNumber,
			PageSize:   pageSize,
			Total:      total,
		},
	}
}

func GetSwitch(ctx context.Context, req GetSwitchRequest) (resp SwitchResponse, err error) {
	//TODO: cache
	vpc, err := model.FindVpcById(ctx, model.FindVpcConditions{
		VpcId: req.VpcId,
	})
	if err != nil {
		logs.Logger.Errorf("FindVpcById failed.err: [%v] req[%v]", err, req)
		return SwitchResponse{}, errs.ErrDBQueryFailed
	}
	if vpc.VpcId == "" {
		return SwitchResponse{}, errs.ErrVpcNotExist
	}
	vpcId := vpc.VpcId
	s, total, err := model.FindSwitchesWithPage(ctx, model.FindSwitchesConditions{
		VpcId:      vpcId,
		SwitchName: req.SwitchName,
		PageNumber: req.PageNumber,
		PageSize:   req.PageSize,
	})
	if err != nil {
		logs.Logger.Errorf("FindSwitchesWithPage failed.err: [%v] req[%v]", err, req)
		return SwitchResponse{}, errs.ErrDBQueryFailed
	}
	return model2SwitchResponse(s, req.PageNumber, req.PageSize, int(total)), nil
}

type CreateSecurityGroupRequest struct {
	VpcId             string
	SecurityGroupName string
	SecurityGroupType string
}

func CreateSecurityGroup(ctx context.Context, req CreateSecurityGroupRequest) (securityGroupId string, err error) {
	vpc, err := model.FindVpcById(ctx, model.FindVpcConditions{
		VpcId: req.VpcId,
	})
	if err != nil {
		logs.Logger.Errorf("FindVpcById failed.err: [%v] req[%v]", err, req)
		return "", errs.ErrDBQueryFailed
	}

	if vpc.VpcId == "" {
		return "", errs.ErrVpcNotExist
	}
	vpcId := vpc.VpcId
	/* name 如果限制了再打开这部分
	groupIdStruct, err := model.FindSecurityId(ctx, model.FindSecurityGroupConditions{VpcId: vpcId, SecurityGroupName: req.SecurityGroupName})
	if err != nil {
		return "", errs.ErrDBQueryFailed
	}
	if groupIdStruct.GroupId != "" {
		return "", errs.ErrSecurityGroupNameExist
	}

	*/

	p, err := getProvider(vpc.Provider, vpc.Ak, vpc.RegionId)
	if err != nil {
		return "", err
	}
	// TODO: lock
	// TODO: defer unlock
	if req.SecurityGroupType == "" {
		req.SecurityGroupType = "normal"
	}
	res, err := p.CreateSecurityGroup(cloud.CreateSecurityGroupRequest{
		RegionId:          vpc.RegionId,
		SecurityGroupName: req.SecurityGroupName,
		VpcId:             vpcId,
		SecurityGroupType: req.SecurityGroupType,
	})
	if err != nil {
		return "", errs.ErrCreateSecurityGroupFailed
	}
	now := time.Now()
	err = model.CreateSecurityGroup(ctx, model.SecurityGroup{
		Base: model.Base{
			CreateAt: &now,
			UpdateAt: &now,
		},
		VpcId:             vpcId,
		SecurityGroupId:   res.SecurityGroupId,
		Name:              req.SecurityGroupName,
		SecurityGroupType: req.SecurityGroupType,
		IsDel:             0,
	})
	if err != nil {
		logs.Logger.Errorf("save security group failed: %v, error: %v", res, err.Error())
		return "", nil
	}
	return res.SecurityGroupId, nil
}

type GetSecurityGroupRequest struct {
	SecurityGroupName string
	VpcId             string
	PageNumber        int
	PageSize          int
}

type SecurityGroupResponse struct {
	Groups []Group
	Pager  types.Pager
}

type Group struct {
	VpcId             string
	SecurityGroupId   string
	SecurityGroupName string
	SecurityGroupType string
	CreateAt          string
}

func model2SecurityGroupResponse(groups []model.SecurityGroup, pageNumber, pageSize, total int) SecurityGroupResponse {
	res := make([]Group, 0, len(groups))
	for _, s := range groups {
		res = append(res, Group{
			VpcId:             s.VpcId,
			SecurityGroupId:   s.SecurityGroupId,
			SecurityGroupName: s.Name,
			SecurityGroupType: s.SecurityGroupType,
			CreateAt:          s.CreateAt.String(),
		})
	}
	return SecurityGroupResponse{
		Groups: res,
		Pager: types.Pager{
			PageNumber: pageNumber,
			PageSize:   pageSize,
			Total:      total,
		},
	}
}

func GetSecurityGroup(ctx context.Context, req GetSecurityGroupRequest) (SecurityGroupResponse, error) {
	//TODO: cache
	vpc, err := model.FindVpcById(ctx, model.FindVpcConditions{
		VpcId: req.VpcId,
	})
	if err != nil {
		logs.Logger.Errorf("FindVpcById failed.err: [%v] req[%v]", err, req)
		return SecurityGroupResponse{}, errs.ErrDBQueryFailed
	}
	if vpc.VpcId == "" {
		return SecurityGroupResponse{}, errs.ErrVpcNotExist
	}
	vpcId := vpc.VpcId
	groups, total, err := model.FindSecurityGroupWithPage(ctx, model.FindSecurityGroupConditions{
		VpcId:             vpcId,
		SecurityGroupName: req.SecurityGroupName,
		PageNumber:        req.PageNumber,
		PageSize:          req.PageSize,
	})
	if err != nil {
		logs.Logger.Errorf("FindSecurityGroupWithPage failed.err: [%v] req[%v]", err, req)
		return SecurityGroupResponse{}, errs.ErrDBQueryFailed
	}
	return model2SecurityGroupResponse(groups, req.PageNumber, req.PageSize, int(total)), nil
}

type AddSecurityGroupRuleRequest struct {
	RegionId        string
	VpcId           string
	SecurityGroupId string
	Rules           []GroupRule
}

type GroupRule struct {
	Protocol     string `json:"protocol"`
	PortRange    string `json:"port_range"`
	Direction    string `json:"direction"`
	GroupId      string `json:"group_id"`
	CidrIp       string `json:"cidr_ip"`
	PrefixListId string `json:"prefix_list_id"`
}

func AddSecurityGroupRule(ctx context.Context, req AddSecurityGroupRuleRequest) (string, error) {
	vpc, err := model.FindVpcById(ctx, model.FindVpcConditions{
		VpcId:    req.VpcId,
		RegionId: req.RegionId,
	})
	if err != nil {
		logs.Logger.Errorf("FindVpcById failed.err: [%v] req[%v]", err, req)
		return "", errs.ErrDBQueryFailed
	}

	if vpc.VpcId == "" {
		return "", errs.ErrVpcNotExist
	}
	vpcId := vpc.VpcId
	groupIdStruct, err := model.FindSecurityId(ctx, model.FindSecurityGroupConditions{
		VpcId: vpcId, SecurityGroupId: req.SecurityGroupId})
	if err != nil {
		logs.Logger.Errorf("FindSecurityId failed.err: [%v] req[%v]", err, req)
		return "", errs.ErrDBQueryFailed
	}
	if groupIdStruct.SecurityGroupId == "" {
		return "", errs.ErrSecurityGroupNotExist
	}

	p, err := getProvider(vpc.Provider, vpc.Ak, req.RegionId)
	if err != nil {
		return "", err
	}
	// TODO: lock
	// TODO: defer unlock
	ruleModels := make([]model.SecurityGroupRule, 0)
	for _, rule := range req.Rules {
		addRuleReq := cloud.AddSecurityGroupRuleRequest{
			RegionId:        req.RegionId,
			VpcId:           vpcId,
			SecurityGroupId: groupIdStruct.SecurityGroupId,
			IpProtocol:      rule.Protocol,
			PortRange:       rule.PortRange,
			GroupId:         rule.GroupId,
			CidrIp:          rule.CidrIp,
			PrefixListId:    rule.PrefixListId,
		}
		switch rule.Direction {
		case DirectionIn:
			err = p.AddIngressSecurityGroupRule(addRuleReq)
		case DirectionOut:
			err = p.AddEgressSecurityGroupRule(addRuleReq)
		}
		if err != nil {
			continue
		}
		now := time.Now()
		ruleModels = append(ruleModels, model.SecurityGroupRule{
			Base: model.Base{
				CreateAt: &now,
				UpdateAt: &now,
			},
			VpcId:           vpcId,
			SecurityGroupId: groupIdStruct.SecurityGroupId,
			PortRange:       rule.PortRange,
			Protocol:        rule.Protocol,
			Direction:       rule.Direction,
			GroupId:         rule.GroupId,
			CidrIp:          rule.CidrIp,
			PrefixListId:    rule.PrefixListId,
		})
	}

	err = model.BatchCreate(ruleModels)
	if err != nil {
		logs.Logger.Errorf("save security group rules failed. error: %s", err.Error())
		return "", nil
	}
	return "", nil
}

type GetRegionsRequest struct {
	Provider string
	Account  *types.OrgKeys
}

func GetRegions(ctx context.Context, req GetRegionsRequest) ([]cloud.Region, error) {
	ak := getFirstAk(req.Account, req.Provider)
	p, err := getProvider(req.Provider, ak, DefaultRegion)
	if err != nil {
		return nil, err
	}
	regions, err := p.GetRegions()
	if err != nil {
		return nil, errs.ErrGetRegionsFailed
	}
	return regions.Regions, nil
}

type GetZonesRequest struct {
	Provider string
	RegionId string
	Account  *types.OrgKeys
}

func GetZones(ctx context.Context, req GetZonesRequest) ([]cloud.Zone, error) {
	ak := getFirstAk(req.Account, req.Provider)
	p, err := getProvider(req.Provider, ak, req.RegionId)
	if err != nil {
		return nil, err
	}
	zones, err := p.GetZones(cloud.GetZonesRequest{
		RegionId: req.RegionId,
	})
	if err != nil {
		return nil, errs.ErrGetZonesFailed
	}
	return zones.Zones, nil
}

func getFirstAk(account *types.OrgKeys, provider string) string {
	for _, a := range account.Info {
		if a.Provider == provider {
			return a.AK
		}
	}
	return ""
}
