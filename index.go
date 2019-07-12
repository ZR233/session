/*
@Time : 2019-07-12 16:30
@Author : zr
*/
package session

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/ZR233/session/adapter"
	"github.com/ZR233/session/serr"
	"github.com/ZR233/session/session"
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
	m := &Manager{}
	return m
}
func (m Manager) genToken() string {
	tokenIdIter := atomic.AddUint64(m.tokenIdIter, 1)
	ctx := md5.New()
	text := strconv.FormatUint(tokenIdIter, 10) + time.Now().String()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

func (m Manager) CreateSession(channel string, expireTime time.Duration) (s *session.Session, err error) {
	s = &session.Session{
		Token:   m.genToken(),
		Channel: channel,
	}
	if err := m.db.CreateTokenMap(s.Token, s.Channel); err != nil {
		return nil, err
	}

	return s, nil
}
func (m Manager) UpdateUserId(s *session.Session, token string, expireTime time.Duration) (*session.Session, error) {
	s.Token = token

	if err := m.db.UpdateTokenMapSetUserId(token, s.UserId); err != nil {
		return nil, err
	}
	expireAt := time.Now().Add(expireTime)

	if err := m.db.UserIdTokenListAppendToken(s.UserId, token, expireAt); err != nil {
		return nil, err
	}
	return s, nil
}

func (m Manager) FindByToken(token string) (s *session.Session, err error) {

	if len(token) < tokenLen {
		return nil, serr.NewErr(errors.New("token not found"), serr.TokenNotFind)
	}
	s, err = m.db.FindByToken(token)
	if err != nil {
		return nil, err
	}
	return s, nil
}
func (m Manager) GetUserAllSessions(userId string) (sessions []*session.Session, err error) {
	tokens, err := m.db.FindTokenByUserId(userId)
	if err != nil {
		return sessions, serr.NewErr(err, serr.RedisErr)
	}
	for _, v := range tokens {
		s, _ := m.FindByToken(v)
		sessions = append(sessions, s)
	}
	return sessions, nil
}

func (m Manager) UpdateJsonField(s *session.Session, jsonField interface{}) error {
	jsonStr, err := json.Marshal(jsonField)
	if err != nil {
		return serr.NewErr(err, serr.JsonErr)
	}

	return m.db.UpdateTokenMapSetJsonField(s.Token, string(jsonStr))
}
