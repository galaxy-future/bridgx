package model

import (
	"context"
	"time"

	"github.com/galaxy-future/BridgX/internal/clients"
	"github.com/galaxy-future/BridgX/internal/constants"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/pkg/utils"
	"gorm.io/gorm/clause"
)

type Network struct {
	Base
	AK                      string `gorm:"ak"`
	RegionId                string
	VpcId                   string
	SubNetId                string
	SecurityGroup           string
	InternetChargeType      string
	InternetMaxBandwidthOut string
}

func (Network) TableName() string {
	return "b_network"
}

type Vpc struct {
	Base
	AK        string `gorm:"ak"`
	RegionId  string
	VpcId     string
	Name      string
	CidrBlock string
	Provider  string
	VStatus   string
	IsDel     int
}

func (Vpc) TableName() string {
	return "b_vpc"
}

type Switch struct {
	Base
	VpcId                   string
	SwitchId                string
	ZoneId                  string
	Name                    string
	CidrBlock               string
	GatewayIp               string
	IsDefault               int
	AvailableIpAddressCount int
	VStatus                 string
	IsDel                   int
}

func (Switch) TableName() string {
	return "b_switch"
}

type SecurityGroup struct {
	Base
	AK                string `gorm:"ak"`
	Provider          string
	RegionId          string
	VpcId             string
	SecurityGroupId   string
	Name              string
	SecurityGroupType string
	IsDel             int
}

func (SecurityGroup) TableName() string {
	return "b_security_group"
}

type SecurityGroupRule struct {
	Base
	VpcId           string
	SecurityGroupId string
	PortRange       string
	Protocol        string
	Direction       string
	GroupId         string `gorm:"column:other_group_id"`
	CidrIp          string
	PrefixListId    string
	IsDel           int
}

func (SecurityGroupRule) TableName() string {
	return "b_security_group_rule"
}

type FindVpcConditions struct {
	AccountKey string
	VpcId      string
	VpcName    string
	RegionId   string
	PageNumber int
	PageSize   int
	Provider   string
}

type VpcIDStruct struct {
	VpcId string
}

type SwitchIdStruct struct {
	SwitchId string
}

type SecurityGroupIDStruct struct {
	SecurityGroupId string
}

func FindVpcById(ctx context.Context, cond FindVpcConditions) (result Vpc, err error) {
	err = clients.ReadDBCli.WithContext(ctx).
		Where("vpc_id = ? and is_del = 0", cond.VpcId).
		First(&result).
		Error
	return result, err
}

func FindVpcsWithPage(ctx context.Context, cond FindVpcConditions) (result []Vpc, total int64, err error) {
	query := clients.ReadDBCli.WithContext(ctx).Table(Vpc{}.TableName()).Where("ak = ? and is_del = 0", cond.AccountKey)
	if cond.RegionId != "" {
		query.Where("region_id = ?", cond.RegionId)
	}
	if cond.VpcName != "" {
		query.Where("name = ?", cond.VpcName)
	}
	if cond.PageNumber <= 0 {
		cond.PageNumber = 1
	}
	if cond.PageSize <= 0 || cond.PageSize > constants.DefaultPageSize {
		cond.PageSize = constants.DefaultPageSize
	}
	offset := (cond.PageNumber - 1) * cond.PageSize
	err = query.Find(&result).Limit(int(cond.PageSize)).Offset(int(offset)).Error
	if err != nil {
		logs.Logger.Errorf("FindVpcsWithPage failed.err: [%v]", err)
		return nil, 0, err
	}
	err = query.Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logs.Logger.Errorf("FindVpcsWithPage failed.err: [%v]", err)
		return nil, 0, err
	}
	return result, total, nil
}

func CreateVpc(ctx context.Context, vpc Vpc) error {
	return clients.WriteDBCli.WithContext(ctx).Create(&vpc).Error
}

func UpdateVpc(ctx context.Context, vpcId, cidrBlock, vStatus string) error {
	now := time.Now()
	queryMap := map[string]interface{}{
		"cidr_block": cidrBlock,
		"v_status":   vStatus,
		"update_at":  &now,
	}

	return clients.WriteDBCli.WithContext(ctx).
		Table(Vpc{}.TableName()).
		Where(`vpc_id = ?`, vpcId).
		Updates(queryMap).
		Error
}

type FindSwitchesConditions struct {
	VpcId      string
	ZoneId     string
	SwitchId   string
	SwitchName string
	PageNumber int
	PageSize   int
}

func FindSwitchesWithPage(ctx context.Context, cond FindSwitchesConditions) (result []Switch, total int64, err error) {
	query := clients.ReadDBCli.WithContext(ctx).Table(Switch{}.TableName()).Where("vpc_id = ? and zone_id = ? and is_del = 0", cond.VpcId, cond.ZoneId)
	if cond.SwitchId != "" {
		query.Where("switch_id = ?", cond.SwitchId)
	}
	if cond.SwitchName != "" {
		query.Where("name = ?", cond.SwitchName)
	}
	if cond.PageNumber <= 0 {
		cond.PageNumber = 1
	}
	if cond.PageSize <= 0 || cond.PageSize > constants.DefaultPageSize {
		cond.PageSize = constants.DefaultPageSize
	}
	offset := (cond.PageNumber - 1) * cond.PageSize
	query = query.Find(&result).Limit(int(cond.PageSize)).Offset(int(offset))
	err = query.Error
	if err != nil {
		logs.Logger.Errorf("FindSwitchesWithPage failed.err: [%v]", err)
		return nil, 0, err
	}
	err = query.Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		return result, 0, err
	}
	return result, total, nil
}

func FindSwitchId(ctx context.Context, cond FindSwitchesConditions) (result SwitchIdStruct, err error) {
	err = clients.ReadDBCli.WithContext(ctx).
		Table("b_switch").
		Select(`switch_id`).
		Where("vpc_id =? and name = ? and is_del = 0", cond.VpcId, cond.SwitchName).
		Scan(&result).
		Error
	return result, nil
}

func FindSwitchById(ctx context.Context, vpcId, switchId string) (result Switch, err error) {
	err = clients.ReadDBCli.WithContext(ctx).
		Where("vpc_id =? and switch_id = ? and is_del = 0", vpcId, switchId).
		First(&result).
		Error
	return result, err
}

func CreateSwitch(ctx context.Context, s Switch) error {
	return clients.WriteDBCli.WithContext(ctx).Create(&s).Error
}

func UpdateSwitch(ctx context.Context, availableIpAddressCount, isDefault int, vpcId, switchId, name, sStatus, cidrBlock, gatewayIp string) error {
	now := time.Now()
	queryMap := map[string]interface{}{
		"available_ip_address_count": availableIpAddressCount,
		"is_default":                 isDefault,
		"v_status":                   sStatus,
		"name":                       name,
		"cidr_block":                 cidrBlock,
		"gateway_ip":                 gatewayIp,
		"update_at":                  &now,
	}
	return clients.WriteDBCli.WithContext(ctx).
		Table(Switch{}.TableName()).
		Where(`vpc_id = ? and switch_id = ?`, vpcId, switchId).
		Updates(queryMap).
		Error
}

type FindSecurityGroupConditions struct {
	AK                string
	Provider          string
	RegionId          string
	VpcId             string
	SecurityGroupId   string
	SecurityGroupName string
	PageNumber        int
	PageSize          int
}

func FindSecurityGroupWithPage(ctx context.Context, cond FindSecurityGroupConditions) (result []SecurityGroup, total int64, err error) {
	query := clients.ReadDBCli.WithContext(ctx).
		Model(&SecurityGroup{}).
		Where("ak = ? and provider = ? and region_id=? and is_del = 0", cond.AK, cond.Provider, cond.RegionId)
	if cond.VpcId != "" {
		query.Where("vpc_id = ?", cond.VpcId)
	}
	if cond.SecurityGroupId != "" {
		query.Where("security_group_id = ?", cond.SecurityGroupId)
	}
	if cond.SecurityGroupName != "" {
		query.Where("name = ?", cond.SecurityGroupName)
	}
	if cond.PageNumber <= 0 {
		cond.PageNumber = 1
	}
	if cond.PageSize <= 0 || cond.PageSize > constants.DefaultPageSize {
		cond.PageSize = constants.DefaultPageSize
	}
	offset := (cond.PageNumber - 1) * cond.PageSize
	query = query.Find(&result).Limit(int(cond.PageSize)).Offset(int(offset))
	err = query.Error
	if err != nil {
		logs.Logger.Errorf("FindSecurityGroupWithPage failed.err: [%v]", err)
		return nil, 0, err
	}
	err = query.Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logs.Logger.Errorf("FindSecurityGroupWithPage failed.err: [%v]", err)
		return nil, 0, err
	}
	return result, total, nil
}

func FindSecurityId(ctx context.Context, cond FindSecurityGroupConditions) (result SecurityGroupIDStruct, err error) {
	sql := clients.ReadDBCli.WithContext(ctx).
		Model(&SecurityGroup{}).
		Select(`security_group_id`).
		Where("ak = ? and provider = ? and region_id=? and is_del = 0", cond.AK, cond.Provider, cond.RegionId)
	if cond.VpcId != "" {
		sql.Where("vpc_id = ?", cond.VpcId)
	}
	if cond.SecurityGroupId != "" {
		sql.Where("security_group_id = ?", cond.SecurityGroupId)
	}
	if cond.SecurityGroupName != "" {
		sql.Where("name = ?", cond.SecurityGroupName)
	}
	err = sql.Scan(&result).Error
	return result, err
}

func FindSecurityGroupById(ctx context.Context, securityGroupId string) (result SecurityGroup, err error) {
	err = clients.ReadDBCli.WithContext(ctx).
		Where("security_group_id = ? and is_del = 0", securityGroupId).
		First(&result).Error
	return result, err
}

func CreateSecurityGroup(ctx context.Context, s SecurityGroup) error {
	return clients.WriteDBCli.WithContext(ctx).Create(&s).Error
}

func AddSecurityGroupRule(ctx context.Context, r SecurityGroupRule) error {
	return clients.WriteDBCli.WithContext(ctx).Create(&r).Error
}

const _effectiveTime = "DATE_ADD(now(),interval 478 minute)"

func UpdateOrCreateVpcs(ctx context.Context, ak, provider string, regionIds []string, vpcs []Vpc) error {
	oldVpcIds := make([]string, 0)
	if err := clients.ReadDBCli.WithContext(ctx).Model(&Vpc{}).
		Select("vpc_id").
		Where("ak=? and provider=? and region_id in (?) and is_del=0 and update_at<?", ak, provider, regionIds, _effectiveTime).
		Scan(&oldVpcIds).Error; err != nil {
		return err
	}
	vpcIds := make([]string, 0, len(vpcs))
	for _, v := range vpcs {
		vpcIds = append(vpcIds, v.VpcId)
	}
	vpcIdDiff := utils.StringSliceDiff(oldVpcIds, vpcIds)
	if len(vpcIdDiff) > 0 {
		logs.Logger.Debugf("%s, vpcIdDiff %v", ak, vpcIdDiff)
		if err := clients.WriteDBCli.WithContext(ctx).
			Where("ak=? and provider=? and vpc_id in (?)", ak, provider, vpcIdDiff).
			Delete(&Vpc{}).Error; err != nil {
			return err
		}
	}

	return clients.WriteDBCli.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "ak"}, {Name: "region_id"}, {Name: "vpc_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "cidr_block", "v_status"}),
	}).Create(&vpcs).Error
}

func UpdateOrCreateSwitches(ctx context.Context, vpcIds []string, switches []Switch) error {
	oldSwitchIds := make([]string, 0)
	if err := clients.ReadDBCli.WithContext(ctx).Model(&Switch{}).
		Select("switch_id").
		Where("vpc_id in (?) and is_del=0 and update_at<?", vpcIds, _effectiveTime).
		Scan(&oldSwitchIds).Error; err != nil {
		return err
	}
	switchIds := make([]string, 0, len(switches))
	for _, v := range switches {
		switchIds = append(switchIds, v.SwitchId)
	}
	switchIdDiff := utils.StringSliceDiff(oldSwitchIds, switchIds)
	if len(switchIdDiff) > 0 {
		logs.Logger.Debugf("switchIdDiff %v", switchIdDiff)
		if err := clients.WriteDBCli.WithContext(ctx).
			Where("switch_id in (?)", switchIdDiff).
			Delete(&Switch{}).Error; err != nil {
			return err
		}
	}

	return clients.WriteDBCli.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "switch_id"}, {Name: "vpc_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "cidr_block", `gateway_ip`, "v_status", "available_ip_address_count", "is_default"}),
	}).Create(&switches).Error
}

func UpdateOrCreateGroups(ctx context.Context, ak, provider string, regionIds []string, groups []SecurityGroup) error {
	oldSecGrpIds := make([]string, 0)
	if err := clients.ReadDBCli.WithContext(ctx).Model(&SecurityGroup{}).
		Select("security_group_id").
		Where("ak=? and provider=? and region_id in (?) and is_del=0 and update_at<?", ak, provider, regionIds, _effectiveTime).
		Scan(&oldSecGrpIds).Error; err != nil {
		return err
	}
	secGrpIds := make([]string, 0, len(groups))
	for _, v := range groups {
		secGrpIds = append(secGrpIds, v.SecurityGroupId)
	}
	secGrpIdDiff := utils.StringSliceDiff(oldSecGrpIds, secGrpIds)
	if len(secGrpIdDiff) > 0 {
		logs.Logger.Debugf("%s, secGrpIdDiff %v", ak, secGrpIdDiff)
		if err := clients.WriteDBCli.WithContext(ctx).
			Where("ak=? and provider=? and security_group_id in (?)", ak, provider, secGrpIdDiff).
			Delete(&SecurityGroup{}).Error; err != nil {
			return err
		}
	}

	return clients.WriteDBCli.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "ak"}, {Name: "provider"}, {Name: "region_id"}, {Name: "security_group_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "security_group_type"}),
	}).Create(&groups).Error
}

func ReplaceRules(ctx context.Context, vpcID, groupId string, rules []SecurityGroupRule) (err error) {
	tx := clients.WriteDBCli.WithContext(ctx).Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	err = tx.WithContext(ctx).
		Where("security_group_id = ? and vpc_id = ?", groupId, vpcID).
		Delete(SecurityGroupRule{}).Error
	if err != nil {
		return err
	}
	err = tx.CreateInBatches(&rules, len(rules)).Error
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func FindSecurityGroupRulesById(ctx context.Context, securityGroupId string) (result []SecurityGroupRule, err error) {
	err = clients.ReadDBCli.WithContext(ctx).
		Where("security_group_id = ? and is_del = 0", securityGroupId).
		Find(&result).Error
	return result, err
}
