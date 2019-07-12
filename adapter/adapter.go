/*
@Time : 2019-07-12 16:42
@Author : zr
*/
package adapter

import (
	"github.com/ZR233/session/session"
	"time"
)

type DB interface {
	//首次插入session, 创建key为token的map
	CreateTokenMap(token string, channel string) error
	//key为token的map 设定超时时间
	TokenMapTokenExpireAt(token string, expireAt time.Time) error
	//userId所有token列表追加token
	UserIdTokenListAppendToken(userId string, token string, expireAt time.Time) error

	FindByToken(token string) (*session.Session, error)
	//更新map的userId字段
	UpdateTokenMapSetUserId(token string, userId string) error
	//更新map的jsonField字段
	UpdateTokenMapSetJsonField(token string, jsonField string) error

	FindTokenByUserId(token string) ([]string, error)
}
