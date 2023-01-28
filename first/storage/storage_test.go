package storage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gnh123/redis-action-practice/first/internal/types"
	"github.com/go-redis/redis/v9"
	"github.com/imdario/mergo"
	"github.com/stretchr/testify/assert"
)

func newRedis() *redis.Client {

	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func Test_HMset(t *testing.T) {
	r := newRedis()
	structVal := types.ArticleReq{Title: "hello", Link: "www.test.com", Poster: "guo", Time: fmt.Sprint(time.Now().Unix()), Votes: 10}
	// 结构体转map
	val := map[string]interface{}{}
	mergo.Map(&val, structVal)

	err := r.HMSet(context.TODO(), "test-hmset-get", val).Err()
	assert.NoError(t, err)
	res, err := r.HGetAll(context.TODO(), "test-hmset-get").Result()
	assert.NoError(t, err)
	fmt.Printf("%v\n", res)
}
