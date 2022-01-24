package authorization

import (
	"errors"
	"time"

	"github.com/galaxy-future/BridgX/config"
	"github.com/golang-jwt/jwt"
)

type CustomClaims struct {
	UserId   int64  `json:"user_id"`
	Name     string `json:"name"`
	UserType string `json:"user_type"`
	OrgId    int64  `json:"org_id"`
	jwt.StandardClaims
}

// GetOrgIdForTest 主要是在测试的时候使用.
func (c *CustomClaims) GetOrgIdForTest() int64 {
	if c == nil {
		return 0
	}
	return c.OrgId
}

func CreateUserTokenFactory() *userToken {
	return &userToken{
		userJwt: CreateJWT(config.GlobalConfig.JwtToken.JwtTokenSignKey),
	}
}

type userToken struct {
	userJwt *JwtSign
}

func (u *userToken) GenerateToken(userid int64, username string, userType string, orgId int64, expireAt int64) (tokens string, err error) {
	customClaims := CustomClaims{
		UserId:   userid,
		Name:     username,
		UserType: userType,
		OrgId:    orgId,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 10,       // 生效开始时间
			ExpiresAt: time.Now().Unix() + expireAt, // 失效截止时间
		},
	}
	return u.userJwt.CreateToken(customClaims)
}

func (u *userToken) IsValid(token string) bool {
	if customClaims, err := u.userJwt.ParseToken(token); err == nil {
		if time.Now().Unix()-(customClaims.ExpiresAt) < 0 {
			// token有效
			return true
		}
	}
	return false
}

func (u *userToken) ParseToken(tokenStr string) (*CustomClaims, error) {
	if customClaims, err := u.userJwt.ParseToken(tokenStr); err == nil {
		return customClaims, nil
	} else {
		return &CustomClaims{}, errors.New("parse token error")
	}
}

func (u *userToken) TokenIsMeetRefreshCondition(token string) bool {
	//目前仅仅判断token是否有效，后续可以考虑对于过期一段时间内的token也允许刷新
	return u.IsValid(token)
}

func (u *userToken) RefreshToken(oldToken string) (newToken string, res bool) {
	var err error
	if newToken, err = u.userJwt.RefreshToken(oldToken, config.GlobalConfig.JwtToken.JwtTokenRefreshExpires); err == nil {
		return newToken, true
	}

	return "", false
}
