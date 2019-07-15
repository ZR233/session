/*
@Time : 2019-07-12 16:42
@Author : zr
*/
package adapter

import (
	session2 "github.com/ZR233/session"
	"time"
)

type DB interface {
	//首次插入session, 创建key为token的map
	CreateTokenMap(token string, channel string, expireTime time.Duration) error
	//key为token的map 设定超时时间
	TokenMapTokenExpireAt(token string, expireAt time.Time) error
	//session 设置userId, userId所有token列表追加token
	SessionUpdateUserIdAndUserTokenSetAppendToken(userId string, token string, expireAt time.Time) error

	FindByToken(token string) (*session2.Session, error)

	//更新map的jsonField字段
	UpdateTokenMapSetJsonField(token string, jsonField string) error

	FindTokenByUserId(id string) ([]string, error)

	DeleteByToken(token string) error
}
