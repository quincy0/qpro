package qRedis

import (
	"context"
	"fmt"
	"time"

	"github.com/quincy0/qpro/qConfig"
	"github.com/quincy0/qpro/qLog"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type AeClient struct {
	client *redis.Client
}

var Client *AeClient

func InitRedis() {
	newClient := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", qConfig.Settings.Redis.Host, qConfig.Settings.Redis.Port),
		Password:     qConfig.Settings.Redis.Password,
		DB:           qConfig.Settings.Redis.DataBase,
		PoolFIFO:     true,
		PoolSize:     qConfig.Settings.Redis.MaxIdleConn,
		MinIdleConns: qConfig.Settings.Redis.MinIdleConn,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			qLog.TraceInfo(ctx, "redis connect")
			return nil
		},
	})
	// Enable tracing instrumentation.
	if err := redisotel.InstrumentTracing(newClient); err != nil {
		panic(err)
	}

	// Enable metrics instrumentation.
	if err := redisotel.InstrumentMetrics(newClient); err != nil {
		panic(err)
	}
	pong := newClient.Ping(context.Background()).Val()
	if pong != "PONG" {
		panic("redis " + pong)
	}
	Client = &AeClient{
		client: newClient,
	}
}

func (r *AeClient) Close() error {
	return r.client.Close()
}

func (r *AeClient) Set(ctx context.Context, key, value string, exp time.Duration) error {
	cmd := r.client.Set(ctx, key, value, exp)
	return cmd.Err()
}

func (r *AeClient) SetNX(ctx context.Context, key string, value interface{}, exp time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, value, exp).Result()
}

func (r *AeClient) SetEx(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
	return r.client.SetEx(ctx, key, value, expiration).Result()
}

func (r *AeClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *AeClient) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *AeClient) HSet(ctx context.Context, key string, values ...interface{}) error {
	return r.client.HSet(ctx, key, values...).Err()
}

func (r *AeClient) HDel(ctx context.Context, key string, field ...string) error {
	return r.client.HDel(ctx, key, field...).Err()
}

func (r *AeClient) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	return r.client.HGet(ctx, key, field)
}

func (r *AeClient) HMGetCmd(ctx context.Context, key string, fields ...string) *redis.SliceCmd {
	return r.client.HMGet(ctx, key, fields...)
}

func (r *AeClient) HGetCmd(ctx context.Context, key string) *redis.MapStringStringCmd {
	return r.client.HGetAll(ctx, key)
}

func (r *AeClient) HGetAll(ctx context.Context, key string, model interface{}) error {
	return r.client.HGetAll(ctx, key).Scan(&model)
}

func (r *AeClient) HKEYS(ctx context.Context, key string) ([]string, error) {
	return r.client.HKeys(ctx, key).Result()
}

func (r *AeClient) HGetJson(ctx context.Context, key, field string) (string, error) {
	cmd := r.client.HGet(ctx, key, field)
	return cmd.Result()
}

func (r *AeClient) LPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	cmd := r.client.LPush(ctx, key, values...)
	return cmd.Result()
}

func (r *AeClient) RPop(ctx context.Context, key string) (string, error) {
	cmd := r.client.RPop(ctx, key)
	return cmd.Result()
}

func (r *AeClient) LPushRPop(ctx context.Context, source, destination string) (string, error) {
	cmd := r.client.RPopLPush(ctx, source, destination)
	return cmd.Result()
}

func (r *AeClient) LLen(ctx context.Context, key string) int64 {
	cmd := r.client.LLen(ctx, key)
	return cmd.Val()
}

func (r *AeClient) Expire(ctx context.Context, key string, expiration time.Duration) bool {
	cmd := r.client.Expire(ctx, key, expiration)
	return cmd.Val()
}

func (r *AeClient) TTL(ctx context.Context, key string) time.Duration {
	cmd := r.client.TTL(ctx, key)
	return cmd.Val()
}

func (r *AeClient) Exists(ctx context.Context, key string) bool {
	cmd := r.client.Exists(ctx, key)
	return cmd.Val() == 1
}

func (r *AeClient) HExists(ctx context.Context, key, filed string) bool {
	cmd := r.client.HExists(ctx, key, filed)
	return cmd.Val()
}

func (r *AeClient) Persist(ctx context.Context, key string) bool {
	return r.client.Persist(ctx, key).Val()
}

func (r *AeClient) Unlink(ctx context.Context, key string) error {
	return r.client.Unlink(ctx, key).Err()
}

func (r *AeClient) Pipeline(ctx context.Context, fn func(redis.Pipeliner) error) error {
	_, err := r.client.Pipelined(ctx, fn)
	return err
}

func (r *AeClient) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return r.client.IncrBy(ctx, key, value).Result()
}

func (r *AeClient) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

func (r *AeClient) ZAdd(ctx context.Context, key string, members ...redis.Z) (int64, error) {
	return r.client.ZAddArgs(ctx, key, redis.ZAddArgs{Members: members}).Result()
}

func (r *AeClient) ZScore(ctx context.Context, key, member string) (float64, error) {
	return r.client.ZScore(ctx, key, member).Result()
}

func (r *AeClient) ZRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return r.client.ZRem(ctx, key, members).Result()
}

func (r *AeClient) ZRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) ([]redis.Z, error) {
	return r.client.ZRangeByScoreWithScores(ctx, key, opt).Result()
}

func (r *AeClient) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error) {
	return r.client.ZRangeByScore(ctx, key, opt).Result()
}

func (r *AeClient) ZRemRangeByScore(ctx context.Context, key, min, max string) (int64, error) {
	return r.client.ZRemRangeByScore(ctx, key, min, max).Result()
}

func (r *AeClient) HMSet(ctx context.Context, key string, values ...interface{}) (bool, error) {
	return r.client.HMSet(ctx, key, values).Result()
}

func (r *AeClient) HMSetSetEx(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return r.client.HMSet(ctx, key, value, expiration).Result()
}

func (r *AeClient) HIncrBy(ctx context.Context, key, field string, incr int64) (int64, error) {
	return r.client.HIncrBy(ctx, key, field, incr).Result()
}
