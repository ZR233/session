/*
@Time : 2019-07-12 16:34
@Author : zr
*/
package adapter

import (
	"github.com/ZR233/session/model"
	"github.com/go-redis/redis"
	"reflect"
	"testing"
	"time"
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
		userid   string
		token    string
		channel  string
		expireAt time.Time
	}
	tests := []struct {
		name    string
		fields  Redis
		args    args
		wantErr bool
	}{
		{"1", newRedisForTest(), args{"1", "111111", "test",
			time.Now().Add(time.Second * 10)}, false},
		{"1", newRedisForTest(), args{"1", "222222", "test",
			time.Now().Add(time.Second * 20)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.fields
			if err := r.CreateTokenMap(tt.args.userid, tt.args.token, tt.args.channel, tt.args.expireAt); (err != nil) != tt.wantErr {
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

func TestRedis_FindByToken(t *testing.T) {

	r := newRedisForTest()
	at := time.Now().Add(time.Second * 5)
_:
	r.CreateTokenMap("99", "1234", "test", at)

	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    *model.Session
		wantErr bool
	}{
		{"未找到token", args{"222222"}, nil, true},
		{"找到", args{"1234"}, &model.Session{
			Token:      "1234",
			UserId:     "99",
			Channel:    "test",
			ExpireAt:   at,
			JsonFields: nil,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := r.FindByToken(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("Redis.FindByToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if got.Token != tt.want.Token ||
					got.UserId != tt.want.UserId ||
					got.Channel != tt.want.Channel ||
					got.JsonFields != tt.want.JsonFields {
					t.Errorf("Redis.FindByToken() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestRedis_FindAllSessionsByUserId(t *testing.T) {

	r := newRedisForTest()
	at1 := time.Now().Add(time.Second * 2)
	at2 := time.Now().Add(time.Millisecond * 2)
_:
	r.CreateTokenMap("98", "1aaaaa", "test", at1)
_:
	r.CreateTokenMap("98", "2bbbbb", "test", at2)

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"1超时1正常", args{"98"}, []string{"1aaaaa"}, false},
	}

	time.Sleep(time.Millisecond * 50)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newRedisForTest()
			got, err := r.FindAllSessionsByUserId(tt.args.id)
			var got2 []string

			for _, g := range got {
				got2 = append(got2, g.Token)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("Redis.FindAllSessionsByUserId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got2, tt.want) {
				t.Errorf("Redis.FindAllSessionsByUserId() = %v, want %v", got, tt.want)
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
