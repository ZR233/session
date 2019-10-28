# session
开始使用

```go
import (
    "github.com/ZR233/session"
    "github.com/go-redis/redis"
    "github.com/ZR233/session/serr"
)

client := redis.NewClient(&redis.Options{
    Addr:     "localhost",
})
profix = "test_project_session"
db := session.NewRedisAdapter(redis.GetRedis(), profix)
sessionManager := session.NewManager(db)

userId := "1"
src := "pc"
expireAt := time.Now().Add(time.Hour*24*5)

//新建session
sess, err := sessionManager.CreateSession(userId, src, expireAt)
if err != nil {
    pamic(err)
}

//查找session
sess, err := sessionManager.FindByToken(token)
if err != nil{
    if err == serr.TokenNotFound {
        println(err)
    }else{
        pamic(err)
    }
}
```