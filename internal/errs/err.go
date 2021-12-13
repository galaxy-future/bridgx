package errs

import "errors"

var (
	ErrDBQueryFailed             = errors.New("查询数据库出错")
	ErrCreateVpcFailed           = errors.New("vpc 创建失败")
	ErrVpcNotExist               = errors.New("vpc 不存在")
	ErrCreateSwitchFailed        = errors.New("switch 创建失败")
	ErrCreateSecurityGroupFailed = errors.New("安全组创建失败")
	ErrSecurityGroupNotExist     = errors.New("安全组不存在")
	ErrGetRegionsFailed          = errors.New("获取地域信息失败")
	ErrGetZonesFailed            = errors.New("获取可用区信息失败")
	ErrVpcPending                = errors.New("pending")
	ErrBadEncryptedText          = errors.New("bad encrypted text")
	ErrSaveAccountFailed         = errors.New("save account falied")
)
