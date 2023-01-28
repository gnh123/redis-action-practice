package storage

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"unsafe"

	"github.com/antlabs/gstl/cmp"
	"github.com/go-redis/redis/v9"
)

const (
	LIMIT = 1000
)

type Second struct {
	rdb *redis.Client
}

func New() *Second {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &Second{rdb: rdb}
}

func (s *Second) CheckToken(token string) string {
	return s.rdb.HGet(context.TODO(), "login:", token).Val()
}

func (s *Second) UpdateToken(token string, user string, item string) {

	timestamp := float64(time.Now().Unix())
	s.rdb.HSet(context.TODO(), "login:", token, user)
	// 记录令牌最后一次登录时间
	s.rdb.ZAdd(context.TODO(), "recent:", redis.Z{Member: token, Score: timestamp})

	if item != "" {
		s.rdb.ZAdd(context.TODO(), "viewed:"+token, redis.Z{Member: item, Score: timestamp}).Err()
		s.rdb.ZRemRangeByRank(context.TODO(), "viewed:"+token, 0, -26).Err()
		s.rdb.ZIncrBy(context.TODO(), "viewed", -1, item).Err()
	}
}

const limit = 100

func (s *Second) CleanSessions() {
	for {
		size := s.rdb.ZCard(context.TODO(), "recent:").Val()
		if size <= limit {
			time.Sleep(time.Second)
			continue
		}

		endIndex := cmp.Min(size-limit, 100)
		tokens := s.rdb.ZRange(context.TODO(), "recent:", 0, endIndex-1).Val()

		if err := s.rdb.Del(context.TODO(), tokens...).Err(); err != nil {

		}

		sessionKeys := make([]string, 0, len(tokens))
		for _, token := range tokens {
			sessionKeys = append(sessionKeys, fmt.Sprintf("viewed:%s"+token))
		}

		if err := s.rdb.Del(context.TODO(), sessionKeys...).Err(); err != nil {
			fmt.Printf("del:%s\n", err)
		}

		if err := s.rdb.HDel(context.TODO(), "login:", tokens...).Err(); err != nil {
			fmt.Printf("del:%s\n", err)
		}

		if err := s.rdb.ZRem(context.TODO(), "recent:", tokens).Err(); err != nil {
			fmt.Printf("del:%s\n", err)
		}
	}
}

func (s *Second) AddToCart(session string, item string, count int) {
	if count <= 0 {
		// 从购物车里面移除指定的商品
		s.rdb.HDel(context.TODO(), "cart:"+session, item)
	} else {
		// 将指定的商品添加到购物车
		s.rdb.HSet(context.TODO(), "cart:"+session, item, count)
	}
}

func (s *Second) CleanFullSessions() {
	for {
		size := s.rdb.ZCard(context.TODO(), "recent:").Val()
		if size < LIMIT {
			time.Sleep(time.Second)
		}

		endIndex := cmp.Min(int(size)-LIMIT, 1000)
		sessions := s.rdb.ZRange(context.TODO(), "recent:", 0, int64(endIndex)).Val()
		sessionsKeys := []string{}
		for _, sess := range sessions {
			sessionsKeys = append(sessionsKeys, "viewed:"+sess)
			sessionsKeys = append(sessionsKeys, "cart:"+sess)
		}

		s.rdb.Del(context.TODO(), sessionsKeys...).Err()
		s.rdb.HDel(context.TODO(), "login:", sessionsKeys...).Err()
		s.rdb.ZRem(context.TODO(), "recent:", sessionsKeys).Err()
	}
}

func (s *Second) CacheRequest(r *http.Request, callback func() string) string {
	if s.CanCache() {
		return callback()
	}

	pageKey := "cache:" + fmt.Sprintf("%x", unsafe.Pointer(r))
	val := s.rdb.Get(context.TODO(), pageKey).Val()
	if len(val) == 0 {
		val = callback()
		s.rdb.SetNX(context.TODO(), pageKey, val, 300)
	}

	return val
}

func extractItemID() string {
	return "x"
}

func (s *Second) RescaleViewed() {
	for {
		s.rdb.ZRemRangeByRank(context.TODO(), "viewed:", 0, -20001).Err()
		s.rdb.ZInterStore(context.TODO(), "viewed:", &redis.ZStore{Keys: []string{"viewed:"}, Weights: []float64{0.5}}).Err()
	}
}

func (s *Second) CanCache() bool {
	itemID := extractItemID()
	rank, err := s.rdb.ZRank(context.TODO(), "viewed:", itemID).Result()

	return err != nil && rank < 10000
}
