package internal

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheClient struct {
	rdsClient *redis.Client
}

type JSONWrapper struct {
	Data Weather
}

func NewCacheClient(client *redis.Client) *CacheClient {
	return &CacheClient{rdsClient: client}
}

func (m *JSONWrapper) MarshalBinary() (data []byte, err error) {
	marshalledData, err := json.Marshal(m.Data)
	if err != nil {
		return nil, err
	}
	return marshalledData, nil
}

func (m *JSONWrapper) UnmarshalBinary(data []byte) (err error) {
	err = json.Unmarshal(data, &m.Data)
	if err != nil {
		return err
	}
	return nil
}

func (c *CacheClient) searchCache(
	ctx context.Context,
	cacheKey string,
	response *JSONWrapper,
) (err error) {
	getTransaction := c.rdsClient.Get(ctx, cacheKey)
	if getTransaction.Err() != nil {
		return errors.New("cache-miss")
	}

	binaryRedisHit, err := getTransaction.Bytes()
	if err != nil {
		return errors.Join(err, errors.New("while converting cache-hit values to bytes"))
	}

	err = response.UnmarshalBinary(binaryRedisHit)
	if err != nil {
		return errors.Join(err, errors.New("while unmarshalling bytes from cache-hit"))
	}
	return nil
}

func (c *CacheClient) storeInCache(
	ctx context.Context,
	cacheKey string,
	response *JSONWrapper,
) (err error) {
	resBytes, err := response.MarshalBinary()
	if err != nil {
		return errors.Join(errors.New("while serializing response for redis"), err)
	}
	transactionStatus := c.rdsClient.Set(ctx, cacheKey, resBytes, time.Hour)
	if transactionStatus.Err() != nil {
		return errors.Join(
			errors.New("while storing response in redis"),
			transactionStatus.Err(),
		)
	}
	return nil
}
