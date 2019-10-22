/*
@Time : 2019-07-12 16:30
@Author : zr
*/
package session

import (
	"github.com/ZR233/session/adapter"
	"github.com/ZR233/session/serr"
	"github.com/go-redis/redis"
	"testing"
	"time"
)

func TestIndex_All(t *testing.T) {
	db := redis.NewClient(&redis.Options{
		Addr:     "192.168.0.3:6379",
		Password: "asdf*123", // no password set
		DB:       0,          // use default DB
	})

	redisClient := adapter.NewRedis(db, "test")

	m := NewManager(redisClient)
	userId := "1"
	src := "a"
	expire := time.Now().Add(time.Second)

	session, err := m.CreateSession(userId, src, expire)
	if err != nil {
		t.Error(err)
	}

	sessionFound, err := m.FindByToken(session.Token)
	if err != nil {
		t.Error(err)
	}

	if session.Token != sessionFound.Token {
		t.Error("session不一致")
	}
	if session.UserId != sessionFound.UserId {
		t.Error("session不一致")
	}
	if session.Channel != sessionFound.Channel {
		t.Error("session不一致")
	}

	sessions, err := m.GetUserAllSessions("1")
	if err != nil {
		t.Error(err)
	}
	if len(sessions) != 1 {
		t.Error("map 错误")
	}

	time.Sleep(time.Second)

	sessionFound, err = m.FindByToken(session.Token)

	if err != nil {
		if err != serr.TokenNotFound {
			t.Error(err)
		}
	} else {
		t.Error("token 未删除")
	}

}
