package config

import (
	"log"
	"os"
	"time"
	"regexp"
	"github.com/go-redis/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"github.com/joho/godotenv"
)

const projectDirName = "trash-separator-api" // change to relevant project name

func goDotEnvVariable(key string) string {
    projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
    currentWorkDirectory, _ := os.Getwd()
    rootPath := projectName.Find([]byte(currentWorkDirectory))

    err := godotenv.Load(string(rootPath) + `/.env`)
  
	if err != nil {
	  log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
  }

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
	dns := goDotEnvVariable("DATABASE_CREDS")
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Printf("%v", err)
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
	address = goDotEnvVariable("REDIS_ADDR")
	password = goDotEnvVariable("REDIS_PASS")
	database = 0

	client = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       database,
	})

	return client
}
