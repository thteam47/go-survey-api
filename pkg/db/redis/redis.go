package redisrepo

import "time"

type RedisRepository interface {
	GetValueCache(key string, result interface{}) error
	SetValueCache(key string, result interface{}, exp time.Duration) error
	RemoveValueCache(key string) error
	SetKeyToListKeyCache(key string, listCache string)
	RemoveListKeyCache(keyListCache string)
}
