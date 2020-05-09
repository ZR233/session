/*
@Time : 2019-07-12 16:34
@Author : zr
*/
package adapter

import (
	"encoding/json"
	"github.com/ZR233/session/model"
	"github.com/ZR233/session/serr"
	"github.com/go-redis/redis/v7"
	"strings"
	"time"
)

const (
	sessionPrefix = "session"
)

type Redis struct {
	db     redis.UniversalClient
	prefix string
}

func (r *Redis) genSessionMapKey(token string) string {
	key := strings.Join([]string{
		r.prefix,
		"token",
		token,
	}, ":")

	return key
}
func (r *Redis) decodeSessionMapKey(key string) string {
	key = key[len(r.prefix)+7:]
	return key
}

func (r Redis) genUserSessionSetKey(userId string) string {
	key := r.prefix + ":user:" + userId
	return key
}

func (r Redis) CreateTokenMap(userId string, token string, channel string, expireAt time.Time) error {
	key := r.genSessionMapKey(token)
	userKey := r.genUserSessionSetKey(userId)

	values := make(map[string]interface{})
	values["userid"] = userId
	values["channel"] = channel

	pipe := r.db.TxPipeline()
	//创建HashSet token-values
	pipe.HMSet(key, values)
	pipe.ExpireAt(key, expireAt)

	//userId-tokenList  userId添加token
	pipe.SAdd(userKey, key)
	setTTL, _ := pipe.TTL(userKey).Result()
	if setTTL < expireAt.Sub(time.Now()) {
		pipe.ExpireAt(userKey, expireAt)
	}

	_, err := pipe.Exec()

	return err
}

func (r Redis) TokenMapTokenExpireAt(token string, expireAt time.Time) error {
	tokenKey := r.genSessionMapKey(token)
	return r.db.ExpireAt(tokenKey, expireAt).Err()
}

func (r Redis) SessionUpdate(s *model.Session) error {

	tokenKey := r.genSessionMapKey(s.Token)
	userKey := r.genUserSessionSetKey(s.UserId)
	values := make(map[string]interface{})
	values["userid"] = s.UserId
	values["channel"] = s.Channel
	jsonStr, err := json.Marshal(s.JsonFields)
	if err != nil {
		return err
	}

	pipe := r.db.TxPipeline()
	pipe.HMSet(tokenKey, values)
	pipe.ExpireAt(tokenKey, s.ExpireAt)
	pipe.HSet(tokenKey, "jsonField", jsonStr)

	setTTL, _ := pipe.TTL(userKey).Result()
	if setTTL < s.ExpireAt.Sub(time.Now()) {
		pipe.ExpireAt(userKey, s.ExpireAt)
	}
	_, err = pipe.Exec()
	return err
}

func (r Redis) FindByToken(token string) (*model.Session, error) {
	tokenKey := r.genSessionMapKey(token)
	data, err := r.db.HGetAll(tokenKey).Result()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, serr.TokenNotFound
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

	return s, nil
}

func (r Redis) FindAllSessionsByUserId(id string) (sessions []*model.Session, err error) {
	userKey := r.genUserSessionSetKey(id)

	data, err := r.db.SMembers(userKey).Result()
	if err != nil {
		return nil, err
	}
	var tokensNotExist []string
	for _, v := range data {
		v = r.decodeSessionMapKey(v)
		var s *model.Session
		s, err = r.FindByToken(v)
		if err != nil {
			if err == serr.TokenNotFound {
				tokensNotExist = append(tokensNotExist, v)
				err = nil
				continue
			} else {
				return
			}
		}
		sessions = append(sessions, s)
	}
	r.db.SRem(userKey, tokensNotExist)
	return
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

func NewRedis(options *redis.UniversalOptions, prefix string) Redis {
	client := redis.NewUniversalClient(options)

	a := Redis{
		client,
		prefix + ":" + sessionPrefix,
	}
	return a
}
