package db

import (
	"gorm.io/gorm/clause"
)

type KeyValue struct {
	Key   string `gorm:"primaryKey"`
	Value string
}

func (c *Client) Get(key string) (string, error) {
	var kv KeyValue
	if err := c.db.Where("key = ?", key).First(&kv).Error; err != nil {
		return "", err
	}
	return kv.Value, nil
}

func (c *Client) Set(key, value string) error {
	kv := KeyValue{Key: key, Value: value}
	return c.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&kv).Error
}
