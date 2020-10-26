package fetcher

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type currency struct {
	Id int
	Code string
}

func (currency) TableName() string {
	return "ISO4217"
}

type sqLiteFetcher struct {
	db *gorm.DB
}

func (s sqLiteFetcher) Fetch(id int) (string, error) {
	item := currency{Id: id}
	if err :=  s.db.Where(item).First(&item).Error; err != nil {
		return "", err
	}
	return item.Code, nil
}

func (s sqLiteFetcher) FetchAll() (map[int]string, error) {
	var items []currency
	if err := s.db.Find(&items).Error; err != nil {
		return nil, err
	}
	result := make(map[int]string)
	for _, item := range items {
		result[item.Id] = item.Code
	}
	return result, nil
}

func NewSQLiteFetcher(path string) (Fetcher, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &sqLiteFetcher{db: db}, nil
}
