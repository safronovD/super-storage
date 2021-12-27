package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
)

var (
	ctx = context.Background()
)

type HashStorage interface {
	WriteNewHash(hashKey string, fileLink string, blockPosition int) error
	WriteBlockPosition(fileID string, hash string) error
	CheckFileID(fileID string) (int64, error)
	CheckHash(hash string) (int64, []string, error)
	Read(fileID string) ([][]string, error)
	Delete()
}

type redisHashStorageImpl struct {
	hashTable *redis.Client
	fileTable *redis.Client
}

func NewRedisHashImpl(cfg *Config) (HashStorage, error) {
	hashdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Hostname,
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.HashDB,
	})

	err := hashdb.Ping(ctx).Err()
	if err != nil {
		wrappedErr := fmt.Errorf("error connecting to hashdb in the redis: %w", err)
		return nil, wrappedErr
	}
	log.Printf("Successfuly connected to Redis Hash Table!")

	filedb := redis.NewClient(&redis.Options{
		Addr:     cfg.Hostname,
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.FileDB,
	})

	err = filedb.Ping(ctx).Err()
	if err != nil {
		wrappedErr := fmt.Errorf("error connecting to filedb in the redis: %w", err)
		return nil, wrappedErr
	}
	log.Printf("Successfuly connected to Redis File Table!")

	return &redisHashStorageImpl{
		hashTable: hashdb,
		fileTable: filedb,
	}, nil
}

func (client *redisHashStorageImpl) WriteNewHash(hashKey string, fileLink string, blockPosition int) error {
	err := client.hashTable.SAdd(ctx, hashKey, fileLink, blockPosition).Err()
	if err != nil {
		wrappedErr := fmt.Errorf("error creating new record in hash table: %w", err)
		return wrappedErr
	}
	log.Printf("New hash record created: %v", hashKey)
	return nil
}

func (client *redisHashStorageImpl) WriteBlockPosition(fileID string, hash string) error {
	err := client.fileTable.RPush(ctx, fileID, hash).Err()
	if err != nil {
		wrappedErr := fmt.Errorf("error adding new item in the file table record: %w", err)
		return wrappedErr
	}
	log.Printf("Added new hash for the file table FileID: %v, Hash: %v", fileID, hash)
	return nil
}

func (client *redisHashStorageImpl) CheckFileID(fileID string) (int64, error) {
	check, err := client.fileTable.Exists(ctx, fileID).Result()
	if err != nil {
		wrappedErr := fmt.Errorf("error finding key in file table: %w", err)
		return 0, wrappedErr
	}

	return check, nil
}

func (client *redisHashStorageImpl) CheckHash(hash string) (int64, []string, error) {
	var link []string
	check, err := client.hashTable.Exists(ctx, hash).Result()
	if err != nil {
		wrappedErr := fmt.Errorf("error finding key in hash table: %w", err)
		return 0, nil, wrappedErr
	}

	if check == 1 {
		link, err = client.hashTable.SMembers(ctx, hash).Result()
		if err != nil {
			wrappedErr := fmt.Errorf("error getting key in hash table: %w", err)
			return 0, nil, wrappedErr
		}
	}

	return check, link, nil
}

func (client *redisHashStorageImpl) Read(fileID string) ([][]string, error) {
	length, err := client.fileTable.LLen(ctx, fileID).Result()
	if err != nil {
		wrappedErr := fmt.Errorf("error reading lenght from file table: %w", err)
		return nil, wrappedErr
	}

	var (
		idx    int64
		result [][]string
	)

	for idx = 0; idx < length; idx++ {
		hash, err := client.fileTable.LIndex(ctx, fileID, idx).Result()
		if err != nil {
			wrappedErr := fmt.Errorf("error taking elem from list by index: %w", err)
			return nil, wrappedErr
		}
		link, err := client.hashTable.SMembers(ctx, hash).Result()

		if err != nil {
			wrappedErr := fmt.Errorf("error taking file link and pos by hash: %w", err)
			return nil, wrappedErr
		}
		result = append(result, link)
	}
	return result, nil
}

func (client *redisHashStorageImpl) Delete() {}
