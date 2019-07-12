/*
@Time : 2019-07-12 16:34
@Author : zr
*/
package adapter

import (
	"github.com/ZR233/session/session"
	"github.com/go-redis/redis"
	"time"
)

type Redis struct {
	db *redis.Client
}

func (r Redis) CreateTokenMap(token string, channel string) error {
	panic("implement me")
}

func (Redis) TokenMapTokenExpireAt(token string, expireAt time.Time) error {
	panic("implement me")
}

func (Redis) UserIdTokenListAppendToken(userId string, token string, expireAt time.Time) error {
	panic("implement me")
}

func (Redis) FindByToken(token string) (*session.Session, error) {
	panic("implement me")
}

func (Redis) UpdateTokenMapSetUserId(token string, userId string) error {
	panic("implement me")
}

func (Redis) UpdateTokenMapSetJsonField(token string, jsonField string) error {
	panic("implement me")
}

func (Redis) FindTokenByUserId(token string) ([]string, error) {
	panic("implement me")
}

func NewRedis(client *redis.Client) Redis {
	a := Redis{
		client,
	}
	return a
}
