package repository

import (
	"basic-rest-api-jwt-mysql/entity"

	"gorm.io/gorm"
)

type ItemRepository interface {
	InsertItem(b entity.Item) entity.Item
	UpdateItem(b entity.Item) entity.Item
	DeleteItem(b entity.Item)
	AllItem() []entity.Item
	FindItemByID(itemID uint64) entity.Item
}

type itemConnection struct {
	connection *gorm.DB
}

func NewItemRepository(dbConn *gorm.DB) ItemRepository {
	return &itemConnection{
		connection: dbConn,
	}
}

func (db *itemConnection) InsertItem(b entity.Item) entity.Item {
	db.connection.Save(&b)
	db.connection.Preload("User").Find(&b)
	return b
}

func (db *itemConnection) UpdateItem(b entity.Item) entity.Item {
	db.connection.Save(&b)
	db.connection.Preload("User").Find(&b)
	return b
}

func (db *itemConnection) DeleteItem(b entity.Item) {
	db.connection.Delete(&b)
}

func (db *itemConnection) FindItemByID(itemID uint64) entity.Item {
	var item entity.Item
	db.connection.Preload("User").Find(&item, itemID)
	return item
}

func (db *itemConnection) AllItem() []entity.Item {
	var items []entity.Item
	db.connection.Preload("User").Find(&items)
	return items
}
