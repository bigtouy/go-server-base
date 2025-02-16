package repo

import "go-server-base/app/model"

type CacheRepo struct{}

type ICacheRepo interface {
	GetById(id string) (model.Cache, error)
	Create(id string) error
	Update(cache model.Cache) error
}

func NewICacheRepo() ICacheRepo {
	return &CacheRepo{}
}

func (c CacheRepo) GetById(id string) (model.Cache, error) {
	cacheFirst := model.Cache{}
	db := getDb().First(&cacheFirst, "id = ?", id)

	if db.Error != nil {
		return cacheFirst, db.Error
	}

	//db.RowsAffected 返回找到的数量
	return cacheFirst, nil
}

func (c CacheRepo) Create(id string) error {
	cache := model.Cache{
		ID:     id,
		Status: model.StatusCreate,
	}

	db := getDb().Create(&cache)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (c CacheRepo) Update(cache model.Cache) error {
	db := getDb().Save(&cache)
	if db.Error != nil {
		return db.Error
	}
	return nil
}
