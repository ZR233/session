/*
@Time : 2019-07-12 16:34
@Author : zr
*/
package adapter

import (
	"github.com/ZR233/session/session"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

const (
	DefaultPrefix = "session"
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
	values["expireAt"] = time.Now().Add(expireTime).Unix()
	return r.db.HMSet(key, values).Err()
}

func (r Redis) TokenMapTokenExpireAt(token string, expireAt time.Time) error {
	tokenKey := r.genSessionMapKey(token)
	return r.db.HSet(tokenKey, "expireAt", expireAt.Unix()).Err()
}

func (r Redis) SessionUpdateUserIdAndUserTokenSetAppendToken(userId string, token string, expireAt time.Time) error {
	userKey := r.genUserSessionSetKey(userId)
	tokenKey := r.genSessionMapKey(token)

	values := make(map[string]interface{})
	values["userid"] = userId
	values["expireAt"] = expireAt.Unix()

	pipe := r.db.TxPipeline()
	pipe.HMSet(tokenKey, values)
	pipe.SAdd(userKey, tokenKey)
	_, err := pipe.Exec()
	return err
}

func (r Redis) FindByToken(token string) (*session.Session, error) {
	tokenKey := r.genSessionMapKey(token)
	data, err := r.db.HGetAll(tokenKey).Result()
	if err != nil {
		return nil, err
	}

	timestanp, _ := strconv.ParseInt(data["expireAt"], 10, 64)
	if timestanp == 0 {
		return nil, nil
	}
	expireAt := time.Unix(timestanp, 0)

	if expireAt.Sub(time.Now()) < 0 {
		if err := r.DeleteByToken(token); err != nil {
			return nil, err
		}
		return nil, nil
	}

	s := &session.Session{
		Token:    token,
		UserId:   data["userid"],
		Channel:  data["channel"],
		ExpireAt: time.Unix(timestanp, 0),
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
