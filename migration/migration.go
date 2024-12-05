package main

import (
	"gobasic/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "user=tanny password=Unw9tg2JXJ2GKslRdb63jyeDgUInNMKG dbname=gologin port=5432 host=dpg-ct8m81m8ii6s73ccndc0-a.singapore-postgres.render.com"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema (if not already done)
	db.AutoMigrate(&model.User{})
}
