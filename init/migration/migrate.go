package migration

import (
	"go-server-base/app/model"
	"go-server-base/global"
)

func Init() {
	db := global.DB
	err := db.AutoMigrate(&model.Cache{})
	if err != nil {
		panic(err)
	}
}
