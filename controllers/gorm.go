package controllers

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type InDB struct {
	DB          *gorm.DB
	RedisClient *redis.Client
}
