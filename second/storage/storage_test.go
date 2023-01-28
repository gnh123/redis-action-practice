package storage

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-redis/redis/v9"
)

func Test_SRem(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	rdb.SAdd(context.TODO(), "mytestset", "1", "2", "3", "4").Err()
	fmt.Println(rdb.SMembers(context.TODO(), "mytestset").Val())
	rdb.Del(context.TODO(), "mytestset", "1", "2", "3", "4").Err()
	fmt.Println(rdb.SMembers(context.TODO(), "mytestset").Val())
}

func Test_ZRem(t *testing.T) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	rdb.ZAdd(context.TODO(), "mytestzset", redis.Z{Member: "baidu.com", Score: 10}, redis.Z{Member: "qq.com", Score: 100}).Err()
	fmt.Println(rdb.ZRange(context.TODO(), "mytestzset", 0, -1).Val())
	fmt.Println(rdb.ZRem(context.TODO(), "mytestzset", []string{"baidu.com", "qq.com"}).Val())
	fmt.Println(rdb.ZRange(context.TODO(), "mytestzset", 0, -1).Val())
}

// TODO
func Test_ZInterStore(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	rdb.ZAdd(context.TODO(), "mytestzset", redis.Z{Member: "baidu.com", Score: 10}, redis.Z{Member: "qq.com", Score: 100}).Err()
	fmt.Println(rdb.ZRange(context.TODO(), "mytestzset", 0, -1).Val())
	rdb.ZInterStore(context.TODO(), "mytestzset", &redis.ZStore{Keys: []string{"mytestzset"}, Weights: []float64{0.5}}).Err()
	fmt.Println(rdb.ZRangeWithScores(context.TODO(), "mytestzset", 0, -1).Val())
}
