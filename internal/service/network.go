package service

import (
	"context"
	"fmt"
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
	"github.com/spf13/cast"
)

type targetType int

const (
	TargetTypeVpc targetType = iota
	TargetTypeSwitch
	TargetTypeSecurityGroup
	TargetTypeNetwork
	TargetTypeAccount
	TargetTypeInstanceType

	DefaultRegion        = "cn-qingdao"
	DefaultRegionHuaWei  = "cn-north-4"
	DefaultRegionTencent = "ap-beijing"
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
	RefreshCache()
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
			if t != nil {
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
		err = RefreshAccount(t)
	case TargetTypeInstanceType:
		err = refreshInstanceType(t)
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

func refreshInstanceType(t *SimpleTask) error {
	err := SyncInstanceTypes(context.Background(), t.ProviderName)
	if err != nil {
		logs.Logger.Error("SyncInstanceTypes failed :%v", err)
		return err
	}
	return RefreshCache()
}

func RefreshAccount(t *SimpleTask) error {
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

	regionIds := make([]string, 0, len(regions))
	for _, region := range regions {
		regionIds = append(regionIds, region.RegionId)
	}
	err = syncNetworkConfig(ctx, regionIds, t.ProviderName, t.AccountKey)
	if err != nil {
		return err
	}
	return nil
}

func syncNetworkConfig(ctx context.Context, regionIds []string, provider, ak string) error {
	vpcs, err := updateOrCreateVpcs(ctx, regionIds, provider, ak)
	if err != nil {
		return err
	}
	updateOrCreateSwitch(ctx, vpcs, provider, ak)

	groups, err := updateOrCreateSecurityGroups(ctx, regionIds, vpcs, provider, ak)
	if err != nil {
		return err
	}
	updateOrCreateSecurityGroupRules(ctx, groups, provider, ak)
	return nil
}

func updateOrCreateVpcs(ctx context.Context, regionIds []string, provider, ak string) ([]cloud.VPC, error) {
	vpcs := make([]cloud.VPC, 0, 64)
	describeVpcReq := cloud.DescribeVpcsRequest{}
	for _, regionId := range regionIds {
		describeVpcReq.RegionId = regionId
		cloudCli, err := getProvider(provider, ak, regionId)
		if err != nil {
			logs.Logger.Errorf("getProvider failed: %s", err.Error())
			continue
		}
		vpcsRes, err := cloudCli.DescribeVpcs(describeVpcReq)
		if err != nil {
			logs.Logger.Errorf("DescribeVpcs failed: %s", err.Error())
			continue
		}
		vpcs = append(vpcs, vpcsRes.Vpcs...)
	}
	vpcModels := cloud2ModelVpc(vpcs, ak, provider)
	err := model.UpdateOrCreateVpcs(ctx, vpcModels)
	return vpcs, err
}

func updateOrCreateSwitch(ctx context.Context, vpcs []cloud.VPC, provider, ak string) {
	switches := make([]cloud.Switch, 0, 64)
	describeSwitchesReq := cloud.DescribeSwitchesRequest{}
	for _, vpc := range vpcs {
		describeSwitchesReq.VpcId = vpc.VpcId
		cloudCli, err := getProvider(provider, ak, vpc.RegionId)
		if err != nil {
			logs.Logger.Errorf("getProvider failed.err: %s", err.Error())
			continue
		}
		switchesRes, err := cloudCli.DescribeSwitches(describeSwitchesReq)
		if err != nil {
			logs.Logger.Errorf("DescribeSwitches failed: %s", err.Error())
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

func DescribeSecurityGroups(provider, ak, regionId, vpcId string) ([]cloud.SecurityGroup, error) {
	groupReq := cloud.DescribeSecurityGroupsRequest{}
	groupReq.VpcId = vpcId
	groupReq.RegionId = regionId
	cloudCli, err := getProvider(provider, ak, regionId)
	if err != nil {
		return nil, err
	}
	groupRes, err := cloudCli.DescribeSecurityGroups(groupReq)
	if err != nil {
		return nil, err
	}
	return groupRes.Groups, nil
}

func DoesSecurityGroupBelongsVpc(provider string) bool {
	if provider == cloud.HuaweiCloud || provider == cloud.TencentCloud {
		return false
	}
	return true
}

func updateOrCreateSecurityGroups(ctx context.Context, regionIds []string, vpcs []cloud.VPC, provider,
	ak string) ([]cloud.SecurityGroup, error) {
	groups := make([]cloud.SecurityGroup, 0, 64)
	if DoesSecurityGroupBelongsVpc(provider) {
		for _, vpc := range vpcs {
			secGroups, err := DescribeSecurityGroups(provider, ak, vpc.RegionId, vpc.VpcId)
			if err != nil {
				logs.Logger.Errorf("DescribeSecurityGroups failed: %s", err.Error())
				continue
			}
			groups = append(groups, secGroups...)
		}
	} else {
		for _, regionId := range regionIds {
			secGroups, err := DescribeSecurityGroups(provider, ak, regionId, "")
			if err != nil {
				logs.Logger.Errorf("DescribeSecurityGroups failed: %s", err.Error())
				continue
			}
			groups = append(groups, secGroups...)
		}
	}

	groupsModels := cloud2ModelGroups(groups, ak, provider)
	err := model.UpdateOrCreateGroups(ctx, groupsModels)
	return groups, err
}

func updateOrCreateSecurityGroupRules(ctx context.Context, groups []cloud.SecurityGroup, provider, ak string) {
	rulesReq := cloud.DescribeGroupRulesRequest{}
	for _, group := range groups {
		rulesReq.RegionId = group.RegionId
		rulesReq.SecurityGroupId = group.SecurityGroupId
		cloudCli, err := getProvider(provider, ak, group.RegionId)
		if err != nil {
			logs.Logger.Errorf("getProvider failed.err: %s", err.Error())
			continue
		}
		rulesRes, err := cloudCli.DescribeGroupRules(rulesReq)
		if err != nil {
			logs.Logger.Errorf("DescribeGroupRules failed.err: %s", err.Error())
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
		now := time.Now()
		createAt, err := time.Parse("2006-01-02T15:04:05Z", vpc.CreateAt)
		if err != nil {
			createAt = now
		} else {
			createAt = createAt.Local()
		}

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
			Provider:  provider,
			VStatus:   vpc.Status,
		})
	}
	return res
}

func cloud2ModelSwitches(switches []cloud.Switch) []model.Switch {
	res := make([]model.Switch, 0, len(switches))
	for _, sw := range switches {
		now := time.Now()
		createAt, err := time.Parse("2006-01-02T15:04:05Z", sw.CreateAt)
		if err != nil {
			createAt = now
		} else {
			createAt = createAt.Local()
		}

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
			GatewayIp:               sw.GatewayIp,
			IsDefault:               sw.IsDefault,
			AvailableIpAddressCount: sw.AvailableIpAddressCount,
			VStatus:                 sw.VStatus,
		})
	}
	return res
}

func cloud2ModelGroups(groups []cloud.SecurityGroup, ak, provider string) []model.SecurityGroup {
	res := make([]model.SecurityGroup, 0, len(groups))
	for _, group := range groups {
		now := time.Now()
		createAt, err := time.Parse("2006-01-02T15:04:05Z", group.CreateAt)
		if err != nil {
			createAt = now
		} else {
			createAt = createAt.Local()
		}

		res = append(res, model.SecurityGroup{
			Base: model.Base{
				CreateAt: &createAt,
				UpdateAt: &now,
			},
			Ak:                ak,
			Provider:          provider,
			RegionId:          group.RegionId,
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
		now := time.Now()
		createAt, err := time.Parse("2006-01-02T15:04:05Z", rule.CreateAt)
		if err != nil {
			createAt = now
		} else {
			createAt = createAt.Local()
		}

		res = append(res, model.SecurityGroupRule{
			Base: model.Base{
				CreateAt: &createAt,
				UpdateAt: &now,
			},
			VpcId:           rule.VpcId,
			SecurityGroupId: rule.SecurityGroupId,
			PortRange:       getPortRange(rule.PortFrom, rule.PortTo),
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
	if t.VpcId == "" {
		return nil
	}

	res, err := t.Provider.GetVPC(cloud.GetVpcRequest{
		VpcId:    t.VpcId,
		RegionId: t.RegionId,
		VpcName:  t.VpcName,
	})
	if err != nil {
		return err
	}
	vpc := res.Vpc
	return model.UpdateVpc(context.Background(), vpc.VpcId, vpc.CidrBlock, vpc.Status)
}

func refreshSwitch(t *SimpleTask) error {
	if t.SwitchId == "" {
		return nil
	}

	res, err := t.Provider.GetSwitch(cloud.GetSwitchRequest{
		SwitchId: t.SwitchId,
	})
	if err != nil {
		return err
	}
	vswitch := res.Switch
	return model.UpdateSwitch(context.Background(),
		vswitch.AvailableIpAddressCount, vswitch.IsDefault,
		vswitch.VpcId, vswitch.SwitchId, vswitch.Name,
		vswitch.VStatus, vswitch.CidrBlock, vswitch.GatewayIp)
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
	GatewayIp         string
	SwitchName        string
	SecurityGroupName string
	SecurityGroupType string
	Ak                string
	Rules             []GroupRule
}

type CreateNetworkResponse struct {
	VpcId           string
	SwitchId        string
	SecurityGroupId string
}

type SyncNetworkRequest struct {
	Provider   string
	RegionId   string
	AccountKey string
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
		Ak:         req.Ak,
		SwitchName: req.SwitchName,
		ZoneId:     req.ZoneId,
		VpcId:      vpcId,
		CidrBlock:  req.SwitchCidrBlock,
		GatewayIp:  req.GatewayIp,
	})
	if err != nil {
		return CreateNetworkResponse{}, err
	}

	groupId, err := CreateSecurityGroup(ctx, CreateSecurityGroupRequest{
		Ak:                req.Ak,
		VpcId:             vpcId,
		SecurityGroupName: req.SecurityGroupName,
		SecurityGroupType: req.SecurityGroupType,
	})
	if err != nil {
		return CreateNetworkResponse{}, err
	}

	if len(req.Rules) == 0 || req.Rules[0].Protocol == "" {
		return CreateNetworkResponse{}, fmt.Errorf("miss rule")
	}
	_, err = AddSecurityGroupRule(ctx, AddSecurityGroupRuleRequest{
		Ak:              req.Ak,
		RegionId:        req.RegionId,
		VpcId:           vpcId,
		SecurityGroupId: groupId,
		Rules:           req.Rules,
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

func SyncNetwork(ctx context.Context, req SyncNetworkRequest) error {
	return syncNetworkConfig(ctx, []string{req.RegionId}, req.Provider, req.AccountKey)
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
	AccountKey string
}
type VPCResponse struct {
	Vpcs  []Vpc
	Pager types.Pager
}

type Vpc struct {
	VpcId     string
	VpcName   string
	CidrBlock string
	Provider  string
	Status    string
	CreateAt  string
}

func model2Vpc(v model.Vpc) Vpc {
	return Vpc{
		VpcId:     v.VpcId,
		VpcName:   v.Name,
		CidrBlock: v.CidrBlock,
		Provider:  v.Provider,
		Status:    v.VStatus,
		CreateAt:  v.CreateAt.String()}
}

func model2VpcResponse(vpcs []model.Vpc, pageNumber, pageSize, total int) VPCResponse {
	vs := make([]Vpc, 0, len(vpcs))
	for _, v := range vpcs {
		if v.VStatus != cloud.VPCStatusAvailable {
			continue
		}
		vs = append(vs, model2Vpc(v))
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

func GetVpcById(ctx context.Context, vpcId string) (Vpc, error) {
	vpc, err := model.FindVpcById(ctx, model.FindVpcConditions{
		VpcId: vpcId,
	})
	if err != nil {
		return Vpc{}, err
	}

	return model2Vpc(vpc), nil
}

func GetVPC(ctx context.Context, req GetVPCRequest) (resp VPCResponse, err error) {
	// TODO: cache
	vs, total, err := model.FindVpcsWithPage(ctx, model.FindVpcConditions{
		AccountKey: req.AccountKey,
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
	Ak         string
	SwitchName string
	ZoneId     string
	VpcId      string
	CidrBlock  string
	GatewayIp  string
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

	p, err := getProvider(vpc.Provider, req.Ak, vpc.RegionId)
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
		GatewayIp:   req.GatewayIp,
	})
	if err != nil {
		logs.Logger.Errorf("CreateSwitch failed, %s", err.Error())
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
		GatewayIp: req.GatewayIp,
		IsDel:     0,
	})
	if err != nil {
		logs.Logger.Errorf("save Switch failed: %v, error: %v", res, err.Error())
		return "", nil
	}
	H.SubmitTask(&SimpleTask{
		VpcId:      req.VpcId,
		SwitchId:   res.SwitchId,
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
	ZoneId     string
	PageNumber int
	PageSize   int
}
type Switch struct {
	VpcId                   string
	SwitchId                string
	ZoneId                  string
	SwitchName              string
	CidrBlock               string
	GatewayIp               string
	VStatus                 string
	CreateAt                string
	IsDefault               string
	AvailableIpAddressCount int
}

type SwitchResponse struct {
	Switches []Switch
	Pager    types.Pager
}

func model2Switch(v model.Switch) Switch {
	isDefault := "N"
	if v.IsDefault == 1 {
		isDefault = "Y"
	}
	return Switch{
		VpcId:                   v.VpcId,
		SwitchId:                v.SwitchId,
		ZoneId:                  v.ZoneId,
		SwitchName:              v.Name,
		CidrBlock:               v.CidrBlock,
		GatewayIp:               v.GatewayIp,
		VStatus:                 v.VStatus,
		CreateAt:                v.CreateAt.String(),
		IsDefault:               isDefault,
		AvailableIpAddressCount: v.AvailableIpAddressCount,
	}
}

func model2SwitchResponse(switches []model.Switch, pageNumber, pageSize, total int) SwitchResponse {
	vs := make([]Switch, 0, len(switches))
	for _, v := range switches {
		if v.VStatus != cloud.SubnetAvailable {
			continue
		}
		vs = append(vs, model2Switch(v))
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

func GetSwitchById(ctx context.Context, vpcId, switchId string) (Switch, error) {
	sw, err := model.FindSwitchById(ctx, vpcId, switchId)
	if err != nil {
		return Switch{}, err
	}

	return model2Switch(sw), nil
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
		ZoneId:     req.ZoneId,
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
	Ak                string
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
	vpcId := ""
	if DoesSecurityGroupBelongsVpc(vpc.Provider) {
		vpcId = vpc.VpcId
	}
	/* name 如果限制了再打开这部分
	groupIdStruct, err := model.FindSecurityId(ctx, model.FindSecurityGroupConditions{VpcId: vpcId, SecurityGroupName: req.SecurityGroupName})
	if err != nil {
		return "", errs.ErrDBQueryFailed
	}
	if groupIdStruct.GroupId != "" {
		return "", errs.ErrSecurityGroupNameExist
	}

	*/

	p, err := getProvider(vpc.Provider, req.Ak, vpc.RegionId)
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
		logs.Logger.Errorf("CreateSecurityGroup failed, %s", err.Error())
		return "", errs.ErrCreateSecurityGroupFailed
	}
	now := time.Now()
	err = model.CreateSecurityGroup(ctx, model.SecurityGroup{
		Base: model.Base{
			CreateAt: &now,
			UpdateAt: &now,
		},
		Ak:                req.Ak,
		Provider:          vpc.Provider,
		RegionId:          vpc.RegionId,
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
	Ak                string
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
	vpcId := ""
	if DoesSecurityGroupBelongsVpc(vpc.Provider) {
		vpcId = vpc.VpcId
	}
	groups, total, err := model.FindSecurityGroupWithPage(ctx, model.FindSecurityGroupConditions{
		Ak:                req.Ak,
		Provider:          vpc.Provider,
		RegionId:          vpc.RegionId,
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
	Ak              string
	RegionId        string
	VpcId           string
	SecurityGroupId string
	Rules           []GroupRule
}

type GroupRule struct {
	Protocol     string `json:"protocol"`
	PortFrom     int    `json:"port_from"`
	PortTo       int    `json:"port_to"`
	Direction    string `json:"direction"`
	GroupId      string `json:"group_id"`
	CidrIp       string `json:"cidr_ip"`
	PrefixListId string `json:"prefix_list_id"`
}

type SecurityGroupWithRule struct {
	SgId   string         `json:"security_group_id"`
	SgName string         `json:"security_group_name"`
	SgType string         `json:"security_group_type"`
	Rules  []GroupRuleRsp `json:"rules"`
}

type GroupRuleRsp struct {
	Protocol     string `json:"protocol"`
	PortRange    string `json:"port_range"`
	Direction    string `json:"direction"`
	GroupId      string `json:"group_id"`
	CidrIp       string `json:"cidr_ip"`
	PrefixListId string `json:"prefix_list_id"`
}

func AddSecurityGroupRule(ctx context.Context, req AddSecurityGroupRuleRequest) (string, error) {
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
	vpcId := ""
	if DoesSecurityGroupBelongsVpc(vpc.Provider) {
		vpcId = vpc.VpcId
	}
	cond := model.FindSecurityGroupConditions{
		Ak:              req.Ak,
		Provider:        vpc.Provider,
		RegionId:        vpc.RegionId,
		VpcId:           vpcId,
		SecurityGroupId: req.SecurityGroupId,
	}
	groupIdStruct, err := model.FindSecurityId(ctx, cond)
	if err != nil {
		logs.Logger.Errorf("FindSecurityId failed.err: [%v] req[%v]", err, req)
		return "", errs.ErrDBQueryFailed
	}
	if groupIdStruct.SecurityGroupId == "" {
		return "", errs.ErrSecurityGroupNotExist
	}

	p, err := getProvider(vpc.Provider, req.Ak, req.RegionId)
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
			PortFrom:        rule.PortFrom,
			PortTo:          rule.PortTo,
			GroupId:         rule.GroupId,
			CidrIp:          rule.CidrIp,
			PrefixListId:    rule.PrefixListId,
		}
		switch rule.Direction {
		case DirectionIn:
			err = p.AddIngressSecurityGroupRule(addRuleReq)
		case DirectionOut:
			err = p.AddEgressSecurityGroupRule(addRuleReq)
		default:
			err = fmt.Errorf("invalid direction %s", rule.Direction)
		}
		if err != nil {
			logs.Logger.Errorf("addSecurityGroupRule failed, %s", err.Error())
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
			PortRange:       getPortRange(rule.PortFrom, rule.PortTo),
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

func GetSecurityGroupWithRules(ctx context.Context, securityGroupId string) (SecurityGroupWithRule, error) {
	sg, err := model.FindSecurityGroupById(ctx, securityGroupId)
	if err != nil {
		return SecurityGroupWithRule{}, err
	}

	rules, err := model.FindSecurityGroupRulesById(ctx, securityGroupId)
	if err != nil {
		return SecurityGroupWithRule{}, err
	}

	sgRules := make([]GroupRuleRsp, 0, len(rules))
	for _, rule := range rules {
		sgRules = append(sgRules, GroupRuleRsp{
			Protocol:     rule.Protocol,
			PortRange:    rule.PortRange,
			Direction:    rule.Direction,
			GroupId:      rule.GroupId,
			CidrIp:       rule.CidrIp,
			PrefixListId: rule.PrefixListId,
		})
	}
	sgWithRule := SecurityGroupWithRule{
		SgId:   sg.SecurityGroupId,
		SgName: sg.Name,
		SgType: sg.SecurityGroupType,
		Rules:  sgRules,
	}
	return sgWithRule, nil
}

type GetRegionsRequest struct {
	Provider string
	Account  *types.OrgKeys
}

func GetRegions(ctx context.Context, req GetRegionsRequest) ([]cloud.Region, error) {
	ak := getFirstAk(req.Account, req.Provider)
	regionId := getDefaultRegion(req.Provider)
	p, err := getProvider(req.Provider, ak, regionId)
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
		logs.Logger.Errorf("GetZones failed, %s", err.Error())
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

func getDefaultRegion(provider string) string {
	regionId := ""
	switch provider {
	case cloud.AlibabaCloud:
		regionId = DefaultRegion
	case cloud.HuaweiCloud:
		regionId = DefaultRegionHuaWei
	case cloud.TencentCloud:
		regionId = DefaultRegionTencent
	}
	return regionId
}

func getPortRange(from, to int) string {
	if from < 1 {
		return ""
	}

	portRange := cast.ToString(from)
	if from != to {
		portRange = fmt.Sprintf("%d-%d", from, to)
	}
	return portRange
}
