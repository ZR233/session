/*
@Time : 2019-07-12 16:34
@Author : zr
*/
package adapter

import (
	"errors"
	"github.com/ZR233/session/model"
	"github.com/ZR233/session/serr"
	"github.com/go-redis/redis"
	"time"
)

const (
	DefaultPrefix = "model"
)

type Redis struct {
	db     *redis.Client
	prefix string
}

func (r *Redis) genSessionMapKey(token string) string {
	if r.prefix == "" {
		r.prefix = DefaultPrefix
	}
	key := r.prefix + "_token_" + token
	return key
}
func (r *Redis) decodeSessionMapKey(key string) string {
	if r.prefix == "" {
		r.prefix = DefaultPrefix
	}
	key = key[len(r.prefix)+7:]
	return key
}

func (r Redis) genUserSessionSetKey(userId string) string {
	if r.prefix == "" {
		r.prefix = DefaultPrefix
	}
	key := r.prefix + "_user_" + userId
	return key
}

func (r Redis) CreateTokenMap(token string, channel string, expireTime time.Duration) error {
	key := r.genSessionMapKey(token)
	values := make(map[string]interface{})
	values["channel"] = channel
	expireAt := time.Now().Add(expireTime)
	pipe := r.db.TxPipeline()
	pipe.HMSet(key, values)
	pipe.ExpireAt(key, expireAt)
	_, err := pipe.Exec()

	return err
}

func (r Redis) TokenMapTokenExpireAt(token string, expireAt time.Time) error {
	tokenKey := r.genSessionMapKey(token)
	return r.db.ExpireAt(tokenKey, expireAt).Err()
}

func (r Redis) SessionUpdate(s *model.Session) error {
	userKey := r.genUserSessionSetKey(s.UserId)
	tokenKey := r.genSessionMapKey(s.Token)

	values := make(map[string]interface{})
	values["userid"] = s.UserId

	pipe := r.db.TxPipeline()
	pipe.HMSet(tokenKey, values)
	pipe.SAdd(userKey, tokenKey)
	pipe.ExpireAt(tokenKey, s.ExpireAt)
	pipe.ExpireAt(userKey, s.ExpireAt)
	_, err := pipe.Exec()
	return err
}

func (r Redis) FindByToken(token string) (*model.Session, error) {
	tokenKey := r.genSessionMapKey(token)
	data, err := r.db.HGetAll(tokenKey).Result()
	if err != nil {
		return nil, err
	}
	expire, err := r.db.TTL(tokenKey).Result()
	if err != nil {
		return nil, err
	}

	timestamp := time.Now().Add(expire)

	s := &model.Session{
		Token:    token,
		UserId:   data["userid"],
		Channel:  data["channel"],
		ExpireAt: timestamp,
	}

	if s.UserId == "" {
		return nil, serr.NewErr(errors.New("token not found"), serr.TokenNotFind)
	}
	return s, nil
}

func (r Redis) UpdateTokenMapSetJsonField(token string, jsonField string) error {
	tokenKey := r.genSessionMapKey(token)
	return r.db.HSet(tokenKey, "jsonField", jsonField).Err()
}

func (r Redis) FindTokenByUserId(id string) ([]string, error) {
	userKey := r.genUserSessionSetKey(id)
	data, err := r.db.SMembers(userKey).Result()
	if err != nil {
		return nil, err
	}
	var tokens []string
	for _, v := range data {
		v = r.decodeSessionMapKey(v)
		tokens = append(tokens, v)
	}

	return data, nil
}

func (r Redis) DeleteByToken(token string) error {
	tokenKey := r.genSessionMapKey(token)
	data, err := r.db.HGetAll(tokenKey).Result()
	if err != nil {
		return err
	}
	if userid, ok := data["userid"]; ok {
		if err := r.db.SRem(r.genUserSessionSetKey(userid), tokenKey).Err(); err != nil {
			return nil
		}
	}

	return r.db.Del(tokenKey).Err()
}

func newRedisForTest() Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     "192.168.0.3:6379",
		Password: "asdf*123", // no password set
		DB:       0,          // use default DB
	})

	a := Redis{
		client,
		"test_session",
	}
	return a
}

func NewRedis(client *redis.Client, prefix string) Redis {
	a := Redis{
		client,
		prefix,
	}
	return a
}
