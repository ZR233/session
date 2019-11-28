/*
@Time : 2019-07-12 16:30
@Author : zr
*/
package session

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/ZR233/session/adapter"
	"github.com/ZR233/session/model"
	"github.com/ZR233/session/serr"
	"github.com/go-redis/redis"
	"strconv"
	"sync/atomic"
	"time"
)

var tokenLen int

func init() {
	ctx := md5.New()
	text := time.Now().String()
	ctx.Write([]byte(text))
	str := hex.EncodeToString(ctx.Sum(nil))
	tokenLen = len(str)
}

type Manager struct {
	tokenIdIter *uint64
	Prefix      string
	db          adapter.DB
}

func NewManager(db adapter.DB) *Manager {
	var iter uint64 = 0

	m := &Manager{
		tokenIdIter: &iter,
		db:          db,
	}
	return m
}

func NewRedisAdapter(client redis.UniversalClient, prefix string) adapter.Redis {
	return adapter.NewRedis(client, prefix)
}

func (m Manager) genToken() string {
	tokenIdIter := atomic.AddUint64(m.tokenIdIter, 1)
	ctx := md5.New()
	text := strconv.FormatUint(tokenIdIter, 10) + time.Now().String()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

func (m Manager) CreateSession(userId string, channel string, expireAt time.Time) (s *model.Session, err error) {
	s = &model.Session{
		UserId:  userId,
		Token:   m.genToken(),
		Channel: channel,
	}
	if err := m.db.CreateTokenMap(userId, s.Token, s.Channel, expireAt); err != nil {
		return nil, err
	}
	return s, nil
}

func (m Manager) FindByToken(token string) (s *model.Session, err error) {

	if len(token) < tokenLen {
		return nil, serr.TokenNotFound
	}
	s, err = m.db.FindByToken(token)
	if err != nil {
		return nil, err
	}

	return s, nil
}
func (m Manager) GetUserAllSessions(userId string) (sessions []*model.Session, err error) {
	return m.db.FindAllSessionsByUserId(userId)
}

func (m Manager) Update(s *model.Session) error {
	return m.db.SessionUpdate(s)
}

func (m Manager) Delete(s *model.Session) error {
	return m.db.DeleteByToken(s.Token)
}

func (m Manager) DeleteByUser(id string) error {
	sessions, err := m.db.FindAllSessionsByUserId(id)
	if err != nil {
		return err
	}
	for _, s := range sessions {
		_ = m.db.DeleteByToken(s.Token)
	}

	return nil
}
