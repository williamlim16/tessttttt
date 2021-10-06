package config

import (
	"log"
	"os"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBInit create connection to database
func DBInit() *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)

	// dns := "root:@(localhost)/trash_separator?charset=utf8&parseTime=True&loc=Local"
	dns := "host=fanny.db.elephantsql.com port=5432 user=nckcqlju dbname=nckcqlju sslmode=disable password=dNBdEXEtLEQ6GAWLaCdPFH-iGJP9biV2"
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic("failed to connect to database")
	}

	return db
}

func RedisInit() *redis.Client {
	var (
		client   *redis.Client
		address  string
		password string
		database int
	)
	address = "localhost:6379"
	password = ""
	database = 0

	client = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       database,
	})

	return client
}
