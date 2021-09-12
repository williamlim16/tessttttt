package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DBInit create connection to database
func DBInit() *gorm.DB {
	dns := "root:@(localhost)/trash_separator?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	return db
}
