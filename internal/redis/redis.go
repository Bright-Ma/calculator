package redis

import (
	"calculator/internal/database"
	"calculator/internal/model"
	"context"
	"fmt"
	"math"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// Redis key 前缀
	QuestionKeyPrefix = "question:"
	HourlyRankKey     = "rank:hourly"
	DailyRankKey      = "rank:daily"
	// 时间衰减因子（24小时）
	TimeDecayFactor = 24 * time.Hour
)

// Redis Redis客户端封装
type Redis struct {
	Client *redis.Client
}

// NewRedis 创建新的Redis客户端
func NewRedis() *Redis {
	client := &Redis{
		Client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // 如果有密码，在这里设置
			DB:       0,  // 使用默认DB
		}),
	}

	// 启动定期更新排行榜的goroutine
	go client.startPeriodicRankingUpdate()

	return client
}

// startPeriodicRankingUpdate 启动定期更新排行榜的goroutine
func (r *Redis) startPeriodicRankingUpdate() {
	ticker := time.NewTicker(5 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		if err := r.InitRankingData(); err != nil {
			fmt.Printf("更新排行榜数据失败: %v\n", err)
		}
	}
}

// InitRankingData 从MySQL初始化排行榜数据
func (r *Redis) InitRankingData() error {
	ctx := context.Background()

	// 清空现有的排行榜数据
	if err := r.Client.Del(ctx, HourlyRankKey, DailyRankKey).Err(); err != nil {
		return fmt.Errorf("清空排行榜数据失败: %v", err)
	}

	// 从MySQL获取所有用户的历史记录
	var historyRecords []model.HistoryRecord
	if err := database.DB.Find(&historyRecords).Error; err != nil {
		return fmt.Errorf("获取历史记录失败: %v", err)
	}

	// 按用户ID分组统计
	userScores := make(map[uint]float64)
	for _, record := range historyRecords {
		// 计算基础分数（答对2分，答错1分）
		score := 1.0
		if record.IsCorrect {
			score = 2.0
		}

		// 应用时间衰减
		timeDecay := math.Exp(-float64(record.CreatedAt.Unix()) / float64(TimeDecayFactor.Seconds()))
		finalScore := score * timeDecay

		// 累加用户分数
		userScores[record.UserID] += finalScore
	}

	// 将用户分数写入Redis
	for userID, score := range userScores {
		// 更新小时榜
		if err := r.Client.ZAdd(ctx, HourlyRankKey, redis.Z{
			Score:  score,
			Member: fmt.Sprintf("%d", userID),
		}).Err(); err != nil {
			return fmt.Errorf("更新小时榜失败: %v", err)
		}

		// 更新日榜
		if err := r.Client.ZAdd(ctx, DailyRankKey, redis.Z{
			Score:  score,
			Member: fmt.Sprintf("%d", userID),
		}).Err(); err != nil {
			return fmt.Errorf("更新日榜失败: %v", err)
		}
	}

	return nil
}

// UpdateUserHotScore 更新用户热度值
func (r *Redis) UpdateUserHotScore(ctx context.Context, userID uint, isCorrect bool) error {
	// 更新所有排行榜
	rankKeys := []string{HourlyRankKey, DailyRankKey}
	for _, key := range rankKeys {
		// 基础分数增量
		scoreIncrement := 50.0 // 每次答题基础加分

		// 答对额外加分
		if isCorrect {
			scoreIncrement += 100.0 // 答对额外加100分
		}

		// 使用 ZINCRBY 增加分数
		err := r.Client.ZIncrBy(ctx, key, scoreIncrement, fmt.Sprintf("%d", userID)).Err()
		if err != nil {
			return fmt.Errorf("更新排行榜 %s 失败: %v", key, err)
		}
	}

	return nil
}

// GetHotRanking 获取热度排行榜
func (r *Redis) GetHotRanking(ctx context.Context, rankType string, limit int64) ([]RankingItem, error) {
	var rankKey string
	switch rankType {
	case "hourly":
		rankKey = HourlyRankKey
	default:
		rankKey = DailyRankKey
	}

	// 获取排行榜数据
	result, err := r.Client.ZRevRangeWithScores(ctx, rankKey, 0, limit-1).Result()
	if err != nil {
		return nil, fmt.Errorf("获取排行榜失败: %v", err)
	}

	// 构建返回数据
	rankings := make([]RankingItem, 0, len(result))
	for i, item := range result {
		// 从用户ID获取用户名
		userID := item.Member.(string)

		var user model.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			// 如果找不到用户，使用默认名称
			rankings = append(rankings, RankingItem{
				Rank:     i + 1,
				UserID:   userID,
				Username: fmt.Sprintf("用户%s", userID),
				HotScore: item.Score,
			})
			continue
		}

		rankings = append(rankings, RankingItem{
			Rank:     i + 1,
			UserID:   userID,
			Username: user.Username,
			HotScore: item.Score,
		})
	}

	return rankings, nil
}

// RankingItem 排行榜项目
type RankingItem struct {
	Rank     int     `json:"rank"`
	UserID   string  `json:"user_id"`
	Username string  `json:"username"`
	HotScore float64 `json:"hot_score"`
}
