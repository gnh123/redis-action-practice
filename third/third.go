package third

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
)

type Third struct {
	rdb *redis.Client
}

func New() *Third {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &Third{rdb: rdb}
}

func (t *Third) stringCmd() {
	t.rdb.Del(context.TODO(), "guo.num").Err()

	t.rdb.Incr(context.TODO(), "guo.num").Err()
	fmt.Println(t.rdb.Get(context.TODO(), "guo.num").Val())

	t.rdb.Decr(context.TODO(), "guo.num").Err()
	fmt.Println(t.rdb.Get(context.TODO(), "guo.num").Val())

	t.rdb.IncrBy(context.TODO(), "guo.num", 10).Err()
	fmt.Println(t.rdb.Get(context.TODO(), "guo.num").Val())

	t.rdb.DecrBy(context.TODO(), "guo.num", 10).Err()
	fmt.Println(t.rdb.Get(context.TODO(), "guo.num").Val())

	t.rdb.IncrByFloat(context.TODO(), "guo.num", 10.1).Err()
	fmt.Println(t.rdb.Get(context.TODO(), "guo.num").Val())

	t.rdb.Append(context.TODO(), "guo.num", "hello world")
	fmt.Println(t.rdb.Get(context.TODO(), "guo.num").Val())

	val := t.rdb.GetRange(context.TODO(), "guo.num", 0, 3).Val()
	fmt.Println(val)

	t.rdb.SetRange(context.TODO(), "guo.num", 0, "0").Err()
	fmt.Println(t.rdb.Get(context.TODO(), "guo.num").Val())

	t.rdb.SetBit(context.TODO(), "guo.bit", 0, 1).Err()
	fmt.Println("get bit")
	fmt.Println(t.rdb.GetBit(context.TODO(), "guo.bit", 0).Val())

	t.rdb.SetBit(context.TODO(), "guo.bit", 0, 1).Err()
	fmt.Println(t.rdb.GetBit(context.TODO(), "guo.bit", 0).Val())

	val2 := t.rdb.BitCount(context.TODO(), "guo.bit", &redis.BitCount{
		Start: 0,
		End:   1,
	}).Val()
	fmt.Println("####", val2)

	t.rdb.SetBit(context.TODO(), "guo.bit1", 0, 1).Err()
	fmt.Println(t.rdb.GetBit(context.TODO(), "guo.bit1", 0).Val())

	t.rdb.SetBit(context.TODO(), "guo.bit2", 1, 1).Err()
	fmt.Println(t.rdb.GetBit(context.TODO(), "guo.bit2", 0).Val())

	t.rdb.BitOpAnd(context.TODO(), "guo.bit.op", "guo.bit1", "guo.bit2")
	fmt.Println("###", t.rdb.GetBit(context.TODO(), "guo.bit.op", 0).Val())
}

func (t *Third) listCmd() {
	t.rdb.Del(context.TODO(), "guo")
	t.rdb.RPush(context.TODO(), "guo", 1, 2, 3, 4, 5).Err()
	t.rdb.LPush(context.TODO(), "guo", -1, -2, -3, -4, -5).Err()
	fmt.Println(t.rdb.LRange(context.TODO(), "guo", 0, -1).Result())
	t.rdb.RPop(context.TODO(), "guo")
	t.rdb.LPop(context.TODO(), "guo")
	fmt.Println(t.rdb.LRange(context.TODO(), "guo", 0, -1).Result())
	fmt.Println(t.rdb.LIndex(context.TODO(), "guo", 0).Result())
	t.rdb.LTrim(context.TODO(), "guo", 2, 3)
	fmt.Println(t.rdb.LRange(context.TODO(), "guo", 0, -1).Result())
}

func (t *Third) setCmd() {
	// 删除
	t.rdb.Del(context.TODO(), "guo").Err()
	// 添加
	t.rdb.SAdd(context.TODO(), "guo", "a", "b", "c", "e", "f").Err()
	fmt.Println(t.rdb.SMembers(context.TODO(), "guo").Result())

	// 删除
	t.rdb.SRem(context.TODO(), "guo", "c", "d").Err()
	fmt.Println(t.rdb.SMembers(context.TODO(), "guo").Result())

	// 计数
	fmt.Println(t.rdb.SCard(context.TODO(), "guo").Result())

	t.rdb.SAdd(context.TODO(), "guo2", "1", "2", "3", "4", "5").Err()
	// 把某个集合的成员移动到另一个成员里面
	t.rdb.SMove(context.TODO(), "guo", "guo2", "a").Err()
	fmt.Println(t.rdb.SMembers(context.TODO(), "guo2").Result())
	// 某个成员在集合中
	fmt.Println(t.rdb.SIsMember(context.TODO(), "guo2", "a").Result())
	t.rdb.SPop(context.TODO(), "guo2").Err()
	fmt.Println(t.rdb.SMembers(context.TODO(), "guo2").Result())
	fmt.Println(t.rdb.SRandMemberN(context.TODO(), "guo2", 3).Result())

	t.rdb.Del(context.TODO(), "guo").Err()
	t.rdb.SAdd(context.TODO(), "guo", "k1", "k2").Err()
	t.rdb.SAdd(context.TODO(), "guo2", "k2", "k3").Err()
	fmt.Println(t.rdb.SDiff(context.TODO(), "guo", "guo2").Result())
	fmt.Println(t.rdb.SInter(context.TODO(), "guo", "guo2").Result())
	fmt.Println(t.rdb.SUnion(context.TODO(), "guo", "guo2").Result())
}
