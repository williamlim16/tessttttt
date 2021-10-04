package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBInit create connection to database
func DBInit() *gorm.DB {
	// dns := "root:@(localhost)/trash_separator?charset=utf8&parseTime=True&loc=Local"
	dns := "host=fanny.db.elephantsql.com port=5432 user=nckcqlju dbname=nckcqlju sslmode=disable password=dNBdEXEtLEQ6GAWLaCdPFH-iGJP9biV2"
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})

	if err != nil {
		panic("failed to connect to database")
	}

	return db
}
