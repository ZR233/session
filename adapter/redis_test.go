/*
@Time : 2019-07-12 16:34
@Author : zr
*/
package adapter

import (
	"github.com/ZR233/session/model"
	"reflect"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

func getTestRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "192.168.0.3:6379",
		Password: "asdf*123", // no password set
		DB:       0,          // use default DB
	})
}

func TestRedis_CreateTokenMap(t *testing.T) {

	type args struct {
		token      string
		channel    string
		expireTime time.Duration
	}
	tests := []struct {
		name    string
		fields  Redis
		args    args
		wantErr bool
	}{
		{"1", newRedisForTest(), args{"111111", "test", time.Second * 10}, false},
		{"1", newRedisForTest(), args{"222222", "test", time.Second * 20}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.fields
			if err := r.CreateTokenMap(tt.args.token, tt.args.channel, tt.args.expireTime); (err != nil) != tt.wantErr {
				t.Errorf("Redis.CreateTokenMap() error = %v, wantErr %v", err, tt.wantErr)
			}
			data, err := tt.fields.db.HMGet(tt.fields.genSessionMapKey(tt.args.token), "channel", "expireAt").Result()
			if err != nil {
				t.Errorf("Redis.CreateTokenMap() error = %v", err)
			}
			if data[0] != tt.args.channel {
				t.Errorf("Redis.CreateTokenMap() %v != %v", data[0], tt.args.channel)
			}
		})
	}
}

func TestRedis_SessionUpdateUserIdAndUserTokenSetAppendToken(t *testing.T) {
	type args struct {
		userId   string
		token    string
		expireAt time.Time
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"1", args{"1", "111111", time.Now().Add(time.Second * 3)}, false},
		{"1", args{"1", "222222", time.Now().Add(time.Second * 3)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newRedisForTest()
			if err := r.SessionUpdateUserIdAndUserTokenSetAppendToken(tt.args.userId, tt.args.token, tt.args.expireAt); (err != nil) != tt.wantErr {
				t.Errorf("Redis.SessionUpdateUserIdAndUserTokenSetAppendToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRedis_FindByToken(t *testing.T) {

	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    *model.Session
		wantErr bool
	}{
		{"", args{"222222"}, &model.Session{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newRedisForTest()
			got, err := r.FindByToken(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("Redis.FindByToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Redis.FindByToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedis_FindTokenByUserId(t *testing.T) {

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"user 1", args{"1"}, []string{"123456"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newRedisForTest()
			got, err := r.FindTokenByUserId(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Redis.FindTokenByUserId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Redis.FindTokenByUserId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedis_decodeSessionMapKey(t *testing.T) {

	r1 := NewRedis(nil, "")
	r2 := NewRedis(nil, "test")
	type args struct {
		key string
	}
	tests := []struct {
		name  string
		redis Redis
		args  args
		want  string
	}{
		{"default", r1, args{r1.genSessionMapKey("123456")}, "123456"},
		{"with prefix", r2, args{r2.genSessionMapKey("123456")}, "123456"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.redis
			if got := r.decodeSessionMapKey(tt.args.key); got != tt.want {
				t.Errorf("Redis.decodeSessionMapKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedis_DeleteByToken(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		fields  Redis
		args    args
		wantErr bool
	}{
		{"12", newRedisForTest(), args{"111111"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Redis{
				db:     tt.fields.db,
				prefix: tt.fields.prefix,
			}
			if err := r.DeleteByToken(tt.args.token); (err != nil) != tt.wantErr {
				t.Errorf("Redis.DeleteByToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
