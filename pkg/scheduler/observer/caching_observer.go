package observer

import (
	"automsg/pkg/cache"
	"context"
	"fmt"
	"log"
	"time"
)

const (
	defaultExpiration = 6 * time.Hour
	messageKeyPrefix  = "sent_message:"
)

type cachingObserver struct {
	redisClient cache.RedisClient
	expiration  time.Duration
}

func NewCachingObserver(redisClient cache.RedisClient) MessageObserver {
	return &cachingObserver{
		redisClient: redisClient,
		expiration:  defaultExpiration,
	}
}

func (r *cachingObserver) OnMessageProcessed(messageID string, success bool, err error) {
	if !success || messageID == "" {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key := fmt.Sprintf("%s%s", messageKeyPrefix, messageID)
	value := time.Now().Format(time.RFC3339)

	if err := r.redisClient.Set(ctx, key, value, r.expiration); err != nil {
		log.Printf("failed to cache message ID %s in Redis: %v", messageID, err)
		return
	}

}
