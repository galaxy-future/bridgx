package model

import (
	"context"
	"strings"
	"time"

	"gorm.io/gorm/clause"

	"github.com/galaxy-future/BridgX/internal/clients"
	"github.com/galaxy-future/BridgX/internal/constants"
	"github.com/galaxy-future/BridgX/internal/logs"
)

type Network struct {
	Base
	Ak                      string
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
	Ak        string
	RegionId  string
	VpcId     string
	Name      string
	CidrBlock string
	SwitchIds string
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
	Ak         string
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
		Table("b_vpc").
		Where("vpc_id = ? and is_del = 0", cond.VpcId).
		First(&result).
		Error
	return result, nil
}

func FindVpcsWithPage(ctx context.Context, cond FindVpcConditions) (result []Vpc, total int64, err error) {
	query := clients.ReadDBCli.WithContext(ctx).Table(Vpc{}.TableName()).Where("ak = ? and is_del = 0", cond.Ak)
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

func UpdateVpc(ctx context.Context, vpcId, cidrBlock, vStatus string, switchIds []string) error {
	now := time.Now()
	queryMap := map[string]interface{}{
		"cidr_block": cidrBlock,
		"v_status":   vStatus,
		"update_at":  &now,
	}
	if len(switchIds) > 0 {
		queryMap["switch_ids"] = strings.Join(switchIds, ",")
	}
	return clients.WriteDBCli.WithContext(ctx).
		Table(Vpc{}.TableName()).
		Where(`vpc_id = ?`, vpcId).
		Updates(queryMap).
		Error
}

type FindSwitchesConditions struct {
	VpcId      string
	SwitchId   string
	SwitchName string
	PageNumber int
	PageSize   int
}

func FindSwitchesWithPage(ctx context.Context, cond FindSwitchesConditions) (result []Switch, total int64, err error) {
	query := clients.ReadDBCli.WithContext(ctx).Table(Switch{}.TableName()).Where("vpc_id = ? and is_del = 0", cond.VpcId)
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
	VpcId             string
	SecurityGroupId   string
	SecurityGroupName string
	PageNumber        int
	PageSize          int
}

func FindSecurityGroupWithPage(ctx context.Context, cond FindSecurityGroupConditions) (result []SecurityGroup, total int64, err error) {
	query := clients.ReadDBCli.WithContext(ctx).Table(SecurityGroup{}.TableName()).Where("vpc_id = ? and is_del = 0", cond.VpcId)
	if cond.SecurityGroupId != "" {
		query.Where("switch_id = ?", cond.SecurityGroupId)
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
	err = clients.ReadDBCli.WithContext(ctx).
		Table("b_security_group").
		Select(`security_group_id`).
		Where("vpc_id = ? and security_group_id = ? and is_del = 0", cond.VpcId, cond.SecurityGroupId).
		Scan(&result).
		Error
	return result, nil
}

func CreateSecurityGroup(ctx context.Context, s SecurityGroup) error {
	return clients.WriteDBCli.WithContext(ctx).Create(&s).Error
}

func AddSecurityGroupRule(ctx context.Context, r SecurityGroupRule) error {
	return clients.WriteDBCli.WithContext(ctx).Create(&r).Error
}

func UpdateOrCreateVpcs(ctx context.Context, vpcs []Vpc) error {
	return clients.WriteDBCli.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "vpc_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "cidr_block", "switch_ids", "v_status"}),
	}).Create(&vpcs).Error
}

func UpdateOrCreateSwitches(ctx context.Context, switches []Switch) error {
	return clients.WriteDBCli.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "switch_id"}, {Name: "vpc_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "cidr_block", "v_status", "available_ip_address_count", "is_default"}),
	}).Create(&switches).Error
}

func UpdateOrCreateGroups(ctx context.Context, groups []SecurityGroup) error {
	return clients.WriteDBCli.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "security_group_id"}, {Name: "vpc_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name"}),
	}).Create(&groups).Error
}

func ReplaceRules(ctx context.Context, vpcID, groupId string, rules []SecurityGroupRule) error {
	tx := clients.WriteDBCli.WithContext(ctx).Begin()
	defer func() {
		tx.Rollback()
	}()
	err := tx.WithContext(ctx).
		Where("security_group_id = ? and vpc_id = ?", groupId, vpcID).
		Delete(SecurityGroupRule{}).Error
	if err != nil {
		tx.Rollback()
	}
	err = tx.CreateInBatches(&rules, len(rules)).Error
	if err != nil {
		tx.Rollback()
	}
	tx.Commit()
	return err
}
