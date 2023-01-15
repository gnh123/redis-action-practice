package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gnh123/redis-action-practice/first/internal/types"
	"github.com/go-redis/redis/v8"
)

type Storage struct {
	rdb *redis.Client
}

const (
	// 86400是一天的秒数 24 * 3600
	ONE_WEEK_IN_SECONDS = 7 * 86400
	VOTE_SCORE          = 86400 / 24
	ARTICLES_PER_PAGE   = 25
)

func New() *Storage {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &Storage{rdb: rdb}
}

// user -> user:id
// article article:id
func (s *Storage) ArticleVote(user string, article string) {

	cutoff := time.Now().Unix() - ONE_WEEK_IN_SECONDS
	// time zset存放所有的文章，值就是时间, 把这篇文章的时间取出来，所以超过7天，直接不能投票
	if s.rdb.ZScore(context.TODO(), "time:", article).Val() < float64(cutoff) {
		return
	}

	// 取出文章id
	articleID := article
	if pos := strings.Index(article, ":"); pos != -1 {
		articleID = article[pos:]
	}

	// 记录每篇文章投过票的用户,如果没有投过票
	if s.rdb.SAdd(context.TODO(), "voted:"+articleID, user).Val() == 1 {
		//这篇文章分数增加, 变更都记录到一个zset里面
		if err := s.rdb.ZIncrBy(context.TODO(), "score:", VOTE_SCORE, article).Err(); err != nil {
			fmt.Printf("score fail:%s\n", err)
		}

		// 修改哈希表里面文章的分数
		if err := s.rdb.HIncrBy(context.TODO(), article, "votes", 1).Err(); err != nil {
			fmt.Printf("")
		}
	}
}

func (s *Storage) PostArticle(req *types.ArticleReq) int64 {
	articleID := s.rdb.Incr(context.TODO(), "article:").Val()

	// user id
	user := req.Poster
	article := fmt.Sprintf("article:%d", articleID)
	// 使用set记录这篇文章被多少用户投过票
	voted := fmt.Sprintf("voted:%d", articleID)
	if err := s.rdb.SAdd(context.TODO(), voted, user).Err(); err != nil {
		fmt.Printf("post Article: fail %s\n", err)
	}

	now := time.Now().Unix()
	req.Time = fmt.Sprintf("%d", now)
	req.Votes = 1

	s.rdb.HMSet(context.TODO(), article, req)

	s.rdb.ZAdd(context.TODO(), "score:", &redis.Z{Score: float64(now + VOTE_SCORE), Member: article})
	s.rdb.ZAdd(context.TODO(), "time:", &redis.Z{Score: float64(now), Member: article})
	return articleID
}

// 获取文章列表，先访问score那个zset。这里面是文章id的zset。再根据zset获取内容
func (s *Storage) GetArticles(page int64, order string) {
	var start int64 = (page - 1) * ARTICLES_PER_PAGE
	var end int64 = start + ARTICLES_PER_PAGE - 1

	if order == "" {
		order = "score:"
	}

	ids := s.rdb.ZRevRange(context.TODO(), order, start, end).Val()
	for _, id := range ids {
		articleData := s.rdb.HGetAll(context.TODO(), id).Val()
		articleData["id"] = id
		fmt.Println(articleData)
	}
}

func (s *Storage) AddRemoveGroups(articleID string, toAdd []string, toRemove []string) {
	article := "article:" + articleID
	for _, group := range toAdd {
		s.rdb.SAdd(context.TODO(), "group:"+group, article)
	}

	for _, group := range toRemove {
		s.rdb.SRem(context.TODO(), "group:"+group, article)
	}
}

func (s *Storage) GetGroupArticles(group string, page int64) {
	key := "score:" + group

	// 不存在
	if s.rdb.Exists(context.TODO(), key).Val() == 0 {
		s.rdb.ZInterStore(context.TODO(), key, &redis.ZStore{Keys: []string{"group:" + group, "score"}})
		s.rdb.Expire(context.TODO(), key, 60)
	}

	s.GetArticles(page, key)
}
