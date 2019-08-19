/*
@Time : 2019-07-12 16:42
@Author : zr
*/
package adapter

import (
	"github.com/ZR233/session/model"
	"time"
)

type DB interface {
	//首次插入session, 创建key为token的map
	CreateTokenMap(userId string, token string, channel string, expireAt time.Time) error

	//更新
	SessionUpdate(*model.Session) error

	//通过token查找session
	FindByToken(token string) (*model.Session, error)

	//通过UserId找到所有token
	FindAllSessionsByUserId(id string) (sessions []*model.Session, err error)

	//根据token删除session
	DeleteByToken(token string) error
}
