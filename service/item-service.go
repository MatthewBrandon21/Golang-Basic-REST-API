package service

import (
	"basic-rest-api-jwt-mysql/dto"
	"basic-rest-api-jwt-mysql/entity"
	"basic-rest-api-jwt-mysql/repository"
	"fmt"
	"log"

	"github.com/mashingan/smapping"
)

type ItemService interface {
	Insert(b dto.ItemCreateDTO) entity.Item
	Update(b dto.ItemUpdateDTO) entity.Item
	Delete(b entity.Item)
	All() []entity.Item
	FindByID(itemID uint64) entity.Item
	IsAllowedToEdit(userID string, itemID uint64) bool
}

type itemService struct {
	itemRepository repository.ItemRepository
}

func NewItemService(itemRepo repository.ItemRepository) ItemService {
	return &itemService{
		itemRepository: itemRepo,
	}
}

func (service *itemService) Insert(b dto.ItemCreateDTO) entity.Item {
	item := entity.Item{}
	err := smapping.FillStruct(&item, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}
	res := service.itemRepository.InsertItem(item)
	return res
}

func (service *itemService) Update(b dto.ItemUpdateDTO) entity.Item {
	item := entity.Item{}
	err := smapping.FillStruct(&item, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}
	res := service.itemRepository.UpdateItem(item)
	return res
}

func (service *itemService) Delete(b entity.Item) {
	service.itemRepository.DeleteItem(b)
}

func (service *itemService) All() []entity.Item {
	return service.itemRepository.AllItem()
}

func (service *itemService) FindByID(itemID uint64) entity.Item {
	return service.itemRepository.FindItemByID(itemID)
}

func (service *itemService) IsAllowedToEdit(userID string, itemID uint64) bool {
	b := service.itemRepository.FindItemByID(itemID)
	id := fmt.Sprintf("%v", b.UserID)
	return userID == id
}
