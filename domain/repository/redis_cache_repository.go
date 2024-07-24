package repository

import "time"

type RedisCacheRepository interface {
	GetHash(key string, property string, src interface{})
	DeleteHash(key string, property string) error
	SetHashObject(key string, property string, object interface{}) error
	SetExpireKey(key string, expiration time.Duration) error
}